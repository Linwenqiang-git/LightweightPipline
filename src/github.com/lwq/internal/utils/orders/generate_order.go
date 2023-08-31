package orders

import (
	"strings"
)

// 解析命令
func ParseCommand(command, branchName, image string) string {
	updatedCommand := strings.Replace(command, "{Branch}", branchName, -1)
	updatedCommand = strings.Replace(updatedCommand, "{Image}", image, -1)
	updatedCommand = strings.Replace(updatedCommand, "{NameSpace}", image, -1)
	updatedCommand = strings.Replace(updatedCommand, "{ServiceName}", image, -1)
	return updatedCommand
}
