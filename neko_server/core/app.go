package core

import (
	"database/sql"
	"log"
	"neko_server_go/neko_server/utils"
	"net/http"
	"time"
)

type App struct {
    server  http.Server
    Setting SettingType
    Router  Router
    Db      map[string]*sql.DB
}

func (self *App) StartApp() {
    // 从基础配置更新配置
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
        App:     self,
        Setting: self.Setting,
        Router:  self.Router,
    }
    server := http.Server{
        Addr:    address,
        Handler: &handler,
        // TODO：增加配置能力
        // 超时设置
        ReadTimeout:  time.Second * 5,
        WriteTimeout: time.Second * 5,
    }
    utils.LogInfo("server start listen ", address)
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}
