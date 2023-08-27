package ipc

type Engine interface {
	RunCommand(command string) (string, error)
}
