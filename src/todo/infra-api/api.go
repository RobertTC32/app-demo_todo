package todo_api

import (
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"encoding/json"
	"log/slog"
	"net/http"
)

// "api" contains adapter for web service with REST API in the "Hexagonal Architecture"

type RestApi struct {
	app    todo_ports.PortInApi
	router *http.ServeMux
}

func NewRestApi(app todo_ports.PortInApi, router *http.ServeMux) (*RestApi, error) {
	slog.Debug("api::NewRestApi() - Executing")
	api := &RestApi{
		app:    app,
		router: router,
	}
	//
	// add todo handlers
	router.HandleFunc("GET /api/todo", api.todoHandler)
	//
	return api, nil
}

func (this *RestApi) Destroy() error {
	slog.Debug("api::Destroy() - Executing")
	return nil
}

func (this *RestApi) todoHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("api::todoHandler() - Executing")
	stodos, err := this.app.GetSimpleTodos()
	if err != nil {
		slog.Error("api::todoHandler() - Failed to execute logic", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stodos)
}
