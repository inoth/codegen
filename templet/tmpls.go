package templet

var (
	Entity = `namespace {{.ProjectName}}.Entity
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
	}`
	Model = `namespace {{.ProjectName}}.Model
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
	}`
	Mapper = `namespace {{.ProjectName}}.AutoMapperProfile
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
	}`
	IRepository = `namespace {{.ProjectName}}.IRepository
	{
		public interface I{{.TableName}}Repository
		{
		}
	}`
	Repository = `namespace {{.ProjectName}}.Repository
	{
		public class {{.TableName}}Repository: BaseRepository<{{if .UseDbName}}{{.DbTableName}}{{else}}{{.TableName}}{{end}}>, I{{.TableName}}Repository
		{
			public {{.TableName}}Repository(){}
		}
	}`
	IService = `namespace {{.ProjectName}}.IService
	{
		public interface I{{.TableName}}Service
		{
		}
	}`
	Service = `namespace {{.ProjectName}}.Service
	{
		public class {{.TableName}}Service: BaseService, I{{.TableName}}Service
		{
			private readonly I{{.TableName}}Repository _repository;
			public {{.TableName}}Service(I{{.TableName}}Repository repository){
				_repository = repository;
			}
		}
	}`
	Controller = `namespace {{.ProjectName}}.Controllers
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
	}`
)
