package main

import (
	"flag"
	"fmt"
	api "github.com/sharch/idserver/api/http"
	"github.com/sharch/idserver/config"
	"github.com/sharch/idserver/internal/log"
	"github.com/sharch/idserver/internal/srv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	// 读取toml中的配置文件
	if err := config.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(config.Conf.Log)
	srv.NewService(config.Conf)

	// 启动http服务
	router := api.GetRouter()
	router.Run(":8089")
	// 启动rpc服务

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Printf("server start success pid:%d\n", os.Getpid())
	for s := range c {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.GetLogger().Info("exit")
			// 这里可以关闭其他内容
			return
		default:
			return
		}
	}
}
