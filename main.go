package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"treafik-api/core"
	"treafik-api/pkg/common/version"
	"treafik-api/pkg/logger"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	version.MustRegisterVersionCollector()
	theApp, err := core.NewApp()
	if err != nil {
		logger.Errorw("app start failed", "err", err)
		theApp.Shutdown()
		return
	}
	// 信号清除
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err = theApp.Run(); err != nil {
		logger.Errorw("app exited with error", "err", err)
		chSig <- syscall.SIGTERM
	}
	<-chSig

	theApp.Shutdown()
}
