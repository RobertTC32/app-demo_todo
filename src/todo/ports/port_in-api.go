package todo_ports

type PortInApi interface {
	Destroy() error
	// todo
	GetSimpleTodos() ([]SimpleTodoDTO, error)
}
