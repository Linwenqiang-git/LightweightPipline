package orders

import (
	provider "lightweightpipline/cmd/wire"
	. "lightweightpipline/internal/biz/aggregates"
	"strings"
	"sync"
)

var commandIdWithOrderFuncDic map[uint]string
var once sync.Once

func GetCommand(orderId uint, branchName string, image string) (string, error) {
	err := initDic()
	if err != nil {
		return "", err
	}
	//参数解析
	command := keyWordReplace(commandIdWithOrderFuncDic[orderId], branchName, image)
	return command, nil
}

func initDic() (err error) {
	once.Do(func() {
		commandIdWithOrderFuncDic = make(map[uint]string)
		dbContext, err := provider.GetDbContext()
		if err != nil {
			return
		}
		var commands []Command
		dbContext.GetDb().Find(&commands)
		if len(commands) == 0 {
			return
		}
		for i, v := range commands {
			commandIdWithOrderFuncDic[uint(i)+1] = v.Detail
		}
	})
	return
}

func keyWordReplace(command string, branchName string, image string) string {
	updatedCommand := strings.Replace(command, "{Branch}", branchName, -1)
	updatedCommand = strings.Replace(updatedCommand, "{Image}", image, -1)
	return updatedCommand
}
