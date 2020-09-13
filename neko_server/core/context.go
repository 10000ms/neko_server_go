package core

import (
	uuid "github.com/satori/go.uuid"
	"neko_server_go/neko_server/utils"
	"net/http"
)

type Context struct {
    RequestID string
    App       *App
}

func InitContext(app *App, r *http.Request) Context {
    // 尝试获取header里面的requestID
    requestID := r.Header.Get(HttpRequestHeaderRequestID)
    utils.LogDebug("original request id is: ", requestID)
    if requestID == "" {
        requestID = uuid.Must(uuid.NewV4()).String()
        utils.LogDebug("generate request id is: ", requestID)
    }
    c := Context{
        App:       app,
        RequestID: requestID,
    }
    return c
}

func (c *Context) LogRequestID() string {
    s := "[" + c.RequestID + "]"
    return s
}
