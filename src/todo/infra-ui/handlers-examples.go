package todo_ui

import (
	"fmt"
	"log/slog"
	"net/http"
)

// examples

// easy template example
func (this *WebUi) listHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::listHandler() - Executing")
	// processing input
	//
	// domain logic
	slog.Debug("ui::listHandler() - Reading data")
	d, err := this.app.GetTodos()
	if err != nil {
		slog.Error("ui::listHandler() - Failed to execute logic", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using template)
	slog.Debug("ui::listHandler() - Building response")
	tmpl := templates["list.html"]
	err = tmpl.Execute(w, d)
	if err != nil {
		slog.Error("ui::listHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// example to show mapping between mysql/mariadb and go data types (including "null")
func (this *WebUi) mappingHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::dbmappingHandler() - Executing")
	// processing input
	//
	// domain logic
	this.app.GetDbmapping1()
	this.app.GetDbmapping2()
	//
	// processing output (using fmt)
	fmt.Fprintf(w, "Results of dbmapping shown in console")
}

// easy templ example
func (this *WebUi) helloHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::helloTemplHandler() - Executing")
	// processing input
	//
	// domain logic
	//
	// processing output (using templ)
	err := HelloPage("Robert The Coder").Render(r.Context(), w)
	if err != nil {
		slog.Error("ui::helloHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
