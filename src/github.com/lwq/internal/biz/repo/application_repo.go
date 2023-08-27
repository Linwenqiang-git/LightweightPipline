package repo

import (
	. "lightweightpipline/internal/biz/aggregates"
)

type IApplicationRepo interface {
	GetApplicationInfo(appId int) Application
	GetApplicationBuildStage(appId int) []AppBuildStep
	GetBuildRecord(appId int, version int) []AppBuildRecord
}
