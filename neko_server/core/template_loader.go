package core

import (
    "io/ioutil"
    "neko_server_go/neko_server/utils"
)


type TemplateLoader struct {
    path string
}


func (self *TemplateLoader) GetSource(relativePath string) string {
    realPath := self.path + relativePath
    b, err := ioutil.ReadFile(realPath)
    if err != nil {
        utils.LogFatal(err)
    }
    r := string(b)
    return r
}
