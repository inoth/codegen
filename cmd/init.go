/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"codegenfornet/components/cache"
	"codegenfornet/components/db"
	"codegenfornet/register"
	"codegenfornet/src/codegen"
	"codegenfornet/util"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "初始化一个项目",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		projectName = "DefaultProject"
		if len(args) > 0 {
			projectName = args[0]
		}
		reg := register.Register(
			&cache.CacheComponents{},
		)
		fordb := len(dbhost) > 0 && len(dbuser) > 0 && len(dbpasswd) > 0 && len(databaseName) > 0 && len(dbtype) > 0
		if fordb {
			switch dbtype {
			case "mysql":
				reg = reg.AddRegister(&db.MysqlConnect{
					User:   dbuser,
					Passwd: dbpasswd,
					Host:   dbhost,
					DbName: databaseName,
				})
			default:
				panic("该数据库尚未支持")
			}
		}

		util.Must(reg.Init().Run(&codegen.CodegenForNet{
			DbType:      dbtype,
			ProjectName: projectName,
			DbName:      databaseName,
			TempPath:    tempPath,
			SourceDb:    fordb,
		}))
	},
}

var (
	projectName  string
	databaseName string
	dbhost       string
	dbuser       string
	dbpasswd     string
	dbtype       string
	tempPath     string
)

func init() {
	initCmd.Flags().StringVar(&dbtype, "db", "mysql", "数据库类型, 目前仅支持mysql")
	initCmd.Flags().StringVar(&dbhost, "host", "", "数据库地址")
	initCmd.Flags().StringVar(&databaseName, "dbname", "", "数据库名称")
	initCmd.Flags().StringVar(&dbuser, "user", "", "数据库用户名")
	initCmd.Flags().StringVar(&dbpasswd, "passwd", "", "数据库密码")

	// 模板存放文件夹
	initCmd.Flags().StringVar(&tempPath, "temp", "temp/temp.zip", "模板文件夹")
	rootCmd.AddCommand(initCmd)
}
