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
		Tmpl: `
        namespace {{.ProjectName}}.Entity
		{
			public class {{if .UseDbName}}{{.DbTableName}}{{else}}{{.TableName}}{{end}}
			{
				{{if .UseDbName}}
					{{range .Fields}}
				/// <summary>
				///	{{.Desc}}
				/// </summary>
				public {{.DbType}} {{.DbField}} {get; set;}
					{{end}}
				{{ else }}
					{{range .Fields}}
				/// <summary>
				///	{{.Desc}}
				/// </summary>
				public {{.DbType}} {{.Field}} {get; set;}
					{{end}}
				{{end}}
			}
		}`}

	FileTemplate["Model"] = &EntityHandle{
		Folder: "Model",
		Tmpl: `
        namespace {{.ProjectName}}.Model
		{
			public class {{.TableName}}VM
			{
				{{range .Fields}}
				/// <summary>
				///	{{.Desc}}
				/// </summary>
				public {{.DbType}} {{.Field}} {get; set;}
				{{end}}
			}
		}`}

	FileTemplate["AutoMapper"] = &EntityHandle{
		Folder: "AutoMapper",
		Tmpl: `
		namespace {{.ProjectName}}.AutoMapperProfile
		{
			public class {{.TableName}}Profile : Profile
			{
				public {{.TableName}}Profile()
				{
					CreateEntityMaps();
				}
		
				private void CreateEntityMaps()
				{
					{{if .UseDbName}}					
					CreateMap<{{.DbTableName}}, {{.TableName}}VM>()
					{{range .Fields}}
					.ForMember(e => e.{{.Field}}, opt => opt.MapFrom(x => x.{{.DbField}}))
					{{end}};			
					{{else}}					
					CreateMap<{{.TableName}}, {{.TableName}}VM>().ReverseMap();
					{{end}}
				}
			}
		}`}

	FileTemplate["IRepository"] = &NormalHandle{
		Folder: "IRepository",
		Tmpl: `
		namespace {{.ProjectName}}.IRepository
		{
			public interface I{{.TableName}}Repository
			{
			}
		}`}

	FileTemplate["Repository"] = &NormalHandle{
		Folder: "Repository",
		Tmpl: `
		namespace {{.ProjectName}}.Repository
		{
			public class {{.TableName}}Repository: BaseRepository<{{if .UseDbName}}{{.DbTableName}}{{else}}{{.TableName}}{{end}}>, I{{.TableName}}Repository
			{
				public {{.TableName}}Repository(){}
			}
		}`}

	FileTemplate["IService"] = &NormalHandle{
		Folder: "IService",
		Tmpl: `
        namespace {{.ProjectName}}.IService
        {
            public interface I{{.TableName}}Service
            {
            }
        }`}

	FileTemplate["Service"] = &NormalHandle{
		Folder: "Service",
		Tmpl: `
        namespace {{.ProjectName}}.Service
        {
            public class {{.TableName}}Service: BaseService, I{{.TableName}}Service
            {
				private readonly I{{.TableName}}Repository _repository;
                public {{.TableName}}Service(I{{.TableName}}Repository repository){
					_repository = repository;
				}
            }
        }`}

	FileTemplate["Controller"] = &NormalHandle{
		Folder: "Controller",
		Tmpl: `
        namespace {{.ProjectName}}.Controllers
        {
            [Route("api/[controller]/[Action]")]
            [ApiController]
            public class {{.TableName}}Controller : ControllerBase
            {
                private readonly I{{.TableName}}Service _service;
                public {{.TableName}}Controller(I{{.TableName}}Service service)
                {
                    _service = service;
                }
            }
        }`}
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
		data.Fields[i] = Field{DbType: col.DataType,
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
