package todo_ports

import (
	"time"
)

// "dto" is not an adapter, but contains "Data Transfert Objects" which are used to transfer data between isolated parts of the application;
//
// to simplify the code, the package is also supports json mapping (in api code),
// but in a more complex application it would be better to define api mapping structs "TodoJSON" in the api package;

type SimpleTodoDTO struct {
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"creation_time"`
}
