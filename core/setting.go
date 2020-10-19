package core

import (
    "errors"
    "neko_server_go/db"
)

type SettingType map[string]interface{}

// 声明一下系统使用的配置项
const (
    SettingServiceName = "ServiceName"
    SettingHost        = "Host"
    SettingPort        = "Port"
    SettingDebug       = "Debug"
    SettingDb          = "Db"
    SettingDbUsername  = "Username"
    SettingDbPassword  = "Password"
    SettingDbNetwork   = "Network"
    SettingDbServer    = "Server"
    SettingDbPort      = "Port"
    SettingDbDatabase  = "Database"
)

func (self SettingType) UpdateSettings(newSettings map[string]interface{}) SettingType {
    for k, v := range newSettings {
        self[k] = v
    }
    return self
}

func (self *SettingType) DefaultSettings() SettingType {
    return SettingType{
        SettingHost:  "127.0.0.1",
        SettingPort:  "11000",
        SettingDebug: false,
    }
}

func (self *SettingType) ToMap() map[string]interface{} {
    r := make(map[string]interface{})
    for k, v := range *self {
        r[k] = v
    }
    return r
}

func (self SettingType) CheckMustContain() error {
    mustContainList := []string{SettingServiceName}
    for _, item := range mustContainList {
        if _, ok := self[item]; !ok {
            return errors.New("missing settings items: " + item)
        }
    }
    return nil
}

func (self SettingType) CheckDbSetting() error {
    if _, ok := self[SettingDb]; ok {
        if _, ok := self[SettingDb].(db.NekoDbSettingType); ok {
            return nil
        } else {
            return errors.New("error db settings")
        }
    } else {
        return errors.New("missing db settings")
    }
}
