package main

import "neko_server_go/neko_server/core"

var TempSettings = map[string]string{
    "Host": "127.0.0.1",
    "Port": "11100",
}

var Settings = core.BaseSettings.UpdateSettings(TempSettings)
