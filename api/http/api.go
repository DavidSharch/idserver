package api

import (
	"github.com/gin-gonic/gin"
)

type IdReqHttp struct {
	// Tag 业务标签
	Tag string `form:"tag"`
	// Step 步长
	Step int `form:"step"`
}

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

// GetIdByHttp 获取id，当id不存在时自动创建（可配置）
func GetIdByHttp(c *gin.Context) {
	var req IdReqHttp
	err := c.ShouldBind(&req)
	if err != nil {
		c.String(501, "error")
		return
	}
	if req.Tag == "" {
		c.String(501, "error")
		return
	}
	c.String(200, "i am a new id")
}

// DeleteIdByHttp 删除指定的tag，使用软删除
func DeleteIdByHttp(c *gin.Context) {
	var req IdReqHttp
	err := c.ShouldBind(&req)
	if err != nil {
		c.String(501, "error")
	}
	if req.Tag == "" {
		c.String(501, "error")
	}
	c.String(200, "pong")
}
