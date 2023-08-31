package services

import (
	"context"
	"encoding/json"
	"fmt"
	provider "lightweightpipline/cmd/wire"
	. "lightweightpipline/configs/settings/db"
	. "lightweightpipline/internal/biz/aggregates"
	. "lightweightpipline/internal/biz/repo"
	buildProgress "lightweightpipline/internal/data/consts"
	buildResult "lightweightpipline/internal/data/consts"
	. "lightweightpipline/internal/ipc"
	. "lightweightpipline/internal/utils/orders"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/log"
)

type CommandResultDto struct {
	Result          string
	Duration        float64
	Error           error
	IsExcuteSuccess bool
}

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

// 应用构建
func (b *AppBuildService) AppBuild(appName, branchName string) error {
	//获取执行应用的信息
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
	var stageCommandInfo []StageCommand
	stageIds := linq.From(buildStepInfo).
		SelectT(func(step AppBuildStep) uint {
			return step.StageId
		}).Results()
	b.dbContext.GetDb().Where("stage_id IN ?", stageIds).Order("sort asc").Find(&stageCommandInfo)
	if len(stageCommandInfo) == 0 {
		return fmt.Errorf("应用%s未设置可用构建命令", appName)
	}
	//获取镜像名称
	imageName, err := b.getImageName(appName)
	if err != nil {
		return err
	}
	//获取构建包含的命令信息
	commandDic, err := b.getCommandDict(stageCommandInfo)
	if err != nil {
		return err
	}

	//添加构建记录明细
	appBuildRecord, buildRecordDetail, err := b.addBuildRecore(appInfo.ID, imageName, branchName, buildStepInfo, stageCommandInfo, commandDic)
	if err != nil {
		return err
	}

	//执行命令
	go b.excuteCommand(appInfo.Path, appBuildRecord, buildRecordDetail, commandDic)
	log.Info(fmt.Sprintf("start excute %s commands", appName))
	return nil
}

// 添加构建记录
func (b *AppBuildService) addBuildRecore(appId uint, imageName string, branch string, buildStepInfo []AppBuildStep, stageCommandInfo []StageCommand, commandDict map[uint]Command) (*AppBuildRecord, []AppBuildRecordDetail, error) {
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
		for _, command := range stageCommandInfo {
			if command.StageId == stage.StageId {
				detail := AppBuildRecordDetail{
					BuildRecordId: buildRecordId,
					StageId:       stage.ID,
					StageName:     stage.StageName,
					CommandId:     command.CommandId,
					Command:       "",
					SecondCommand: "",
					Output:        "",
					Result:        buildResult.Start,
				}
				if commandInfo, ok := commandDict[command.CommandId]; ok {
					detail.Command = ParseCommand(commandInfo.Detail, branch, imageName)
					detail.SecondCommand = ParseCommand(commandInfo.SecondCommand, branch, imageName)
					buildRecordDetail = append(buildRecordDetail, detail)
				} else {
					tx.Rollback()
					return nil, nil, fmt.Errorf(fmt.Sprintf("未找到相关命令信息，id:%d", command.CommandId))
				}
			}
		}
	}
	if len(buildRecordDetail) == 0 {
		tx.Rollback()
		return nil, nil, fmt.Errorf("未生成构建记录明细")
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
	if value == 1 {
		expiration := 86410 //24小时零10秒
		//设置key过期时间
		_, err = redis.Expire(ctx, key, time.Duration(expiration)*time.Second).Result()
		if err != nil {
			return 0, fmt.Errorf("Error setting expiration:%s", err.Error())
		}
	}
	return uint(value), nil
}

// 获取本地构建命令信息
func (b *AppBuildService) getCommandDict(stageCommandInfo []StageCommand) (map[uint]Command, error) {
	commandIds := linq.From(stageCommandInfo).
		SelectT(func(stageCommand StageCommand) uint {
			return stageCommand.CommandId
		}).Results()
	commandRepo, err := provider.NewCommandRepo()
	if err != nil {
		return nil, err
	}
	commandDic, err := commandRepo.GetAllCommands(commandIds)
	if err != nil {
		return nil, err
	}
	return commandDic, nil
}

// 执行命令
func (b *AppBuildService) excuteCommand(projectPath string, buildRecore *AppBuildRecord, excuteCommands []AppBuildRecordDetail, commandDict map[uint]Command) {
	engine := NewEngine(projectPath)
	isSuccess := true
	for _, command := range excuteCommands {
		if commandInfo, ok := commandDict[command.CommandId]; ok {
			result := b.excute(engine, command.Command, commandInfo)
			if !result.IsExcuteSuccess && command.SecondCommand != "" {
				result = b.excute(engine, command.SecondCommand, commandInfo)
			}
			command.Output = result.Result
			if result.IsExcuteSuccess {
				command.Result = buildResult.Success
			} else {
				command.Result = buildResult.Fail
				command.Output = result.Error.Error()
			}
			command.RunTime = math.Round(result.Duration*1000) / 1000
		} else {
			command.Output = fmt.Sprintf("构建过程中存在命令变动，命令[%d]不存在", command.CommandId)
			command.Result = buildResult.Fail
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

// 分析命令返回结果是否出错
func (b *AppBuildService) analyzeCommandResult(commandResult string, commandInfo Command) (hasError bool) {
	if commandInfo.ErrorKeyword != "" {
		var errorKeyword []string
		err := json.Unmarshal([]byte(commandInfo.ErrorKeyword), &errorKeyword)
		if err != nil {
			log.Error("Unmarshal ErrorKeyword error:", err)
			return true
		}
		for _, key := range errorKeyword {
			if strings.Contains(strings.ToLower(commandResult), strings.ToLower(key)) {
				log.Info(fmt.Sprintf("%s contains error keyword:%s", commandResult, key))
				return true
			}
		}
	}
	return false
}

// 执行命令
func (b *AppBuildService) excute(engine Engine, command string, commandInfo Command) CommandResultDto {
	result := CommandResultDto{}
	result.Result, result.Duration, result.Error = engine.RunCommand(command)
	if result.Error != nil {
		result.IsExcuteSuccess = false
	} else {
		//校验命令是否包含错误信息
		hasError := b.analyzeCommandResult(result.Result, commandInfo)
		if hasError {
			result.IsExcuteSuccess = false
		} else {
			result.IsExcuteSuccess = true
		}
	}
	return result
}
