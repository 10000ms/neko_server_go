package core

import (
	uuid "github.com/satori/go.uuid"
	"neko_server_go/enum"
	"neko_server_go/utils"
	"net/http"
)

type Context struct {
	RequestID  string
	App        *App
	Request    *http.Request
	PathParams *map[string]string
}

func InitContext(app *App, r *http.Request) Context {
	// 尝试获取header里面的requestID
	requestID := r.Header.Get(enum.HttpRequestHeaderRequestID)
	utils.LogDebug("original request id is: ", requestID)
	if requestID == "" {
		requestID = uuid.Must(uuid.NewV4(), nil).String()
		utils.LogDebug("generate request id is: ", requestID)
	}
	c := Context{
		App:       app,
		RequestID: requestID,
		Request:   r,
	}
	return c
}

func (c *Context) LogRequestID() string {
	s := "[" + c.RequestID + "]"
	return s
}
