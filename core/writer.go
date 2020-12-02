package core

import "net/http"

type ResWriter interface {
	http.ResponseWriter
}
