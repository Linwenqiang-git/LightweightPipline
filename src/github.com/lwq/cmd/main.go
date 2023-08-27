package main

import (
	build_app "lightweightpipline/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	RoutingBinding(r)
	r.Run("127.0.0.1:8090") // 监听并在 0.0.0.0:8080 上启动服务
}

func RoutingBinding(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	build_app.BindingBuildAppService(g)
}
