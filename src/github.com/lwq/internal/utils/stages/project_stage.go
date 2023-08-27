package configs

import (
	. "lightweightpipline/internal/data/models"
)

type IProjectStage interface {
	//拉取
	Pull(branchName string) ([]OrderResponse, error)
	//构建
	Build() ([]OrderResponse, error)
	//推送
	Push() ([]OrderResponse, error)
	//发布
	Publish() ([]OrderResponse, error)
}
