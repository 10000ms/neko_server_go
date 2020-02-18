package core

type SettingType map[string]interface{}

func (self SettingType) UpdateSettings(newSettings map[string]interface{}) SettingType {
    for k, v := range newSettings {
        self[k] = v
    }
    return self
}


var BaseSettings = SettingType{
    "Host":  "127.0.0.1",
    "Port":  "11000",
    "Debug": false,
}
