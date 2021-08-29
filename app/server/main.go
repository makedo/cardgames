package main

import (
	controllers "api/app/http/controllers"
	"net/http"
)

func initRoutes() {
	http.HandleFunc("/deck", controllers.GetDeck)
	http.HandleFunc("/durak", controllers.GetState)
	http.HandleFunc("/ws", controllers.Websocket())
}

func main() {
	initRoutes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
