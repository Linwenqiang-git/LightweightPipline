package repo

import (
	. "lightweightpipline/internal/biz/aggregates"
)

type ICommandRepo interface {
	GetAllCommands(commandIds []interface{}) (map[uint]Command, error)
}
