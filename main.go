package main

import (
	"codegen/db"
	"codegen/handle"
	"codegen/templet"
	"context"
	"flag"
	"os"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	MaxHandle = 100
)

var (
	dbConStr    string
	projectName string
	tableName   string
	fmtDbName   int
)

func init() {
	flag.StringVar(&dbConStr, "db", "", "database connect string 数据库连接串")
	flag.StringVar(&projectName, "name", "DefualtProject", "project name 项目名称")
	flag.StringVar(&tableName, "t", "", "tables; --t='tableA,tableB,tableC' 需要单独生成表名称，默认全表")
	flag.IntVar(&fmtDbName, "f", 0, "defualt:0；是否格式化处理数据命名格式，1：处理，0：不处理")
	flag.Parse()
}

func main() {
	log.Info("Initialization is in progress....")
	// 测试连接串
	if dbConStr == "" {
		log.Info("数据库连接串必传")
		os.Exit(1)
		return
	}

	db.InitDatabase(dbConStr)
	dbNameRegexp := regexp.MustCompile(`/(\w+)\?`)
	dbName := dbNameRegexp.FindStringSubmatch(dbConStr)[1]

	var tables []db.TableInfo
	if len(tableName) <= 0 {
		tables = db.GetTables(dbName)
	} else {
		// tables = strings.Split(tableName, ",")
		for _, tbl := range strings.Split(tableName, ",") {
			if len(tbl) <= 0 {
				continue
			}
			tables = append(tables, db.TableInfo{TableName: tbl})
		}
	}

	templet.SetTmplGlobleVal(projectName, dbName, fmtDbName)

	tn := len(tables)
	log.Infof("开始处理，一共%d个表", tn)

	if tn <= 0 {
		log.Warn("没有找到需要生成表")
		return
	}
	log.Info("预创建对应文件夹")
	handle.CreateFolder(projectName)

	t1 := time.Now()
	log.Info("开始执行表对应实体生成....")
	ch_progress := make(chan string, tn)

	for k, table := range tables {
		if k > 0 && k%MaxHandle == 0 {
			time.Sleep(time.Second * 1)
		}
		log.Infof("开始处理%v", table.TableName)
		go handle.HandlerTable(context.Background(), ch_progress, table)
	}

	curProgress := 0
	for curProgress < tn {
		r := <-ch_progress
		log.Info(r)
		curProgress++
	}
	log.Infof("finish; use: %s", time.Since(t1))
}
