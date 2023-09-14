package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	Addr       string
	User       string
	Password   string
	DbName     string
	Parameters string

	MaxConn      int
	IdleConn     int
	Debug        bool
	IdleTimeout  int
	QueryTimeout int //查询时间
	ExecTimeout  int //执行时间
}

func NewMysql(c *MysqlConfig) (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Addr, c.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	return
}
