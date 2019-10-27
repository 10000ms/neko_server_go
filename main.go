package main

import (
	"neko_server_go/handler"
	"neko_server_go/neko_server/core"
)


func main() {
	setting := core.Setting{
		Host: "127.0.0.1",
		Port: "11000",
	}

	router := core.Router{
		core.Route{Pattern: "/", HandlerFunc: handler.Index},
	}

	app := core.App{
		Setting: setting,
		Router:  router,
	}

    app.StartApp()
}
