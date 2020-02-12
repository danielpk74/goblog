package api

import (
	"api/router"
	"config"
	"fmt"
	"log"
	"net/http"
)

// Run starts a server
func Run() {
	config.Load()
	fmt.Printf("Server is running [::]:%v", config.PORT)
	listen(config.PORT)
}

func listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
