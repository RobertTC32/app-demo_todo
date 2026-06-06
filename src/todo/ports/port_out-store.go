package todo_ports

// there are 3 possible reasons for using pointers:
// 1 - we need to change the state of function parameter
//    (parameter values in functions will always be copied on stack frame, and use in function)
// 2 - size of function parameter can be very big
//    (copying big parameter values in functions on stack frame is not efficient for cpu and memory)
// 3 - to be able to return or store a nil value
//    (alternative for returning a nil value is returning an initialized value with error)
// best to minimize pointer usage

type PortOutStore interface {
	Destroy() error
	// todo
	GetTodos() ([]TodoDTO, error)
	GetTodoById(id int32) (*TodoDTO, error)
	CreateTodo(t TodoDTO) (*TodoDTO, error)
	CreateTodoTitle(title string) (*TodoDTO, error)
	UpdateTodo(t TodoDTO) (*TodoDTO, error)
	UpdateTodoCompleted(id uint32, version uint32, completed bool) (*TodoDTO, error)
	DeleteTodo(id uint32, version uint32) (*TodoDTO, error)
	// examples
	GetDbmapping1() ([]Dbmapping1DTO, error)
	GetDbmapping2() ([]Dbmapping2DTO, error)
}
