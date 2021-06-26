package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ash3798/Social-Network/config"
)

func StartServer() {
	api := NewAPI()

	mux := http.NewServeMux()

	//fill endpoint and handler info into mux
	for _, endp := range api.endpoints {
		mux.HandleFunc(endp.path, endp.handler)
	}

	addr := fmt.Sprintf(":%d", config.Manager.AppPort)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting http server at localhost%s \n", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println("HTTP server stopped . Error : ", err.Error())
	}
}
