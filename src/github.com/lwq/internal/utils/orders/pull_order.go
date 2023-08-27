package orders

import "fmt"

//代码拉取
func pullCode(args interface{}) (string, error) {
	if branchName, ok := args.(string); ok {
		return fmt.Sprintf("git pull origin %s", branchName), nil
	} else {
		return "", fmt.Errorf("parameters error")
	}
}
