package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type connection struct {
	conn *websocket.Conn
	user string
}

type Server struct {
	connections []connection
}

func NewServer() *Server {
	return &Server{connections: make([]connection, 0)}
}

func (s *Server) Close() {
	for _, conn := range s.connections {
		conn.conn.Close()
	}
}

func (s *Server) TakeConnection(username string, conn *websocket.Conn) {
	c := connection{conn: conn, user: username}
	s.connections = append(s.connections, c)
	go s.handleConnection(c)
}

func (s *Server) handleConnection(c connection) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message for user '%s' with error '%s' \n", c.user, err)
			break
		}

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("failed to unmarshal event %s \n", err.Error())
			continue
		}

		switch event.EventType {
		case "task_created":
			log.Println("received task_created event from user ", c.user)
			taskCreated, err := event.AsEventTaskCreated()
			if err != nil {
				log.Println("failed to parse task_created event ", err, " ", string(message))
				continue
			}

			log.Println(taskCreated)

			// s.store.tasks.Add(task)

			s.broadcast(taskCreated, c)
		default:
			log.Println("unknown event type ", event.EventType)
			log.Println(string(message))
		}
	}
}

func (s *Server) broadcast(taskCreated EventTaskCreated, broadcaster connection) {
	for _, conn := range s.connections {
		if conn.user == broadcaster.user {
			continue
		}
		log.Printf("broadcasting to user '%s' \n", conn.user)
		if err := conn.conn.WriteJSON(taskCreated); err != nil {
			log.Println(err)
		}
	}
}
