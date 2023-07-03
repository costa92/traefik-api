package core

import (
	"time"

	"go.uber.org/dig"

	"treafik-api/config"
	"treafik-api/controller"
	"treafik-api/db"
	commonConfig "treafik-api/pkg/common/config"
	logger2 "treafik-api/pkg/logger"
	"treafik-api/pkg/server"
	ginServer "treafik-api/pkg/server/gin-server"
	"treafik-api/routers"
)

type App struct {
	ctr *dig.Container
}

func NewApp() (*App, error) {
	app := &App{}
	err := app.buildContainer()
	if err != nil {
		return nil, err
	}
	return app, nil
}

// 容器
func (a *App) buildContainer() error {
	container := dig.New()
	if err := container.Provide(newConfig); err != nil {
		return err
	}
	if err := container.Provide(db.NewDatabases); err != nil {
		return err
	}
	if err := container.Provide(initApi); err != nil {
		return err
	}
	a.ctr = container
	return nil
}

func initApi(cfg *config.Config, dbs *db.Databases) (*controller.ApiHttp, error) {
	apiHttp := controller.NewApiHttp(cfg, dbs)
	return apiHttp, nil
}

// newConfig 初始化配置
func newConfig() *config.Config {
	var dotGraph bool
	var cfg config.Config
	initCfgLogger := initLogger()
	err := commonConfig.New(
		commonConfig.WithProviders(&commonConfig.FileProvider{
			SkipIfPathEmpty: true,
		}),
		commonConfig.WithRegisterFlags(func(flag *commonConfig.FlagSet) {
			flag.BoolVar(&dotGraph, "graph", false, "parse the graph in Container into DOT format and writes it to stdout")
		}),
		commonConfig.WithLogger(initCfgLogger),
	).Load(&cfg)
	if err != nil {
		panic(err)
	}
	// init app logger
	logger2.SetConfig(logger2.NewConfigFromInterface(&cfg.Log))
	return &cfg
}

func initLogger() logger2.Logger {
	lc := logger2.NewDefaultConfig()
	lc.OutputPaths = []string{"stderr"}
	lc.DisableStacktrace = true
	lc.EnableColor = true
	lc.Encoding = string(commonConfig.LogEncodingConsole)
	return logger2.NewLogger(lc)
}

// Run 运行配置
func (a *App) Run() error {
	return a.ctr.Invoke(func(cfg *config.Config, dbs *db.Databases, api *controller.ApiHttp) error {
		// 实例化 AppServer
		httpServer := ginServer.NewAppServer(
			ginServer.WithServiceConfig(&cfg.Server),
			ginServer.WithMiddleware(cfg.Server.Middlewares),
			ginServer.Timeout(5*time.Second),
		)
		// PreRun
		engine := httpServer.(*ginServer.AppServer).Engine
		routers.RegisterRouter(engine, api)
		httpServer.PreRun(engine)
		// 实例化 App
		app := server.New(
			server.Server(httpServer),
		)
		// 启动
		return app.Run()
	})
}

func (a *App) Shutdown() {
}
