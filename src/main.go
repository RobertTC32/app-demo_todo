package main

import (
	"RobertTC32/example-demo_hello/src/commons"
	todo_app "RobertTC32/example-demo_hello/src/todo/app"
	todo_api "RobertTC32/example-demo_hello/src/todo/infra-api"
	todo_store "RobertTC32/example-demo_hello/src/todo/infra-store"
	todo_ui "RobertTC32/example-demo_hello/src/todo/infra-ui"
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// src folder (or special "src/cmd/todo" subfolder) contains logic to handle the application startup/shutdown,
// and logic to coordinate between bounded context or feature (use case) based modules of VERTICAL SLICE ARCHITECTURE;
// BUSINESS MODULES can be used as isolated parts of the business logic with HIGH COHESION and LOW COUPLING between modules,
// which makes them easier to separately test, maintain and evolve;
// the integration of multiple modules is only needed in bigger applications;
// this simple crud application contains only one module called "todo" (all parts have "todo_" prefix in their names)
//
// each module has a HEXAGONAL ARCHITECTURE to CONTROL TECHNICAL DEPENDENCIES and support TESTABILITY:
// - system inside:
// 	Which parts?
// 	"application/services" consists of features (logic for use cases) and model (business domain logic and entities)
// - system outside:
// 	system inside (containing the application/domain logic) needs to be protected from uncontrolled outside (containing external I/O) access;,
// 	using contracts (called 'ports') and technical implementations (called 'adapters');
// 	this makes the system also easier to maintain and easier to test (by using mock implementations of the adapters);
// 	two types of adapters are used:
//  * "inbound"/"driving" adapters (to receive requests):
// 	  eg "store" (data persistence, eg DB, file system), "called api/service" (eg REST API, message queue, email, sms);
// 	* "outbound"/"driven" adapters (to send requests):
//     eg "view/controller/presentation" (e.g. web pages), "calling api/service" (eg REST API, message queue, email, sms);
// because all technologies are used in the adapters, the system inside is independent of any technology and can be tested without any technology;
// the system inside can be used in different applications with different adapters (eg web, cli, api, etc.),
// which makes it easier to test, maintain and evolve the system inside;
// the application uses dependency injection for the "store", "api" and "ui" adapter;

//go:embed resources
var resourcesEmbed embed.FS

// store (database) has application lifetime,
// and needs to be gracefully closed

func main() {
	commons.LoadEnvFile()
	commons.InitLoggerFromEnv()
	slog.Debug("main::main() - Executing")
	//
	// create (and destroy) "store" port implementation
	var store todo_ports.PortOutStore
	store, err := todo_store.NewMysqlStore()
	if err != nil {
		slog.Error("main::main() - Failed to create store", "error", err)
		return
	}
	defer store.Destroy()
	//
	// create (and destroy) "app4ui" and "app4api" port implementations with "store" injected
	app, err := todo_app.NewAppService(store)
	if err != nil {
		slog.Error("main::main() - Failed to create app", "error", err)
		return
	}
	defer app.Destroy()
	//
	// create router for "ui" and "api" implementations on web server
	router := http.NewServeMux()
	//
	// create (and destroy) "api" port implementation with "app4api" injected
	var app4api todo_ports.PortInApi
	app4api = app
	api, err := todo_api.NewRestApi(app4api, router)
	if err != nil {
		slog.Error("main::main() - Failed to create API", "error", err)
		return
	}
	defer api.Destroy()
	//
	// get access to the embedded resources folder
	//resourcesFs := os.DirFS("src/resources")
	resourcesFs, _ := fs.Sub(resourcesEmbed, "resources")
	//
	// create (and destroy) "ui" port implementation with "app4ui" injected
	var app4ui todo_ports.PortInUi
	app4ui = app
	ui, err := todo_ui.NewWebUi(app4ui, resourcesFs, router)
	if err != nil {
		slog.Error("main::main() - Failed to create UI", "error", err)
		return
	}
	defer ui.Destroy()
	//
	// start ui using web server
	srv, _ := commons.NewServer(router)
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	slog.Info("main::main() - Web Server is available at http://" + host + ":" + port)
	slog.Info("main::main() - Press Ctrl+C to stop")
	if err := srv.RunServer(context.Background(), 5*time.Second); err != nil {
		slog.Error("main::main() - Server error", "error", err)
	}
	//
	slog.Info("main::main() - Stopped")
}
