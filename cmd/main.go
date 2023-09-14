package main

import (
	"flag"
	"fmt"
	"github.com/sharch/idserver/config"
	"github.com/sharch/idserver/internal/log"
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

	//srv.NewService(config.Conf)
	//
	//// 启动http服务
	//router := api.GetRouter()
	//router.Run(config.Conf.Server.Addr)

	//// 启动rpc服务
	//grpc.Init(configs.Conf, s)
	//// 设置etcd模式
	//if err := tool.InitMasterNode(configs.Conf.Etcd, configs.Conf.Server.Addr, 30); err != nil {
	//	panic(err)
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.GetLogger().Info(fmt.Sprintf("server start success pid:%d\n", os.Getpid()))
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
