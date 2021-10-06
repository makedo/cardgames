package main

import (
	"cardgames/app/controllers"
	"net/http"
)

func initRoutes() {
	http.HandleFunc("/ws", controllers.Websocket())
}

func main() {
	initRoutes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
