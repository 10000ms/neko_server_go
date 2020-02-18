package core

import (
    "log"
    "net/http"
    "regexp"
)

type Route struct {
    HandlerFunc http.HandlerFunc
}

type Router map[string]func(http.ResponseWriter, *http.Request)

/*
更新router
 */
func (self Router) UpdateHandler(router Router) Router {
    for pattern, handler := range router {
        self[pattern] = handler
    }
    return self
}


/**
根据pattern获取对应的路由，找不到就返回defaultHandler路由
现在的匹配模式暂时是全匹配
*/
func (self Router) MatchHandler(path string, defaultHandler http.HandlerFunc) http.HandlerFunc {
    for pattern, handler := range self {
        match, err := regexp.MatchString(pattern, path)
        if err != nil {
            log.Fatal(err)
        } else {
            if match == true {
                return handler
            }
        }
    }
    return defaultHandler
}


var DefaultRouter = Router{
    "^/$": DefaultIndexHandler,
}
