package ipc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	. "lightweightpipline/configs/settings/log"
)

type CmdEngine struct {
	projectPath string
}

func NewEngine(excutePath string) Engine {
	engine := &CmdEngine{
		projectPath: excutePath,
	}
	return engine
}

func (c *CmdEngine) RunCommand(command string) (string, float64, error) {
	command = "cd " + c.projectPath + " &&" + command
	println(command + "\n")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Env = append(os.Environ(), "LC_CTYPE=en_US.UTF-8")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	startTime := time.Now()
	Logger.Info(fmt.Sprintf("start excute command:%s \n", command))
	err := cmd.Run()
	if err != nil {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		errInfo := stdout.String()
		Logger.Error(fmt.Sprintf("excute command:%s \n error:%s \n", command, errInfo))
		return "", duration.Seconds(), fmt.Errorf("excute error:%s,detail:%s", err.Error(), errInfo)
	}
	cmd.Wait()
	Logger.Info(fmt.Sprintf("end excute command:%s \n", command))
	result := stdout.String()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	return result, duration.Seconds(), nil
}
