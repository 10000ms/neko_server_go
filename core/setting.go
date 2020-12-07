package core

import (
	"errors"
	"neko_server_go/db"
	"neko_server_go/enum"
)

type SettingType map[string]interface{}

func (s SettingType) UpdateSettings(newSettings map[string]interface{}) SettingType {
	for k, v := range newSettings {
		s[k] = v
	}
	return s
}

func (s *SettingType) DefaultSettings() SettingType {
	return SettingType{
		enum.SettingHost:  "127.0.0.1",
		enum.SettingPort:  "11000",
		enum.SettingDebug: false,
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
	mustContainList := []string{enum.SettingServiceName}
	for _, item := range mustContainList {
		if _, ok := s[item]; !ok {
			return errors.New("missing settings items: " + item)
		}
	}
	return nil
}

func (s SettingType) CheckDbSetting() error {
	if _, ok := s[enum.SettingDb]; ok {
		if _, ok := s[enum.SettingDb].(db.NekoDbSettingType); ok {
			return nil
		} else {
			return errors.New("error db settings")
		}
	} else {
		return errors.New("missing db settings")
	}
}
