package main

import (
	build_app "lightweightpipline/api"
	provider "lightweightpipline/cmd/wire"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	RoutingBinding(r)
	config, err := provider.GetConfigure()
	if err != nil {
		panic("read config err:" + err.Error())
	}
	r.Run(config.System.Addr) // 监听并在 0.0.0.0:8080 上启动服务
}

func RoutingBinding(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	build_app.BindingBuildAppService(g)
}
