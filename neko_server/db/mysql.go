package db

import (
    "database/sql"
    "neko_server_go/neko_server/core"
    "neko_server_go/neko_server/utils"
)

type MysqlOperate struct {
    Db *sql.DB
}

func (self *MysqlOperate) Init(setting core.SettingType) {
    var err error
    dataSourceName :=
        setting["DbUser"].(string) +
        setting["DbPassword"].(string) +
        "@tcp(" +
        setting["DbHost"].(string) +
        ":" +
        setting["DbPort"].(string) +
        ")/" +
        setting["DbDb"].(string)
    self.Db, err = sql.Open(setting["DbType"].(string), dataSourceName)
    if err != nil {
        utils.LogFatal(err)
        panic(err)
        return
    }
}
