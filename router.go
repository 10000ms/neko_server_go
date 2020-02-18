package main

import (
    "neko_server_go/handler"
    "neko_server_go/neko_server/core"
)

var Router = core.Router{
    "^/$": handler.Index,
}
