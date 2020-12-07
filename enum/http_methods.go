package enum

type HttpMethods string

const (
	HttpMethodsGet     HttpMethods = "GET"
	HttpMethodsHead    HttpMethods = "HEAD"
	HttpMethodsPost    HttpMethods = "POST"
	HttpMethodsPut     HttpMethods = "PUT"
	HttpMethodsDelete  HttpMethods = "DELETE"
	HttpMethodsConnect HttpMethods = "CONNECT"
	HttpMethodsOptions HttpMethods = "OPTIONS"
	HttpMethodsTrace   HttpMethods = "TRACE"
	HttpMethodsPatch   HttpMethods = "PATCH"
)
