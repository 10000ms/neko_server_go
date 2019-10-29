package core


type SettingType map[string]string


func (self SettingType) UpdateSettings(newSettings map[string]string) SettingType {
	for k ,v := range newSettings {
		self[k] = v
	}
	return self
}


var BaseSettings = SettingType{
	"Host": "127.0.0.1",
	"Port": "11000",
	"Debug": "False",
}
