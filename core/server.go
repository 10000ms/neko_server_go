package core

import (
	"fmt"
	"github.com/10000ms/neko_server_go/utils"
	"net/http"
)

func RequestLog(r *http.Request) string {
	s := fmt.Sprintf("%s %s (%s)", r.Method, r.RequestURI, r.RemoteAddr)
	return s
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 错误处理
	defer recoveryHandlerPanic(w, r)

	utils.LogInfo(RequestLog(r))

	var notFoundHandler NekoHandlerFunc
	// NotFoundHandler配置可以更改默认的NotFoundHandler
	settingNotFoundHandler := self.Setting["NotFoundHandler"]
	if settingNotFoundHandler == nil {
		notFoundHandler = NotFoundHandler
	} else {
		notFoundHandler = settingNotFoundHandler.(NekoHandlerFunc)
	}

	// 获取writer
	writer := w.(ResWriter)

	// 处理context
	context := InitContext(self.App, r)

	handler := self.Router.MatchHandler(r.URL.Path, notFoundHandler)
	handler(&context, writer)
}
