package handler

import (
    "fmt"
    "log"
    "net/http"
    "strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()       //解析参数，默认是不会解析的
    if err != nil {
        log.Fatal(err)
        return
    }
    fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    _, err = fmt.Fprintf(w, "Hello Wrold!") //这个写入到w的是输出到客户端的
    if err != nil {
        log.Fatal(err)
        return
    }
}
