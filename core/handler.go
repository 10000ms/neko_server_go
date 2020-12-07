package core

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	App *App
}

type HttpMethods string

type NekoHandlerFunc func(*Context, ResWriter)

const (
	GetMethodsHandler     HttpMethods = "GET"
	HeadMethodsHandler    HttpMethods = "HEAD"
	PostMethodsHandler    HttpMethods = "POST"
	PutMethodsHandler     HttpMethods = "PUT"
	DeleteMethodsHandler  HttpMethods = "DELETE"
	ConnectMethodsHandler HttpMethods = "CONNECT"
	OptionsMethodsHandler HttpMethods = "OPTIONS"
	TraceMethodsHandler   HttpMethods = "TRACE"
	PatchMethodsHandler   HttpMethods = "PATCH"
)

type MethodsHandler struct {
	HttpMethod HttpMethods
	Handler *NekoHandlerFunc
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

func NotFoundHandler(c *Context, w ResWriter) {
	w.WriteHeader(404)
	_, err := fmt.Fprintf(w, "Not Found")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func DefaultIndexHandler(c *Context, w ResWriter) {
	_, err := fmt.Fprintf(w, "neko server go")
	if err != nil {
		log.Fatal(err)
		return
	}
}
