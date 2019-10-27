package core

import (
	"log"
	"net/http"
	"time"
)

type App struct {
	server  http.Server
	Setting Setting
	Router  Router
}


func (self *App) StartApp() {
	address := self.Setting.Host + ":" + self.Setting.Port
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
