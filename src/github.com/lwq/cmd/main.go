package main

import (
	build_app "lightweightpipline/api"
	provider "lightweightpipline/cmd/wire"
	. "lightweightpipline/configs/settings/log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	RoutingBinding(r)
	config, err := provider.GetConfigure()
	if err != nil {
		panic("read config err:" + err.Error())
	}
	NewLogOption().SetGlobalLogger()
	Logger.Info("Init Log")
	r.Run(config.System.Addr)
}

func RoutingBinding(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	build_app.BindingBuildAppService(g)
}
