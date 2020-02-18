package core

import (
    "log"
    "neko_server_go/neko_server/utils"
    "net/http"
    "time"
)

type App struct {
    server  http.Server
    Setting SettingType
    Router  Router
}

func (self *App) StartApp() {
    // 从基础配置更新配置
    self.Setting = BaseSettings.UpdateSettings(self.Setting)
    if self.Setting["Debug"] == true {
        utils.LogInfo("Debug模式开启")
    }
    var host, port string
    host = self.Setting["Host"].(string)
    port = self.Setting["Port"].(string)
    address := host + ":" + port
    // 更新router
    self.Router = DefaultRouter.UpdateHandler(self.Router)
    handler := Handler{
        Setting: self.Setting,
        Router:  self.Router,
    }
    server := http.Server{
        Addr:        address,
        Handler:     &handler,
        ReadTimeout: time.Second * 5, // 超时设置
    }
    utils.LogInfo("server start listen ", address)
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}
