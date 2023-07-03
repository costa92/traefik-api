package gin_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"treafik-api/config"
	"treafik-api/pkg/logger"
	"treafik-api/pkg/server"
)

func GenericEngineServer(serverConfig *config.Server) *gin.Engine {
	gin.ForceConsoleColor()
	if serverConfig.Mode == "" {
		serverConfig.Mode = gin.DebugMode
	}
	gin.SetMode(serverConfig.Mode)
	r := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Infow("endpoint", "httpMethod", httpMethod, "absolutePath",
			absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}
	return r
}

type AppServer struct {
	GlobalConfig *config.Config
	Logger       logger.Logger
	Engine       *gin.Engine

	timeout      time.Duration
	secureServer *http.Server
}

func NewAppServer(cfg *config.Config, opts ...ServerOption) server.IAppServer {
	appSrv := &AppServer{
		GlobalConfig: cfg,
		Engine:       GenericEngineServer(&cfg.Server),
	}
	for _, o := range opts {
		o(appSrv)
	}
	appSrv.Run()
	return appSrv
}

func (s *AppServer) PreRun(router http.Handler) {
	serverConfig := s.GlobalConfig.Server
	s.secureServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", serverConfig.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *AppServer) Run() {
}

func (s *AppServer) Start(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorw("appService recover err", "err", err)
		}
	}()
	secureServer := s.secureServer
	logger.Infow("[HTTP] server started", "listen_addr", s.secureServer.Addr)
	// http 启动
	if err := secureServer.ListenAndServe(); err != nil {
		logger.Fatalf("secureServer ListenAndServe failed", "err", err)
		return err
	}
	return nil
}

func (s *AppServer) Stop(ctx context.Context) error {
	logger.Info("[HTTP] server stopping")
	return s.secureServer.Shutdown(ctx)
}
