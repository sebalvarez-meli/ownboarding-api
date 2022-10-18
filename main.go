package main

import (
	"github.com/api_base/config"
	"github.com/api_base/internal/conectivity"
	"github.com/api_base/internal/domain"
	"github.com/api_base/internal/domain/user"
	"log"
	"net/http"
)

func main() {
	//Configuration
	conf := config.NewConfig()
	//Dependencies
	ctn := domain.NewContainer(conf)
	srv := user.NewService(ctn)
	hdlFunc := conectivity.NewHandlerFunc(srv)
	//Router
	router := conectivity.NewRouterHandler(hdlFunc)
	//Start server
	err := http.ListenAndServe(":3000", router.Handler())
	if err != nil {
		log.Fatal("initialize router fail: ", err)
	}
}
