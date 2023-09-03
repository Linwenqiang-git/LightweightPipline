package ipc

type Engine interface {
	RunCommand(command string) (string, float64, error)
	Exit()
}
