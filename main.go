package main

import (
    "neko_server_go/neko_server"
)

func main() {
    o := neko_server.Options{}
    neko_server.StartAPP(Settings, Router, o)
}
