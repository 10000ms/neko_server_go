package core

import (
	"database/sql"
	"log"
	"neko_server_go/utils"
	"net/http"
	"time"
)

type App struct {
	server        http.Server
	Setting       SettingType
	Router        *Router
	RouterManager *RouterManager
	Db            map[string]*sql.DB
}

func (a *App) StartServer() {
	// 从基础配置更新配置
	if a.Setting["Debug"] == true {
		utils.LogInfo("Debug模式开启")
	}
	var host, port string
	host = a.Setting["Host"].(string)
	port = a.Setting["Port"].(string)
	address := host + ":" + port

	// 更新router
	a.Router = DefaultRouter.MergeRouter(a.Router)
	a.RouterManager = initRouterManager(a.Router)

	handler := Handler{
		App: a,
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
