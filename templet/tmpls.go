package templet

var (
	Model = "package model\n" +
		"type {{.TableName}} struct {\n" +
		"{{range .Fields}}\n" +
		"{{.Field}} {{.DbType}} `gorm:\"{{.DbField}}\"`" +
		"{{end}}\n}"

	// {
	// 	public class {{if .UseDbName}}{{.DbTableName}}{{else}}{{.TableName}}{{end}}
	// 	{
	// 		{{if .UseDbName}}
	// 			{{range .Fields}}
	// 		/// <summary>
	// 		///	{{.Desc}}
	// 		/// </summary>
	// 		public {{.DbType}} {{.DbField}} {get; set;}
	// 			{{end}}
	// 		{{ else }}
	// 			{{range .Fields}}
	// 		/// <summary>
	// 		///	{{.Desc}}
	// 		/// </summary>
	// 		public {{.DbType}} {{.Field}} {get; set;}
	// 			{{end}}
	// 		{{end}}
	// 	}
	// }`
)
