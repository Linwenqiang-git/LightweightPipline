package api

import (
	"fmt"
	service "lightweightpipline/internal/biz/services"
	. "lightweightpipline/internal/data/models"

	"github.com/gin-gonic/gin"
)

// 构建大货前服务
func buildDevelopeCenter(c *gin.Context) {
	buildProject(c, "ufx-scm-cloud-developecenter", "stage")
}

// 构建订单服务
func buildProductOrderCenter(c *gin.Context) {
	buildProject(c, "ufx-scm-cloud-productordercenter", "stage")
}

// 构建全部应用
func buildAllService(c *gin.Context) {
	appbuildService, err := service.GetAppBuildService()
	if err != nil {
		c.JSON(200, Fail(200, err.Error()))
	}
	err = appbuildService.AllAppBuild()
	var result Response
	if err != nil {
		result = Fail(500, err.Error())
	} else {
		result = SuccessWithMessage("任务已全部成功构建", nil)
	}
	c.JSON(200, result)
}

// 绑定构建服务
func BindingBuildAppService(g *gin.Engine) {
	engineApi := g.Group("/api/build")
	{
		engineApi.GET("/developcenter", buildDevelopeCenter)
		engineApi.GET("/productordercenter", buildProductOrderCenter)
		engineApi.GET("/allservice", buildAllService)
	}
}

// 构建项目
func buildProject(c *gin.Context, appName, branchName string) {
	appbuildService, err := service.GetAppBuildService()
	if err != nil {
		c.JSON(200, Fail(200, err.Error()))
	}
	err = appbuildService.AppBuild(appName, branchName)
	var result Response
	if err != nil {
		result = Fail(500, err.Error())
	} else {
		result = SuccessWithMessage(fmt.Sprintf("任务：%s 已成功构建", appName), nil)
	}
	c.JSON(200, result)
}
