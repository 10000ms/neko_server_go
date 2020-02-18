package core

import (
    "fmt"
    "log"
    "neko_server_go/neko_server/utils"
    "net/http"
    "runtime"
)

type Handler struct {
    Setting map[string]interface{}
    Router  Router
}

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

    var notFoundHandler func(http.ResponseWriter, *http.Request)
    // NotFoundHandler配置可以更改默认的NotFoundHandler
    settingNotFoundHandler := self.Setting["NotFoundHandler"]
    if settingNotFoundHandler == nil {
        notFoundHandler = NotFoundHandler
    } else {
        notFoundHandler = settingNotFoundHandler.(func(http.ResponseWriter, *http.Request))
    }

    handler := self.Router.MatchHandler(r.URL.Path, notFoundHandler)
    handler(w, r)
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


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
    _, err := fmt.Fprintf(w, "Not Found")
    if err != nil {
        log.Fatal(err)
        return
    }
}

func DefaultIndexHandler(w http.ResponseWriter, r *http.Request) {
    _, err := fmt.Fprintf(w, "neko server go")
    if err != nil {
        log.Fatal(err)
        return
    }
}
