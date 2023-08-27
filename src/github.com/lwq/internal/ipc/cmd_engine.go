package ipc

import (
	"bytes"
	"fmt"
	"os/exec"
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

func (c *CmdEngine) RunCommand(command string) (string, error) {
	command = "cd " + c.projectPath + " ;" + command
	println(command + "\n")
	cmd := exec.Command("sh", command)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Command output:", stdout.String())
		return "", fmt.Errorf("excute error:%s,detail:%s", err.Error(), stdout.String())
	}
	cmd.Wait()
	result := stdout.String()
	return result, nil
}
