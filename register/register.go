package register

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	g_model  *GlobalRegister
	once     sync.Once
	initOnce sync.Once
)

type IRegister interface {
	Init() error
}
type IServeStart interface {
	ServeStart() error
}

type GlobalRegister struct {
	servers []IRegister
}

func instance() *GlobalRegister {
	once.Do(func() {
		g_model = &GlobalRegister{}
	})
	return g_model
}

// 注册组件
func Register(models ...IRegister) *GlobalRegister {
	if len(models) <= 0 {
		fmt.Printf("%v\n", errors.New("No services have been loaded yet."))
		os.Exit(1)
	}
	model := instance()
	model.servers = make([]IRegister, len(models))
	for i, m := range models {
		model.servers[i] = m
	}
	return model
}

// 追加注册
func (g *GlobalRegister) AddRegister(models ...IRegister) *GlobalRegister {
	g.servers = append(g.servers, models...)
	return g
}

// 根据注册顺序，配置时注意引用优先级，初始化组件模块
func (g *GlobalRegister) Init() *GlobalRegister {
	initOnce.Do(func() {
		for _, svc := range g.servers {
			must(svc.Init())
		}
	})
	return g
}

// 运行子服务，比如性能分析，或者websocket之类的
func (g *GlobalRegister) SubServe(serve ...IServeStart) *GlobalRegister {
	for _, subSvc := range serve {
		go func(svc IServeStart) {
			defer func() {
				// TODO:协程内单独的异常捕获if exception := recover(); exception != nil {
				if exception := recover(); exception != nil {
					if err, ok := exception.(error); ok {
						fmt.Printf("%v\n", err)
					} else {
						panic(exception)
					}
					os.Exit(1)
				}
			}()
			must(svc.ServeStart())
		}(subSvc)
	}
	return g
}

// 运行服务
func (g *GlobalRegister) Run(serve IServeStart) error {
	return serve.ServeStart()
}

func must(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
