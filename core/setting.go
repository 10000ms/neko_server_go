package core

import (
	"errors"
	"neko_server_go/db"
)

type SettingType map[string]interface{}

// 声明一下系统使用的配置项
const (
	SettingServiceName     = "ServiceName"
	SettingHost            = "Host"
	SettingPort            = "Port"
	SettingDebug           = "Debug"
	SettingDb              = "Db"
	SettingDbUsername      = "Username"
	SettingDbPassword      = "Password"
	SettingDbNetwork       = "Network"
	SettingDbServer        = "Server"
	SettingDbPort          = "Port"
	SettingDbDatabase      = "Database"
	SettingNotFoundHandler = "NotFoundHandler"
)

func (s SettingType) UpdateSettings(newSettings map[string]interface{}) SettingType {
	for k, v := range newSettings {
		s[k] = v
	}
	return s
}

func (s *SettingType) DefaultSettings() SettingType {
	return SettingType{
		SettingHost:  "127.0.0.1",
		SettingPort:  "11000",
		SettingDebug: false,
	}
}

func (s *SettingType) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range *s {
		r[k] = v
	}
	return r
}

func (s SettingType) CheckMustContain() error {
	mustContainList := []string{SettingServiceName}
	for _, item := range mustContainList {
		if _, ok := s[item]; !ok {
			return errors.New("missing settings items: " + item)
		}
	}
	return nil
}

func (s SettingType) CheckDbSetting() error {
	if _, ok := s[SettingDb]; ok {
		if _, ok := s[SettingDb].(db.NekoDbSettingType); ok {
			return nil
		} else {
			return errors.New("error db settings")
		}
	} else {
		return errors.New("missing db settings")
	}
}
