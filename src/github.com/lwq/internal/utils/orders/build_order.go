package orders

import "fmt"

//镜像构建
func buildImage(appName string, version int) string {
	return fmt.Sprintf("docker build -t %s:%d .", appName, version)
}

//镜像标签
func imageTag(appName string, version int) string {
	return fmt.Sprintf("docker tag %s:%d images.codefr.com/ufx/dotnet/%s:%d", appName, version, appName, version)
}
