package mysql_struct

import (
	"codegenfornet/components/db"
	plugin "codegenfornet/src/Plugin"
	"codegenfornet/src/Plugin/struct_crawl"
	"fmt"
)

var typeMap = map[string]string{
	"varchar":   "string",
	"decimal":   "decimal",
	"int":       "int",
	"tinyint":   "int",
	"datetime":  "DateTime",
	"timestamp": "DateTime",
}

type MysqlCrawl struct{}

func (mc MysqlCrawl) GetTables(dbName string) []plugin.Tables {
	var tables []plugin.Tables
	err := db.DB.Raw("SELECT TABLE_NAME as `TableName` FROM INFORMATION_SCHEMA.`TABLES` WHERE TABLE_SCHEMA = ?", dbName).Scan(&tables).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	return tables
}
func (mc MysqlCrawl) GetColumns(dbName, table string) []plugin.Columns {
	var cols []plugin.Columns
	db.DB.Raw("SELECT COLUMN_NAME as `Field`,DATA_TYPE as `Type`,COLUMN_COMMENT as `Desc`,IS_NULLABLE AS `IsNull`,COLUMN_KEY AS `Key` FROM INFORMATION_SCHEMA.`COLUMNS` WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", dbName, table).Scan(&cols)
	return cols
}

func (mc MysqlCrawl) TypeSwitch(dbtype string) string {
	t, ok := typeMap[dbtype]
	if !ok {
		return "string"
	}
	return t
}

func init() {
	struct_crawl.AddCrawl("mysql", &MysqlCrawl{})
}
