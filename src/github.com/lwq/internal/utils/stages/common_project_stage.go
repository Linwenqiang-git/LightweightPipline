package configs

import (
	. "lightweightpipline/internal/data/models"
	. "lightweightpipline/internal/ipc"
)

// 项目构建阶段（通用）
type CommonProjectStage struct {
	appName             string
	projectPath         string
	currentBuildVersion string
	engine              Engine
}

func NewProjectStage(appName string, projectPath string) (IProjectStage, error) {
	version, err := getVersion()
	if err != nil {
		return nil, err
	}
	engine := NewEngine(projectPath)
	stage := &CommonProjectStage{
		appName:             appName,
		currentBuildVersion: version,
		projectPath:         projectPath,
		engine:              engine,
	}
	return stage, nil
}
func getVersion() (string, error) {
	return "1", nil
}

func (c *CommonProjectStage) Pull(branchName string) ([]OrderResponse, error) {
	pullOrders := []string{
		"git pull origin " + branchName,
	}
	return c.runOrder(pullOrders)
}

func (c *CommonProjectStage) Build() ([]OrderResponse, error) {
	buildOrders := []string{
		"docker build -t " + c.appName + ":" + c.currentBuildVersion + " .",
		"docker tag" + c.appName + ":" + c.currentBuildVersion + " images.codefr.com/ufx/dotnet/" + c.appName + ":" + c.currentBuildVersion,
	}
	return c.runOrder(buildOrders)
}

func (c *CommonProjectStage) Push() ([]OrderResponse, error) {
	pushOrders := []string{
		"docker login  --username ufdeveloper -p '123456' http://images.codefr.com",
		"docker push images.codefr.com/ufx/dotnet/" + c.appName + ":" + c.currentBuildVersion,
		"docker rmi " + c.appName + ":" + c.currentBuildVersion + " -f ",
		"docker rmi images.codefr.com/ufx/dotnet/" + c.appName + ":" + c.currentBuildVersion + " -f ",
	}
	return c.runOrder(pushOrders)
}

func (c *CommonProjectStage) Publish() ([]OrderResponse, error) {
	publishOrders := []string{
		"kubectl set env deployment/" + c.appName + " CURRENTVERSION=" + c.currentBuildVersion + " -n scm-cloud-test --kubeconfig /home/jenkins/.kube/config",
		"kubectl set image deploy " + c.appName + " *=images.codefr.com/ufx/dotnet/" + c.appName + ":" + c.currentBuildVersion + " -n scm-cloud-test --kubeconfig /home/jenkins/.kube/config",
	}
	return c.runOrder(publishOrders)
}

func (c *CommonProjectStage) runOrder(orders []string) (orderDetail []OrderResponse, err error) {
	for _, order := range orders {
		line_data, err := c.engine.RunOrder(order)
		if err != nil {
			orderDetail = append(orderDetail, OrderResponse{
				Status:  ErrorStatus,
				Order:   order,
				Message: err.Error(),
			})
			break
		}
		orderDetail = append(orderDetail, OrderResponse{
			Status:  SuccessStatus,
			Order:   order,
			Message: line_data,
		})
	}
	return orderDetail, err
}
