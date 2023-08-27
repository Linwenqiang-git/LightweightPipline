package services

import (
	"context"
	"fmt"
	provider "lightweightpipline/cmd/wire"
	. "lightweightpipline/configs/settings/db"
	. "lightweightpipline/internal/biz/aggregates"
	. "lightweightpipline/internal/biz/repo"
	buildProgress "lightweightpipline/internal/data/consts"
	buildResult "lightweightpipline/internal/data/consts"
	. "lightweightpipline/internal/ipc"
	. "lightweightpipline/internal/utils/orders"
	"strconv"
	"time"

	"github.com/ahmetb/go-linq/v3"
)

type AppBuildService struct {
	dbContext *DbContext
}

func GetAppBuildService() (*AppBuildService, error) {
	dbContext, err := provider.GetDbContext()
	if err != nil {
		return nil, err
	}
	service := &AppBuildService{
		dbContext: dbContext,
	}
	return service, err
}

func (b *AppBuildService) AppBuild(appName, branchName string) error {
	var appInfo Application
	b.dbContext.GetDb().Where("name=?", appName).First(&appInfo)
	if appInfo.ID == 0 {
		return fmt.Errorf("获取应用信息失败：%s", appName)
	}

	//获取应用构建步骤
	var buildStepInfo []AppBuildStep
	b.dbContext.GetDb().Where("app_id=? AND is_open=?", appInfo.ID, true).Order("sort asc").Find(&buildStepInfo)
	if len(buildStepInfo) == 0 {
		return fmt.Errorf("应用%s未设置可用构建步骤", appName)
	}

	//获取构建步骤对应的命令
	var commandInfo []StageCommand
	stageIds := linq.From(buildStepInfo).
		SelectT(func(step AppBuildStep) uint {
			return step.StageId
		}).Results()
	b.dbContext.GetDb().Where("stage_id IN ?", stageIds).Order("sort asc").Find(&commandInfo)
	if len(commandInfo) == 0 {
		return fmt.Errorf("应用%s未设置可用构建命令", appName)
	}
	imageName, err := b.getImageName(appName)
	if err != nil {
		return err
	}
	appBuildRecord, buildRecordDetail, err := b.addBuildRecore(appInfo.ID, imageName, branchName, buildStepInfo, commandInfo)
	if err != nil {
		return err
	}
	go b.excuteCommand(appInfo.Path, appBuildRecord, buildRecordDetail)
	return nil
}

// 添加执行记录
func (b *AppBuildService) addBuildRecore(appId uint, imageName string, branch string, buildStepInfo []AppBuildStep, commandInfo []StageCommand) (*AppBuildRecord, []AppBuildRecordDetail, error) {
	//生成构建记录主表
	appBuildRecord := AppBuildRecord{
		AppId:     appId,
		ImageName: imageName,
		Branch:    branch,
		Progress:  buildProgress.Start,
		Result:    buildResult.InProgress,
	}
	//开启事务
	tx := b.dbContext.GetDb().Begin()
	result := tx.Create(&appBuildRecord)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	//生成构建记录明细表
	buildRecordId := appBuildRecord.ID
	var buildRecordDetail []AppBuildRecordDetail
	for _, stage := range buildStepInfo {
		for _, command := range commandInfo {
			if command.StageId == stage.ID {
				detail := AppBuildRecordDetail{
					BuildRecordId: buildRecordId,
					StageId:       stage.ID,
					StageName:     stage.StageName,
					CommandId:     command.ID,
					Command:       "",
					Output:        "",
					Result:        buildResult.Start,
				}
				commandStr, err := GetCommand(command.CommandId, branch, imageName)
				if err != nil {
					tx.Rollback()
					return nil, nil, err
				}
				detail.Command = commandStr
				buildRecordDetail = append(buildRecordDetail, detail)
			}
		}
	}
	result = tx.Create(&buildRecordDetail)
	if result.Error != nil {
		tx.Rollback()
		return nil, nil, result.Error
	}
	tx.Commit()
	return &appBuildRecord, buildRecordDetail, nil
}

// 获取镜像名称
func (b *AppBuildService) getImageName(appName string) (string, error) {
	buildVersion, err := b.getBuildVersion(appName)
	if err != nil {
		return "", fmt.Errorf("获取构建版本失败：%s", err.Error())
	}
	date := time.Now().Format("2006-01-02")
	return appName + ":" + date + "-" + strconv.Itoa(int(buildVersion)), nil
}

// 获取构建版本
func (b *AppBuildService) getBuildVersion(appName string) (uint, error) {

	ctx := context.Background()
	resdisOption := NewRedisOption()
	date := time.Now().Format("2006-01-02")
	key := resdisOption.RedisPrefix + ":" + appName + date
	redis := resdisOption.CreateConnect()
	value, err := redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取版本号失败：%s", err.Error())
	}
	return uint(value), nil
}

// 执行命令
func (b *AppBuildService) excuteCommand(projectPath string, buildRecore *AppBuildRecord, excuteCommands []AppBuildRecordDetail) {
	// defer func() {
	// 	fileManage.DeleteFilesInDirectory(projectPath)
	// }()
	engine := NewEngine(projectPath)
	isSuccess := true
	for _, command := range excuteCommands {
		//执行命令
		result, err := engine.RunCommand(command.Command)
		if err != nil {
			command.Output = err.Error()
			command.Result = buildResult.Fail
		} else {
			command.Output = result
			command.Result = buildResult.Success
		}
		b.dbContext.GetDb().Save(command)
		if command.Result == buildResult.Fail {
			isSuccess = false
			break
		}
	}
	if isSuccess {
		buildRecore.Result = buildResult.Success
	} else {
		buildRecore.Result = buildResult.Fail
	}
	buildRecore.Progress = buildProgress.End
	b.dbContext.GetDb().Save(buildRecore)
}
