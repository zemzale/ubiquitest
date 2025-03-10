package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/zemzale/ubiquitest/domain/tasks"
)

type connection struct {
	conn *websocket.Conn
	user string
}

type Server struct {
	connections []connection

	storeTask *tasks.Store
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

	for i, c := range s.connections {
		if c.user != username {
			continue
		}

		if err := c.conn.Close(); err != nil {
			log.Printf("failed to close connection for user '%s' with error '%s' \n", c.user, err)
		}

		s.connections = slices.Delete(s.connections, i, i+1)
		break
	}

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

			if err := s.storeTask.Run(tasks.Task{
				ID:        taskCreated.Id,
				Title:     taskCreated.Title,
				CreatedBy: c.user,
			}); err != nil {
				log.Println("failed to store task ", err)
				if replyErr := s.reply(c, EventTypeTaskStoreFailure, err.Error()); replyErr != nil {
					log.Println("failed to reply with error ", replyErr)
				}
			}

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

func (s *Server) reply(c connection, replyEventType EventType, replyEventData any) error {
	switch replyEventType {
	case EventTypeTaskStoreFailure:
		e, ok := replyEventData.(EventTaskStoreFailure)
		if !ok {
			return fmt.Errorf("failed to cast replyEventData to EventTaskStoreFailure")
		}
		event, err := FromEventStoreFailure(e)
		if err != nil {
			return fmt.Errorf("failed to create Event from EventTaskStoreFailure: %w", err)
		}

		return c.conn.WriteJSON(event)
	default:
		return fmt.Errorf("unknown replyEventType %s", replyEventType)
	}
}
