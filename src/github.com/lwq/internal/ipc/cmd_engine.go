package ipc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
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
	err := cmd.Run()
	if err != nil {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		errInfo := stdout.String()
		return "", duration.Seconds(), fmt.Errorf("excute error:%s,detail:%s", err.Error(), errInfo)
	}
	cmd.Wait()
	result := stdout.String()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	return result, duration.Seconds(), nil
}
