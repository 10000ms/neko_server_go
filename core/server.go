package core

import (
	"fmt"
	"neko_server_go/utils"
	"net/http"
)

func RequestLog(c *Context, r *http.Request) string {
	s := fmt.Sprintf("%s %s %s (%s)", c.LogRequestID(), r.Method, r.RequestURI, r.RemoteAddr)
	return s
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 错误处理
	defer recoveryHandlerPanic(w, r)

	// 获取writer
	writer := w.(ResWriter)

	// 处理context
	context := InitContext(h.App, r)
	utils.LogInfo(RequestLog(&context, r))

	NotFoundHandler := h.App.Setting[SettingNotFoundHandler].(NekoHandlerFunc)
	handler, pathParams := h.App.RouterManager.MatchHandler(
		context.Request.Method,
		context.Request.URL.Path,
		&NotFoundHandler,
	)
	context.PathParams = pathParams
	c := *handler
	c(&context, writer)
}
