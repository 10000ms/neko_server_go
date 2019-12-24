package core

import (
    "io"
    "net/http"
)

type Handler struct {
    Setting map[string]string
    Router  Router
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    handler := self.Router.MatchHandler(r.URL.Path, NotFoundHandler)
    handler(w, r)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
    _, _ = io.WriteString(w, "Not Found")
}
