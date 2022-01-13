# codegen
```shell
PS C:> codegen.exe --help
Usage of C:codegen.exe:
  -db string
        database connect string 数据库连接串
  -f int
        defualt:0；是否格式化处理数据命名格式，1：处理，0：不处理
  -name string
        project name 项目名称 (default "DefualtProject")
  -t string
        tables; --t='tableA,tableB,tableC' 需要单独生成表名称，默认全表
```
.\codegen.exe --db='user:passwd@(host:port)/databse?charset=utf8&parseTime=True&loc=Local' --name=projectName --t='tableA'

.\codegen.exe --db='user:passwd@(host:port)/databse?charset=utf8&parseTime=True&loc=Local' --name=projectName

.\codegen.exe --db='user:passwd@(host:port)/databse?charset=utf8&parseTime=True&loc=Local'