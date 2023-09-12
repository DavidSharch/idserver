package main

import (
	"github.com/gin-gonic/gin"
	api "github.com/sharch/idserver/api/http"
)

func main() {
	// 1. flag
	// 2. logger
	// 3. rpc和http服务
	// 4. 监听退出信号
	router := gin.Default()
	g1 := router.Group("/idserver")
	{
		// http://localhost:8089/idserver/id?tag=test
		g1.GET("/id", api.GetIdByHttp)
		// http://localhost:8089/idserver/ping
		g1.GET("/ping", api.Ping)
		g1.DELETE("/del", api.DeleteIdByHttp)
	}
	router.Run(":8089")
}
