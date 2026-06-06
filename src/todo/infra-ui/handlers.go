package todo_ui

import (
	"log/slog"
	"net/http"
	"strconv"
)

// todo

func (this *WebUi) todoHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::todoHandler() - Executing")
	// processing input
	//
	// domain logic
	todos, err := this.app.GetTodos()
	if err != nil {
		slog.Error("ui::todoHandler() - Failed to execute logic", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using templ)
	err = TodoPage(todos).Render(r.Context(), w)
	if err != nil {
		slog.Error("ui::todoHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (this *WebUi) todoPostHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::todoPostHandler() - Executing")
	// processing input
	err := r.ParseForm()
	if err != nil {
		slog.Error("ui::todoPostHandler() - Failed to process input", "error", err)
		// error is system problem and which result in a http error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.Form.Get("id")
	version := r.Form.Get("version")
	action_delete := r.Form.Get("action_delete")
	action_update := r.Form.Get("action_update")
	title := r.Form.Get("title")
	action_insert := r.Form.Get("action_insert")
	//header := r.Header.Values("Content-Type")
	//slog.Debug("ui::todoPostHandler() - Form values", "header", header, "id", id, "version", version, "title", title)
	//slog.Debug("ui::todoPostHandler() - Form actions", "insert", action_insert, "update", action_update, "delete", action_delete)
	//
	// domain logic
	if len(action_insert) > 0 {
		slog.Debug("ui::todoPostHandler() - Inserting", "title", title)
		_, err := this.app.CreateTodoTitle(title)
		if err != nil {
			slog.Error("ui::todoPostHandler() - Failed to execute logic", "error", err)
			// error still results in a view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(action_update) > 0 {
		idNum, _ := strconv.ParseUint(id, 10, 64)
		versionNum, _ := strconv.ParseUint(version, 10, 64)
		slog.Debug("ui::todoPostHandler() - Updating", "id", idNum, "version", versionNum)
		_, err = this.app.UpdateTodoCompleted(uint32(idNum), uint32(versionNum))
		if err != nil {
			slog.Error("ui::todoPostHandler() - Failed to execute logic", "error", err)
			// error still results in a view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(action_delete) > 0 {
		idNum, _ := strconv.ParseUint(id, 10, 64)
		versionNum, _ := strconv.ParseUint(version, 10, 64)
		slog.Debug("ui::todoPostHandler() - Deleting", "id", idNum, "version", versionNum)
		_, err = this.app.DeleteTodo(uint32(idNum), uint32(versionNum))
		if err != nil {
			slog.Error("ui::todoPostHandler() - Failed to execute logic", "error", err)
			// error still results in a view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	todos, err := this.app.GetTodos()
	if err != nil {
		slog.Error("ui::todoPostHandler() - Failed to execute logic", "error", err)
		// error is system problem and which result in a http error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using templ)
	err = TodoList(todos).Render(r.Context(), w)
	if err != nil {
		slog.Error("ui::todoPostHandler() - Failed to process output", "error", err)
		// error is system problem and which result in a http error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
