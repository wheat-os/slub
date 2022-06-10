package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"douban/douban"
	"github.com/spf13/viper"
	"github.com/wheat-os/wlog"
)

var projectConfFile = "/home/bandl/go/src/github.com/wheat-os/slub/example/douban/conf.toml"

func signalClose(cannel context.CancelFunc) {

	sig := make(chan os.Signal, 1)
	// 监听 退出
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	wlog.Info("the listener hears the exit signal, exiting, please do not kill the process directly.")
	cannel()
}

func main() {
	viper.SetConfigFile(projectConfFile)

	engine := douban.DefaultEngine
	ctx, cannel := context.WithCancel(context.Background())
	go signalClose(cannel)

	engine.Start(ctx)

	engine.Close()
}
