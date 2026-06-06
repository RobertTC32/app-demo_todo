package todo_store

import (
	"RobertTC32/example-demo_hello/src/commons"
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// todo

func showErrorDetails(err error) {
	// this function is only used to experiment with database errors
	if err == nil {
		slog.Error("NoError")
		return
	}
	if strings.HasPrefix(err.Error(), "sql: ") {
		// sentinel errors returned from "database/sql" package
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("SQLError", "ErrNoRows", err)
		} else if errors.Is(err, sql.ErrConnDone) {
			slog.Error("SQLError", "ErrConnDone", err)
		} else if errors.Is(err, sql.ErrTxDone) {
			slog.Error("SQLError", "ErrTxDone", err)
		} else {
			slog.Error("SQLError", "NoName", err)
		}
		return
	}
	number, values, err2 := ParseMysqlError(err)
	if err2 == nil {
		slog.Error("MySQLError", "number", number, "values", values)
		return
	}
	slog.Error("OtherError", "err", err)
}

func (this *MysqlStore) GetTodos() ([]todo_ports.TodoDTO, error) {
	slog.Debug("store::GetTodos() - Executing")
	var tt []todo_ports.TodoDTO
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos ORDER BY id`
	if err := this.Db.Select(&tt, sqlSelect); err != nil {
		slog.Error("store::GetTodos() - Failed to find todos", "error", err)
		return []todo_ports.TodoDTO{}, err
	}
	// no mapping, no null handling, no rows closing and no error handling needed with sqlx
	slog.Debug("store::GetTodos() - Retrieved number", "todos", strconv.Itoa(len(tt)))
	return tt, nil
}

func (this *MysqlStore) GetTodoById(id int32) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::GetTodoById() - Executing")
	var t todo_ports.TodoDTO
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	// MySQL/MariaDB only accept "?" and no "$1" in sql syntax
	if err := this.Db.Get(&t, sqlSelect, id); err != nil {
		slog.Error("store::GetTodoById() - Failed to find todo", "error", err)
		if err == sql.ErrNoRows {
			return nil, commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
		}
		return nil, err
	}
	return &t, nil
}

const TODOS_MAX = 20

func (this *MysqlStore) CreateTodo(s todo_ports.TodoDTO) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::CreateTodo() - Executing")
	//
	var idCount int
	if err := this.Db.Get(&idCount, "SELECT count(id) as id_count FROM todos"); err != nil {
		slog.Error("store::CreateTodo() - Failed to count todo", "error", err)
		return nil, err
	}
	if idCount >= TODOS_MAX {
		slog.Error("store::CreateTodo() - Max todo exceeded", "maximum", TODOS_MAX)
		return nil, commons.NewWrappingError(todo_ports.ErrMaxRowsExceededInStore, nil)
	}
	//
	sqlInsert := `INSERT INTO todos (title, completed, created_at, completed_at) 
		VALUES (?, ?, ?, ?) RETURNING *`
	// id and version are calculated by the database, so we need to get them back from database;
	// INSERT with RETURNING is allowed since MariaDB 10.5
	var t todo_ports.TodoDTO
	if err := this.Db.Get(&t, sqlInsert, t.Title, t.Completed, t.CreatedAt, t.CompletedAt); err != nil {
		slog.Error("store::CreateTodo() - Failed to insert todo", "error", err)
		if de, ok := err.(*mysql.MySQLError); ok && (de.Number == 1062) {
			return nil, commons.NewWrappingError(todo_ports.ErrAlreadyExistsInStore, err)
		}
		return nil, commons.NewWrappingError(todo_ports.ErrConstraintViolatedInStore, err)
	}
	return &t, nil
}

func (this *MysqlStore) CreateTodoTitle(title string) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::CreateTodoTitle() - Executing")
	//
	var idCount int
	if err := this.Db.Get(&idCount, "SELECT count(id) as id_count FROM todos"); err != nil {
		slog.Error("store::CreateTodoTitle() - Failed to count todo", "error", err)
		return nil, err
	}
	if idCount >= TODOS_MAX {
		slog.Error("store::CreateTodoTitle() - Max todo exceeded", "maximum", TODOS_MAX)
		return nil, commons.NewWrappingError(todo_ports.ErrMaxRowsExceededInStore, nil)
	}
	//
	sqlInsert := `INSERT INTO todos (title, completed, created_at) 
		VALUES (?, FALSE, NOW()) RETURNING *`
	// id and version are calculated by the database, so we need to get them back from database;
	// INSERT with RETURNING is allowed since MariaDB 10.5
	var t todo_ports.TodoDTO
	if err := this.Db.Get(&t, sqlInsert, title); err != nil {
		slog.Error("store::CreateTodoTitle() - Failed to insert todo", "error", err)
		if de, ok := err.(*mysql.MySQLError); ok && (de.Number == 1062) {
			showErrorDetails(err)
			return nil, commons.NewWrappingError(todo_ports.ErrAlreadyExistsInStore, err)
		}
		return nil, commons.NewWrappingError(todo_ports.ErrConstraintViolatedInStore, err)
	}
	return &t, nil
}

func (this *MysqlStore) UpdateTodo(s todo_ports.TodoDTO) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::UpdateTodo() - Executing")
	//
	sqlUpdate := `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
	WHERE id = ?`
	// optimistic locking
	if s.Version > 0 {
		sqlUpdate = `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
		WHERE id = ? and version = ?`
	}
	// id is calculated by the database, and version must be changed by client;
	// UPDATE with RETURNING is not allowed for MariaDB;
	// REPLACE with RETURNING is allowed since MariaDB 10.5, but can not be used for optimistic locking
	result, err := this.Db.Exec(sqlUpdate, s.Title, s.Completed, s.CreatedAt, s.CompletedAt, s.Id, s.Version)
	rowsAffected, _ := result.RowsAffected()
	//
	if err != nil {
		slog.Error("store::UpdateTodo() - Failed to update todo", "error", err)
		return nil, commons.NewWrappingError(todo_ports.ErrConstraintViolatedInStore, err)
	}
	// check if the todo was found and updated
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	var t todo_ports.TodoDTO
	if err2 := this.Db.Get(&t, sqlSelect, s.Id); err2 != nil {
		slog.Error("store::UpdateTodo() - Failed to find todo by id", "error", err2)
		return nil, commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	if rowsAffected == 0 {
		slog.Error("store::UpdateTodo() - Failed to find todo by id and rowversion", "error", err)
		return nil, commons.NewWrappingError(todo_ports.ErrChangedByUserInStore, err)
	}
	return &t, nil
}

func (this *MysqlStore) UpdateTodoCompleted(id uint32, version uint32, completed bool) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::UpdateTodoCompleted() - Executing")
	sqlUpdate := `UPDATE todos SET completed = ?, completed_at = ?, version = version + 1 
		WHERE id = ?`
	// optimistic locking
	if version > 0 {
		sqlUpdate = `UPDATE todos SET completed = ?, completed_at = ?, version = version + 1 
			WHERE id = ? and version = ?`
	}
	// id is calculated by the database, and version must be changed by client;
	// UPDATE with RETURNING is not allowed for MariaDB;
	// REPLACE with RETURNING is allowed since MariaDB 10.5, but can not be used for optimistic locking
	var completedAt *time.Time
	completedAt = nil
	if completed {
		now := time.Now()
		completedAt = &now
	}
	result, err := this.Db.Exec(sqlUpdate, completed, completedAt, id, version)
	rowsAffected, _ := result.RowsAffected()
	//
	if err != nil {
		slog.Error("store::UpdateTodoCompleted() - Failed to update todo", "error", err)
		return nil, commons.NewWrappingError(todo_ports.ErrConstraintViolatedInStore, err)
	}
	var t todo_ports.TodoDTO
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	if err2 := this.Db.Get(&t, sqlSelect, id); err2 != nil {
		slog.Error("store::UpdateTodoCompleted() - Failed to find todo by id", "error", err2)
		return nil, commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	if rowsAffected == 0 {
		slog.Error("store::UpdateTodoCompleted() - Failed to find todo by id and rowversion", "error", err)
		return nil, commons.NewWrappingError(todo_ports.ErrChangedByUserInStore, err)
	}
	return &t, nil
}

func (this *MysqlStore) DeleteTodo(id uint32, version uint32) (*todo_ports.TodoDTO, error) {
	slog.Debug("store::DeleteTodo() - Executing")
	sqlDelete := `DELETE FROM todos WHERE id = ? RETURNING *`
	// optimistic locking
	if version > 0 {
		sqlDelete = `DELETE FROM todos WHERE id = ? AND version = ? RETURNING *`
	}
	// DELETE with RETURNING is allowed since MariaDB 10.0
	var t todo_ports.TodoDTO
	if err := this.Db.Get(&t, sqlDelete, id, version); err != nil {
		// not found or sql error
		sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
		if err2 := this.Db.Get(&t, sqlSelect, id); err2 != nil {
			slog.Error("store::DeleteTodo() - Failed to find todo by id", "error", err2)
			return nil, commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
		}
		slog.Error("store::DeleteTodo() - Failed to delete todo", "error", err)
		showErrorDetails(err)
		return nil, commons.NewWrappingError(todo_ports.ErrChangedByUserInStore, err)
	}
	return &t, nil
}

func (this *MysqlStore) SwitchTodoTitles(id1 uint32, id2 uint32) error {
	// example using transactions with optimistic locking
	slog.Debug("store::SwitchTodoTitles() - Executing")
	// we want to switch the titles of two todos
	tx, err := this.Db.BeginTxx(context.Background(), nil)
	if err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback()
	//
	var t1 todo_ports.TodoDTO
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	if err := this.Db.Get(&t1, sqlSelect, id1); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	var t2 todo_ports.TodoDTO
	if err := this.Db.Get(&t2, sqlSelect, id2); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	sqlUpdate := `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
			WHERE id = ? and version = ?`
	if _, err := this.Db.Exec(sqlUpdate, t2.Title, t2.Completed, t2.CreatedAt, t2.CompletedAt, t1.Id, t1.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	if _, err := this.Db.Exec(sqlUpdate, t1.Title, t1.Completed, t1.CreatedAt, t1.CompletedAt, t2.Id, t2.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return commons.NewWrappingError(todo_ports.ErrNotFoundInStore, err)
	}
	if err := tx.Commit(); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to commit", "error", err)
		return err
	}
	return nil
}
