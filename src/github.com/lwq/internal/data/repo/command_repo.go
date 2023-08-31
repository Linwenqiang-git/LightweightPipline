package repo

import (
	"fmt"
	. "lightweightpipline/internal/biz/aggregates"
	. "lightweightpipline/internal/biz/repo"
)

type CommandRepo struct {
	context *DbContext
}

func NewCommandRepo(dbContext *DbContext) (ICommandRepo, error) {
	repo := &CommandRepo{
		context: dbContext,
	}
	return repo, nil
}

// 获取全部命令
func (o *CommandRepo) GetAllCommands(commandIds []interface{}) (map[uint]Command, error) {
	var commandInfo []Command
	tx := o.context.GetDb()
	if len(commandIds) > 0 {
		tx = tx.Where("id IN ?", commandIds)
	}
	tx.Find(&commandInfo)
	if len(commandInfo) == 0 {
		return nil, fmt.Errorf("未获取到相关命令信息")
	}
	commandDict := make(map[uint]Command)
	for _, command := range commandInfo {
		commandDict[command.ID] = command
	}
	return commandDict, nil
}
