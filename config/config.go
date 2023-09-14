package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sharch/idserver/internal/log"
	"github.com/sharch/idserver/internal/orm"
)

var (
	// confPath 在init中初始化
	confPath string
	// Conf 定义一个全局变量
	Conf = new(Config)
)

type Config struct {
	Etcd   []string
	Log    *log.Options
	Mysql  *orm.MysqlConfig
	Server *ServerConfig
}

type ServerConfig struct {
	Addr string
}

func init() {
	flag.StringVar(&confPath, "conf", "./cmd/config.toml", "default config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	fmt.Printf("读取的配置如下:\n")
	fmt.Printf("etcd: %+v \n", Conf.Etcd)
	fmt.Printf("log: %+v \n", Conf.Log)
	fmt.Printf("mysql: %+v \n", Conf.Mysql)
	fmt.Printf("server: %+v \n", Conf.Server)
	return
}
