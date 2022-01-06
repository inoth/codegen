package templet

import (
	"codegen/db"
	"codegen/util"
	"os"
	"text/template"

	"fmt"
	"unicode"

	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	FileTemplate map[string]IHandler
	DataTypeMap  map[string]string
	projectName  string
	dbName       string
	useDbName    bool
)

func init() {
	DataTypeMap = make(map[string]string)
	DataTypeMap["varchar"] = "string"
	DataTypeMap["decimal"] = "decimal"
	DataTypeMap["int"] = "int"
	DataTypeMap["tinyint"] = "int"
	DataTypeMap["datetime"] = "DateTime"

	FileTemplate = make(map[string]IHandler)
	FileTemplate["Entity"] = &EntityHandle{
		Folder: "Entity",
		Tmpl:   Entity}

	FileTemplate["Model"] = &EntityHandle{
		Folder: "Model",
		Tmpl:   Model}

	FileTemplate["AutoMapper"] = &EntityHandle{
		Folder: "AutoMapper",
		Tmpl:   Mapper}

	FileTemplate["IRepository"] = &NormalHandle{
		Folder: "IRepository",
		Tmpl:   IRepository}

	FileTemplate["Repository"] = &NormalHandle{
		Folder: "Repository",
		Tmpl:   Repository}

	FileTemplate["IService"] = &NormalHandle{
		Folder: "IService",
		Tmpl:   IService}

	FileTemplate["Service"] = &NormalHandle{
		Folder: "Service",
		Tmpl:   Service}

	FileTemplate["Controller"] = &NormalHandle{
		Folder: "Controller",
		Tmpl:   Controller}
}

type TmplData struct {
	ProjectName string
	TableName   string
	DbTableName string
	UseDbName   bool
	Fields      []Field
}
type Field struct {
	DbField string
	Field   string
	Key     string
	IsNull  string
	Desc    string
	DbType  string
}

type IHandler interface {
	Process(wg *sync.WaitGroup, tableName string)
}

type NormalHandle struct {
	Folder string
	Tmpl   string
}

func (e *NormalHandle) Process(wg *sync.WaitGroup, tableName string) {
	defer wg.Done()
	tmpl, err := template.New(tableName).Parse(e.Tmpl)
	if err != nil {
		log.Error(err.Error())
		return
	}
	data := &TmplData{
		ProjectName: projectName,
		TableName:   NameHandler(tableName),
		DbTableName: tableName,
		UseDbName:   useDbName,
	}
	path := fmt.Sprintf("./%v/%v/%v%v.cs", projectName, e.Folder, NameHandler(tableName), e.Folder)
	err = util.CreateFileBytes(path, func(f *os.File) error {
		e := tmpl.Execute(f, data)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
	}
}

type EntityHandle struct {
	Folder string
	Tmpl   string
}

func (e *EntityHandle) Process(wg *sync.WaitGroup, tableName string) {
	defer wg.Done()

	cols := db.GetColumns(dbName, tableName)
	if len(cols) <= 0 {
		log.Errorf("%v: 未找到有效列", tableName)
		return
	}
	tmpl, err := template.New(tableName).Parse(e.Tmpl)
	if err != nil {
		log.Error(err.Error())
		return
	}
	data := &TmplData{
		ProjectName: projectName,
		TableName:   NameHandler(tableName),
		DbTableName: tableName,
		UseDbName:   useDbName,
	}
	data.Fields = make([]Field, len(cols))
	for i, col := range cols {
		data.Fields[i] = Field{DbType: MatchType(col.DataType),
			DbField: col.ColName,
			Field:   NameHandler(col.ColName),
			Desc:    col.ColDesc}
	}
	path := fmt.Sprintf("./%v/%v/%v%v.cs", projectName, e.Folder, NameHandler(tableName), e.Folder)
	err = util.CreateFileBytes(path, func(f *os.File) error {
		e := tmpl.Execute(f, data)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func MatchType(dbType string) string {
	if val, ok := DataTypeMap[dbType]; ok {
		return val
	} else {
		return "string"
	}
}

func NameHandler(name string) string {
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

func SetTmplGlobleVal(prjName, dbnm string, fmtDbName int) {
	projectName = prjName
	dbName = dbnm
	useDbName = (fmtDbName <= 0)
}
