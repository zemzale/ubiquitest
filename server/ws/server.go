package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Server struct {
	connections []*websocket.Conn
}

func NewServer() *Server {
	return &Server{connections: make([]*websocket.Conn, 0)}
}

func (s *Server) Close() {
	for _, conn := range s.connections {
		conn.Close()
	}
}

func (s *Server) TakeConnection(conn *websocket.Conn) {
	s.connections = append(s.connections, conn)
	go s.handleConnection(conn)
}

func (s *Server) handleConnection(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message ", err)
			break
		}

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println("failed to unmarshal event ", err)
			continue
		}

		switch event.EventType {
		case "task_created":
			log.Println("received task_created event")
			taskCreated, err := event.AsEventTaskCreated()
			if err != nil {
				log.Println("failed to parse task_created event ", err, " ", string(message))
				continue
			}

			log.Println(taskCreated)

			s.broadcast(taskCreated)
		default:
			log.Println("unknown event type ", event.EventType)
			log.Println(string(message))
		}
	}
}

func (s *Server) broadcast(taskCreated EventTaskCreated) {
	for _, conn := range s.connections {
		if err := conn.WriteJSON(taskCreated); err != nil {
			log.Println(err)
		}
	}
}
