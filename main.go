package main

import (
	"neko_server_go/neko_server/core"
)


func main() {
	app := core.App{
		Setting: Settings,
		Router:  Router,
	}

    app.StartApp()
}
