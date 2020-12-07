package neko_server_go

import (
	"database/sql"
	"neko_server_go/core"
	"neko_server_go/db"
	"neko_server_go/enum"
	"neko_server_go/utils"
)

/*
创建app的方法

1. setting的初始化
2. db连接初始化
*/
func StartAPP(settings core.SettingType, router *core.Router, option *core.Options) {
	for _, f := range option.InitFunc {
		f()
	}

	defer func() {
		for _, f := range option.DeferFunc {
			f()
		}
	}()

	var err error
	err = settings.CheckMustContain()
	if err != nil {
		utils.LogError(err)
		return
	}

	serviceName := settings[enum.SettingServiceName].(string)

	// ..
	utils.LogSystem("                                                                                                \n               ,--.                                                                             \n,--,--,  ,---. |  |,-. ,---.      ,---.  ,---. ,--.--.,--.  ,--.,---. ,--.--.     ,---.  ,---.  \n|      \\| .-. :|     /| .-. |    (  .-' | .-. :|  .--' \\  `'  /| .-. :|  .--'    | .-. || .-. | \n|  ||  |\\   --.|  \\  \\' '-' '    .-'  `)\\   --.|  |     \\    / \\   --.|  |       ' '-' '' '-' ' \n`--''--' `----'`--'`--'`---'     `----'  `----'`--'      `--'   `----'`--'       .`-  /  `---'  \n                                                                                 `---'          ")

	utils.LogSystem("////////////////////////////")
	utils.LogSystem("Service Name: " + serviceName)
	utils.LogSystem("////////////////////////////")

	// 初始化db
	Db := make(map[string]*sql.DB)
	if _, ok := settings[enum.SettingDb]; ok {
		err := settings.CheckDbSetting()
		if err == nil {
			for dbName, dbConfig := range settings[enum.SettingDb].(map[string]map[string]string) {
				username := dbConfig[enum.SettingDbUsername]
				password := dbConfig[enum.SettingDbPassword]
				network := dbConfig[enum.SettingDbNetwork]
				server := dbConfig[enum.SettingDbServer]
				port := dbConfig[enum.SettingDbPort]
				database := dbConfig[enum.SettingDbDatabase]
				dbConn := db.CreateDbConnect(username, password, network, server, port, database)
				if dbConn != nil {
					Db[dbName] = dbConn
					utils.LogSystem("connect to db: " + dbName)
				}
			}
		} else {
			utils.LogError(err)
		}
		utils.LogSystem("////////////////////////////")
	}

	// 处理404handler
	var notFoundHandler NekoHandlerFunc
	// NotFoundHandler配置可以更改默认的NotFoundHandler
	settingNotFoundHandler := settings[enum.SettingNotFoundHandler]
	if settingNotFoundHandler == nil {
		notFoundHandler = core.NotFoundHandler
	} else {
		notFoundHandler = settingNotFoundHandler.(NekoHandlerFunc)
	}
	settings[enum.SettingNotFoundHandler] = notFoundHandler

	app := core.App{
		Setting: settings,
		Router:  router,
		Db:      Db,
	}
	app.StartServer()
}

type Setting = core.SettingType
type Router = core.Router
type Options = core.Options
type Context = core.Context
type NekoHandlerFunc = core.NekoHandlerFunc
type ResWriter = core.ResWriter
