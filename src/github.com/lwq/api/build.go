package api

import (
	service "lightweightpipline/internal/biz/services"
	. "lightweightpipline/internal/data/models"

	"github.com/gin-gonic/gin"
)

// 构建大货前服务
func BuildDevelopeCenter(c *gin.Context) {
	appbuildService, err := service.GetAppBuildService()
	if err != nil {
		c.JSON(200, Fail(200, err.Error()))
	}
	err = appbuildService.AppBuild("ufx-scm-cloud-developecenter", "stage")
	var result Response
	if err != nil {
		result = Fail(200, err.Error())
	} else {
		result = SuccessWithMessage("任务：ufx-scm-cloud-developecenter 已成功构建", nil)
	}
	c.JSON(200, result)
}

// 构建订单服务
func BuildProductOrderCenter(c *gin.Context) {
	appbuildService, err := service.GetAppBuildService()
	if err != nil {
		c.JSON(200, Fail(200, err.Error()))
	}
	go appbuildService.AppBuild("ufx-scm-cloud-productordercenter", "stage")
	result := SuccessWithMessage("任务：ufx-scm-cloud-productordercenter 已成功构建", nil)
	c.JSON(200, result)
}

func BindingBuildAppService(g *gin.Engine) {
	engineApi := g.Group("/api/build")
	{
		engineApi.GET("/developcenter", BuildDevelopeCenter)
		engineApi.GET("/productordercenter", BuildProductOrderCenter)
	}
}
