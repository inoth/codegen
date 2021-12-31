package test

import (
	"codegen/db"
	"regexp"
	"testing"
	"time"
)

func TestGetTables(t *testing.T) {
	var (
		constr = ""
	)
	t1 := time.Now()

	db.InitDatabase(constr)
	dbNameRegexp := regexp.MustCompile(`/(\w+)\?`)
	dbName := dbNameRegexp.FindStringSubmatch(constr)[1]

	t.Logf("dbname: %v", dbName)

	tables := db.GetTables(dbName)
	if len(tables) <= 0 {
		t.Error("not found tables")
	}
	for _, table := range tables {
		t.Log(table.TableName, table.TableDesc)
	}
	t.Logf("ok; time: %v", time.Since(t1))
}
