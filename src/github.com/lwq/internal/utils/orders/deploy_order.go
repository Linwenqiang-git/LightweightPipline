package orders

import "fmt"

var nameSpace string = "scm-cloud-test"

func kubectlSetEnv(appName string, version int) string {
	return fmt.Sprintf("kubectl set env deployment/%s CURRENTVERSION=%d -n %s --kubeconfig /home/jenkins/.kube/config", appName, version, nameSpace)
}

func kubectlSetImage(appName string, version int) string {
	return fmt.Sprintf("kubectl set image deploy %s *=images.codefr.com/ufx/dotnet/%s:%d -n %s --kubeconfig /home/jenkins/.kube/config", appName, version, nameSpace)
}
