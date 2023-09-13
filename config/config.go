package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sharch/idserver/internal/log"
	"github.com/sharch/idserver/internal/orm"
)

var (
	confPath string
	Conf     = new(Config)
)

type Config struct {
	Etcd   []string
	Log    *log.Options
	Mysql  *orm.Config
	Server *Srv
}

type Srv struct {
	Ip   string
	Port int
}

func init() {
	flag.StringVar(&confPath, "conf", "./config.toml", "default config path")
}

func Init() error {
	return local()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
