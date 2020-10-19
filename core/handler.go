package core

import (
	"fmt"
	"log"
	"neko_server_go/utils"
	"net/http"
	"runtime"
)

type Handler struct {
    App     *App
    Setting map[string]interface{}
    Router  Router
}

const (
    HttpRequestHeaderRequestID = "request-id"
)

type NekoHandlerFunc func(*Context, http.ResponseWriter, *http.Request)

func RequestLog(r *http.Request) string {
    s := fmt.Sprintf("%s %s (%s)", r.Method, r.RequestURI, r.RemoteAddr)
    return s
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 错误处理
    defer func() {
        if errorRecover := recover(); errorRecover != nil {
            buf := make([]byte, 2048)
            n := runtime.Stack(buf, false)
            stackInfo := fmt.Sprintf("%s", buf[:n])
            utils.LogError("Internal Error: ", errorRecover, stackInfo)
            if self.Setting["Debug"] == true {
                InternalErrorHandler(w, r, errorRecover, stackInfo)
            } else {
                InternalErrorHandler(w, r)
            }
        }
    }()

    utils.LogInfo(RequestLog(r))

    var notFoundHandler NekoHandlerFunc
    // NotFoundHandler配置可以更改默认的NotFoundHandler
    settingNotFoundHandler := self.Setting["NotFoundHandler"]
    if settingNotFoundHandler == nil {
        notFoundHandler = NotFoundHandler
    } else {
        notFoundHandler = settingNotFoundHandler.(NekoHandlerFunc)
    }

    // 处理context
    context := InitContext(self.App, r)

    handler := self.Router.MatchHandler(r.URL.Path, notFoundHandler)
    handler(&context, w, r)
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request, msg ...interface{}) {
    w.WriteHeader(500)
    var m string
    if len(msg) >= 1 {
        for _, i := range msg {
            m += fmt.Sprintf("%s", i)
        }
    } else {
        m = "Internal Error"
    }
    _, err := fmt.Fprintf(w, m)
    if err != nil {
        log.Fatal(err)
        return
    }
}

func NotFoundHandler(c *Context, w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
    _, err := fmt.Fprintf(w, "Not Found")
    if err != nil {
        log.Fatal(err)
        return
    }
}

func DefaultIndexHandler(c *Context, w http.ResponseWriter, r *http.Request) {
    _, err := fmt.Fprintf(w, "neko server go")
    if err != nil {
        log.Fatal(err)
        return
    }
}
