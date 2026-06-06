package todo_ports

import (
	"time"
)

// these "Data Transfert Objects" are not ports/adapters, but are used to transfer data between different tiers of the application;
// this "ports" package should not depend on any other packages in the application (to prevent circular dependencies),
// because they can be used in any tier (eg store, ui, api, etc.);
//
// the DTO's in this package should not contain any business logic,
// but are only data structures (eg structs) and data types (eg enums, constants, etc.) used in contracts (interfaces);
//
// to simplify the code, the package is also used inside app code (as entities) and supports db mapping (in store code),
// but in a more complex application it would be better to define separate entities "TodoEntity"
// and db mapping structs in the app and store packages, respectively;

type TodoDTO struct {
	Id          uint32 // mapping "int32(go) <-> int(mysql)" and "int64(go) <-> bigint(mysql)" ("int" in go is platform dependent)
	Title       string
	Completed   bool
	CreatedAt   time.Time  `db:"created_at"`
	CompletedAt *time.Time `db:"completed_at"` // nullable in the database, so we use a pointer
	Version     uint32     // used for optimistic locking
}
