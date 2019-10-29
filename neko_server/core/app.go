package core

import (
	"log"
	"net/http"
	"time"
)

type App struct {
	server  http.Server
	Setting map[string]string
	Router  Router
}


func (self *App) StartApp() {
	var host, port string
	if self.Setting == nil {
		host = BaseSettings["Host"]
		port = BaseSettings["Port"]
	} else {
		host = self.Setting["Host"]
		port = self.Setting["Port"]
	}
	address := host + ":" + port
	handler := Handler{
		Setting: self.Setting,
		Router:  self.Router,
	}
	server := http.Server{
		Addr:        address,
		Handler:     &handler,
		ReadTimeout: time.Second * 5, // 超时设置
	}
	log.Print("server start listen " + address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
