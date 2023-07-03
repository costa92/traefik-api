package core

import (
	"time"

	"go.uber.org/dig"

	"treafik-api/config"
	config2 "treafik-api/pkg/common/config"
	logger2 "treafik-api/pkg/logger"
	"treafik-api/pkg/server"
	gin_server "treafik-api/pkg/server/gin-server"
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
	if err := container.Provide(NewDatabases); err != nil {
		return err
	}
	a.ctr = container
	return nil
}

// newConfig 初始化配置
func newConfig() *config.Config {
	var dotGraph bool
	var cfg config.Config
	initCfgLogger := initLogger()
	err := config2.New(
		config2.WithProviders(&config2.FileProvider{
			SkipIfPathEmpty: true,
		}),
		config2.WithRegisterFlags(func(flag *config2.FlagSet) {
			flag.BoolVar(&dotGraph, "graph", false, "parse the graph in Container into DOT format and writes it to stdout")
		}),
		config2.WithLogger(initCfgLogger),
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
	lc.Encoding = string(config2.LogEncodingConsole)
	return logger2.NewLogger(lc)
}

// Run 运行配置
func (a *App) Run() error {
	return a.ctr.Invoke(func(cfg *config.Config, dbs *Databases) error {
		// 实例化 AppServer
		httpServer := gin_server.NewAppServer(
			cfg,
			gin_server.Timeout(5*time.Second),
		)
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
