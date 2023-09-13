package api

import "github.com/gin-gonic/gin"

func GetRouter() *gin.Engine {
	router := gin.Default()
	g1 := router.Group("/idserver")
	{
		// http://localhost:8089/idserver/id?tag=test
		g1.GET("/id", GetIdByHttp)
		// http://localhost:8089/idserver/ping
		g1.GET("/ping", Ping)
		g1.DELETE("/del", DeleteIdByHttp)
	}
	return router
}
