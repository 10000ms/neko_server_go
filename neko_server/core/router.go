package core

import (
	"net/http"
)


type Route struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
}


type Router []Route

/**
根据pattern获取对应的路由，找不到就返回defaultHandler路由
现在的匹配模式暂时是全匹配
 */
func (self *Router) MatchHandler(path string, defaultHandler http.HandlerFunc) http.HandlerFunc {
	for _, route := range *self {
		if path == route.Pattern {
			return route.HandlerFunc
		}
	}
	return defaultHandler
}
