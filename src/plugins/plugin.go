package plugin

type DbStructureGet interface {
	GetTables(dbName string) []Tables
	GetColumns(dbName, table string) []Columns
	TypeSwitch(dbtype string) string
}

type Tables struct {
	TableName string
	Desc      string
}

type Columns struct {
	Key    string
	Field  string
	Type   string
	Desc   string
	IsNull string
}
