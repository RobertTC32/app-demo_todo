package todo_ports

type PortInUi interface {
	Destroy() error
	// todo
	GetTodos() ([]TodoDTO, error)
	CreateTodoTitle(title string) (*TodoDTO, error)
	UpdateTodoCompleted(id uint32, version uint32) (*TodoDTO, error)
	DeleteTodo(id uint32, version uint32) (*TodoDTO, error)
	// examples
	GetDbmapping1() ([]Dbmapping1DTO, error)
	GetDbmapping2() ([]Dbmapping2DTO, error)
}
