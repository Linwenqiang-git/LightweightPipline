package ipc

import (
	"bytes"
	"fmt"
	provider "lightweightpipline/cmd/wire"
	. "lightweightpipline/configs/settings/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type CmdEngine struct {
	projectPath    string
	commandExePath string
}

func NewEngine(projectName, projectAddr string) (Engine, error) {
	config, _ := provider.GetConfigure()
	basePath := config.System.BuildRootDir
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)
	excutePath := filepath.Join(basePath, strconv.FormatInt(milliseconds, 10))
	// 创建项目目录
	err := os.MkdirAll(excutePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create build root Path error:" + err.Error())
	}
	//clone 项目
	err = cloneProject(excutePath, projectAddr)
	if err != nil {
		return nil, err
	}
	engine := &CmdEngine{
		projectPath:    excutePath,
		commandExePath: filepath.Join(excutePath, projectName),
	}
	return engine, nil
}

func cloneProject(projectBaseDir, projectAddr string) error {
	command := "cd " + projectBaseDir + " && " + "git clone " + projectAddr
	cmd := getCmd(command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	err := cmd.Run()
	if err != nil {
		errInfo := stdout.String()
		Logger.Error(fmt.Sprintf("excute command:%s \n error:%s \n", command, errInfo))
		return fmt.Errorf("git clone error:%s,detail:%s", err.Error(), errInfo)
	}
	cmd.Wait()
	return nil
}

func getCmd(command string) *exec.Cmd {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	return cmd
}

// 运行命令
func (c *CmdEngine) RunCommand(command string) (string, float64, error) {
	command = "cd " + c.commandExePath + " && " + command
	println(command + "\n")
	cmd := getCmd(command)
	cmd.Env = append(os.Environ(), "LC_CTYPE=en_US.UTF-8")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	startTime := time.Now()
	err := cmd.Run()
	if err != nil {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		errInfo := stdout.String()
		Logger.Error(fmt.Sprintf("excute command:%s \n error:%s \n", command, errInfo))
		return "", duration.Seconds(), fmt.Errorf("excute error:%s,detail:%s", err.Error(), errInfo)
	}
	cmd.Wait()
	result := stdout.String()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	return result, duration.Seconds(), nil
}

// 退出引擎
func (c *CmdEngine) Exit() {
	err := os.RemoveAll(c.projectPath)
	if err != nil {
		Logger.Error(fmt.Sprintf("Engine exist error,delete excute path fail:%s,Error:%s", c.projectPath, err.Error()))
	} else {
		Logger.Info("Directory deleted successfully:" + c.projectPath)
	}
}
