package router

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Router) WsTodos(writer http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("user")
	if username == "" {
		log.Println("No user name provided")
		return
	}

	log.Println("User name:", username)

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	r.ws.TakeConnection(username, conn)
}
