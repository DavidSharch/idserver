package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sharch/idserver/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGorm(t *testing.T) {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := 3306
	dbname := "idserver"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	var s entity.Segments
	db.Where("id = ?", 1).First(&s)

	router := gin.Default()
	g2 := router.Group("/db")
	{
		g2.GET("/test", func(c *gin.Context) {
			c.JSON(200, s)
		})
	}
	router.Run(":8089")
}
