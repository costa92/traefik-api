package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"treafik-api/config"
	logger2 "treafik-api/pkg/logger"
	"treafik-api/pkg/middlewares"
	"treafik-api/routers"
)

func GenericEngineServer() *gin.Engine {
	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger2.Infow("endpoint", "httpMethod", httpMethod, "absolutePath",
			absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}
	return r
}

type IAppService interface {
	InstallLogger(logger logger2.Logger)
	InstallDatabase(dbs *Databases)
	InstallRoutes()
	installMiddlewares()
	PreRun()
	InitRun()
	Start()
}

type AppService struct {
	GlobalConfig *config.Config
	Logger       logger2.Logger
	*gin.Engine
	secureServer *http.Server
	Databases    *Databases
}

func NewAppService(cfg *config.Config) IAppService {
	return &AppService{
		GlobalConfig: cfg,
		Engine:       GenericEngineServer(),
	}
}

func (a *AppService) InstallDatabase(dbs *Databases) {
	a.Databases = dbs
}

func (a *AppService) InstallLogger(logger logger2.Logger) {
	a.Logger = logger
}

// PreRun 预运行
func (a *AppService) PreRun() {
	serverConfig := a.GlobalConfig.Server
	router := a.Engine
	a.secureServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", serverConfig.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger2.Infof("server running on port %s", serverConfig.Port)
}

// InitRun 进入进程
func (a *AppService) InitRun() {
	// 初始中间
	a.installMiddlewares()
	// 后置跟服务
	a.PreRun()
}

func (a *AppService) Start() {
	defer func() {
		if err := recover(); err != nil {
			logger2.Errorw("appService recover err", "err", err)
		}
	}()
	a.InitRun()
	server := a.secureServer
	// http 启动
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger2.Fatal(err)
		}
	}()
	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit
	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		logger2.Fatal("Server Shutdown:", err)
	}
}

// InstallMiddlewares 安装中间件
func (a *AppService) installMiddlewares() {
	serverConfig := a.GlobalConfig.Server
	s := a.Engine
	for _, m := range serverConfig.Middlewares {
		mw, ok := middlewares.Middlewares[m]
		if !ok {
			logger2.Errorf("can not find middleware: %s", m)
			continue
		}
		logger2.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// InstallRoutes 安装路由
func (a *AppService) InstallRoutes() {
	// apiRoute := routers.NewServerApiRoute(a.MysqlDb)
	routers.InitApiRouter(a.Engine)
}
