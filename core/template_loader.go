package core

import (
    "io/ioutil"
    "neko_server_go/utils"
)

type TemplateLoader struct {
    Path string
}

func (self *TemplateLoader) GetSource(relativePath string) string {
    realPath := self.Path + relativePath
    b, err := ioutil.ReadFile(realPath)
    if err != nil {
        utils.LogFatal(err)
    }
    r := string(b)
    return r
}
