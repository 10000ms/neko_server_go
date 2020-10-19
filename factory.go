package neko_server

import (
    "database/sql"
    "neko_server_go/core"
    "neko_server_go/db"
    "neko_server_go/utils"
)

/*
创建app的方法

1. setting的初始化
2. db连接初始化
*/
func StartAPP(settings core.SettingType, router core.Router, option core.Options) {
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

    serviceName := settings[core.SettingServiceName].(string)

    utils.LogSystem("////////////////////////////")
    utils.LogSystem("Service Name: " + serviceName)
    utils.LogSystem("////////////////////////////")

    // 初始化db
    Db := make(map[string]*sql.DB)
    if _, ok := settings[core.SettingDb]; ok {
        err := settings.CheckDbSetting()
        if err == nil {
            for dbName, dbConfig := range settings[core.SettingDb].(map[string]map[string]string) {
                username := dbConfig[core.SettingDbUsername]
                password := dbConfig[core.SettingDbPassword]
                network := dbConfig[core.SettingDbNetwork]
                server := dbConfig[core.SettingDbServer]
                port := dbConfig[core.SettingDbPort]
                database := dbConfig[core.SettingDbDatabase]
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

    app := core.App{
        Setting: settings,
        Router:  router,
        Db:      Db,
    }
    app.StartApp()
}

type Setting = core.SettingType
type Router = core.Router
type Options = core.Options
type Context = core.Context
type NekoHandlerFunc = core.NekoHandlerFunc
