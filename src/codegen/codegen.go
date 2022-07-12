package codegen

import (
	"archive/zip"

	plugin "codegenfornet/src/plugins"
	"codegenfornet/src/plugins/struct_crawl"
	"codegenfornet/util"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"unicode"
)

type CodegenForNet struct {
	DbType      string
	DbName      string
	SourceDb    bool
	ProjectName string
	TempPath    string
	carwl       plugin.DbStructureGet
}

type TableTempData struct {
	TableName   string
	ProjectName string
	DbName      string
	PKList      []TableField
	Fields      []TableField
}
type TableField struct {
	IsPK  string
	Type  string
	Field string
	Desc  string
}

type TableList struct {
	ProjectName string
	DbName      string
	Tables      []string
}

// 需要根据表生成多个的文件
var TableUnionFile = map[string]struct{}{
	"Request.cs":    {},
	"Service.cs":    {},
	"Controller.cs": {},
	"Model.cs":      {},
}

// 单文件内表信息生成
var TableFieldOnce = map[string]struct{}{
	"Context.cs": {},
}

var handleLimit = 100

func (cg *CodegenForNet) ServeStart() error {
	var tables []plugin.Tables
	if cg.SourceDb {
		var err error
		cg.carwl, err = struct_crawl.GetCrawl(cg.DbType)
		if err != nil {
			return err
		}
		tables = cg.carwl.GetTables(cg.DbName)
		if len(tables) <= 0 {
			return errors.New("没有可执行的表")
		}
	} else {
		tables = make([]plugin.Tables, 1)
		tables[0] = plugin.Tables{TableName: "temp_table", Desc: ""}
	}

	zr, err := zip.OpenReader(cg.TempPath)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(cg.ProjectName, 0755); err != nil {
		return err
	}

	for _, item := range zr.File {
		path := filepath.Join(cg.ProjectName, item.Name)
		path = strings.Replace(path, "TempTable", cg.ProjectName, -1)

		if item.FileInfo().IsDir() {
			err = os.MkdirAll(cg.ProjectName, 0755)
			if err != nil {
				return err
			}
			continue
		}
		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("模板文件:%v -> %v\n", item.FileInfo().Name(), path)
		fr, err := item.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		tmplContent, _ := ioutil.ReadAll(fr)
		fileName := item.FileInfo().Name()

		if _, ok := TableUnionFile[fileName]; ok {
			// 需要与表一对一生成的文件
			wg := &sync.WaitGroup{}
			ch_run := make(chan struct{}, handleLimit)
			for _, table := range tables {
				wg.Add(1)
				ch_run <- struct{}{}

				fmt.Printf("开始渲染：%v\n", table.TableName)
				go func(s_fileName, s_path, s_tmplContent, s_table string) {
					defer func() {
						<-ch_run
						wg.Done()
					}()
					err := cg.genTableFile(s_fileName, s_path, s_tmplContent, s_table)
					if err != nil {
						fmt.Printf("%v\n", err)
					}
				}(fileName, path, string(tmplContent), table.TableName)
			}
			wg.Wait()
		} else if _, ok := TableFieldOnce[fileName]; ok {
			tmplData := &TableList{
				ProjectName: cg.ProjectName,
				DbName:      cg.DbName,
				Tables:      make([]string, 0),
			}
			for _, t := range tables {
				tmplData.Tables = append(tmplData.Tables, nameHandler(t.TableName))
			}
			tmpl, err := template.New(fileName).Parse(string(tmplContent))
			if err != nil {
				fmt.Printf("%v\n", err)
				return err
			}
			path = strings.Replace(path, fileName, nameHandler(tmplData.DbName)+fileName, -1)
			err = util.CreateFileBytes(path, func(f *os.File) error {
				e := tmpl.Execute(f, tmplData)
				if e != nil {
					return e
				}
				return nil
			})
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			// fmt.Printf("一次性模板文件:%v -> %v\n", fileName, path)
			tmplData := &TableTempData{
				ProjectName: cg.ProjectName,
				DbName:      cg.DbName,
			}
			tmpl, err := template.New(fileName).Parse(string(tmplContent))
			if err != nil {
				fmt.Printf("%v\n", err)
				return err
			}
			err = util.CreateFileBytes(path, func(f *os.File) error {
				e := tmpl.Execute(f, tmplData)
				if e != nil {
					return e
				}
				return nil
			})
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
	return nil
}

func (cg *CodegenForNet) genTableFile(fileName, path, tmplContent, table string) error {
	var tmplData *TableTempData
	if !cg.SourceDb {
		tmplData = &TableTempData{
			ProjectName: cg.ProjectName,
			DbName:      cg.DbName,
			TableName:   nameHandler("default_table"),
			Fields:      []TableField{},
		}
	} else {
		tmplData = &TableTempData{
			ProjectName: cg.ProjectName,
			DbName:      cg.DbName,
			TableName:   nameHandler(table),
		}

		columns := cg.carwl.GetColumns(cg.DbName, table)
		tmplData.Fields = make([]TableField, len(columns))
		tmplData.PKList = make([]TableField, 0)

		for i := 0; i < len(columns); i++ {
			if columns[i].Key == "PRI" {
				tmplData.PKList = append(tmplData.PKList, TableField{
					Type:  cg.carwl.TypeSwitch(columns[i].Type),
					Field: nameHandler(columns[i].Field),
					Desc:  columns[i].Desc,
					IsPK:  columns[i].Key,
				})
			}
			tmplData.Fields[i] = TableField{
				Type:  cg.carwl.TypeSwitch(columns[i].Type),
				Field: nameHandler(columns[i].Field),
				Desc:  columns[i].Desc,
				IsPK:  columns[i].Key,
			}
		}

	}
	tmpl, err := template.New(table).Parse(string(tmplContent))
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	path = strings.Replace(path, fileName, tmplData.TableName+fileName, -1)
	err = util.CreateFileBytes(path, func(f *os.File) error {
		e := tmpl.Execute(f, tmplData)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return nil
}

func nameHandler(name string) string {
	var tmp string
	strs := strings.Split(name, "_")
	for _, str := range strs {
		if len(str) <= 0 {
			continue
		}
		b := []byte(str)
		b[0] = byte(unicode.ToUpper(rune(b[0])))
		tmp += string(b)
	}
	return tmp
}
