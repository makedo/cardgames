package main

import (
	controllers "cardgames/app/http/controllers"
	"net/http"
)

func initRoutes() {
	http.HandleFunc("/deck", controllers.GetDeck)
	http.HandleFunc("/ws", controllers.Websocket())
}

func main() {
	initRoutes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
