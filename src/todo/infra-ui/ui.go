package todo_ui

import (
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

// "ui" (with views and handlers) contains adapter for web gui in the "Hexagonal Architecture"

type WebUi struct {
	app         todo_ports.PortInUi
	resourcesFs fs.FS
	router      *http.ServeMux
}

// prevent parsing template files on every request by caching them;
// only used in listHandler example where template library is used
var templates = make(map[string]*template.Template)

func NewWebUi(app todo_ports.PortInUi, resourcesFs fs.FS, router *http.ServeMux) (*WebUi, error) {
	slog.Debug("ui::NewWebUi() - Executing")
	ui := &WebUi{
		app:         app,
		resourcesFs: resourcesFs,
		router:      router,
	}
	//
	// add todo handlers
	router.Handle("GET /public/", http.FileServerFS(resourcesFs))
	router.HandleFunc("GET /", ui.defaultHandler)
	router.HandleFunc("GET /todo", ui.todoHandler)
	router.HandleFunc("POST /todo", ui.todoPostHandler)
	//
	// add extra example handlers
	router.HandleFunc("GET /list", ui.listHandler)
	router.HandleFunc("GET /mapping", ui.mappingHandler)
	router.HandleFunc("GET /hello", ui.helloHandler)
	// prevent parsing template files on every request by caching them
	if len(templates) == 0 {
		templates["list.html"] = template.Must(template.ParseFS(resourcesFs, "templates/list.html"))
	}
	//
	return ui, nil
}

func (this *WebUi) Destroy() error {
	slog.Debug("ui::Destroy() - Executing")
	return nil
}

func (this *WebUi) defaultHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::defaultHandler() - Executing")
	http.Redirect(w, r, "/todo", http.StatusMovedPermanently)
}
