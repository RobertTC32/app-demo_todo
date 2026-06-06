package todo_app

import (
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"log/slog"
)

// "app" contains the core with application- and domain-services (aka use cases and entities) in the "Hexagonal Architecture";
// this simple crud application contains no real business logic

type AppService struct {
	store todo_ports.PortOutStore
}

func NewAppService(store todo_ports.PortOutStore) (*AppService, error) {
	slog.Debug("app::NewAppService() - Executing")
	return &AppService{
		store: store,
	}, nil
}

func (this *AppService) Destroy() error {
	slog.Debug("app::Destroy() - Executing")
	return nil
}

func (this *AppService) GetSimpleTodos() ([]todo_ports.SimpleTodoDTO, error) {
	slog.Debug("app::GetSimpleTodos() - Executing")
	todos, err := this.store.GetTodos()
	if err != nil {
		return nil, err
	}
	var stodos []todo_ports.SimpleTodoDTO
	for _, item := range todos {
		stodo := todo_ports.SimpleTodoDTO{Title: item.Title, Completed: item.Completed, CreatedAt: item.CreatedAt}
		stodos = append(stodos, stodo)
	}
	return stodos, nil
}

func (this *AppService) GetTodos() ([]todo_ports.TodoDTO, error) {
	return this.store.GetTodos()
}

func (this *AppService) CreateTodoTitle(title string) (*todo_ports.TodoDTO, error) {
	return this.store.CreateTodoTitle(title)
}

func (this *AppService) UpdateTodoCompleted(id uint32, version uint32) (*todo_ports.TodoDTO, error) {
	return this.store.UpdateTodoCompleted(id, version, true)
}

func (this *AppService) DeleteTodo(id uint32, version uint32) (*todo_ports.TodoDTO, error) {
	return this.store.DeleteTodo(id, version)
}

// examples

func (this *AppService) GetDbmapping1() ([]todo_ports.Dbmapping1DTO, error) {
	return this.store.GetDbmapping1()
}

func (this *AppService) GetDbmapping2() ([]todo_ports.Dbmapping2DTO, error) {
	return this.store.GetDbmapping2()
}
