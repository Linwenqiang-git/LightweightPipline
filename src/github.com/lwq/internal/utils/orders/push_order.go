package orders

import "fmt"

//登录镜像仓库
func loginImageRepo(username, password, addr string) string {
	username = "ufdeveloper"
	password = "123456"
	addr = "http://images.codefr.com"
	return fmt.Sprintf("docker login  --username %s -p '%s' %s", username, password, addr)
}

//镜像推送
func pushImage(appName string, version int) string {
	return fmt.Sprintf("docker push images.codefr.com/ufx/dotnet/%s:%d", appName, version)
}

//删除本地镜像
func DeleteLocalImage(appName string, version int) string {
	return fmt.Sprintf("docker rmi %s:%d -f ", appName, version)
}
