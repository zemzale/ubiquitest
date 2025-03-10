package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/domain/tasks"
)

type connection struct {
	conn *websocket.Conn
	user string
}

type Server struct {
	connections map[string]connection
	rwmutex     *sync.RWMutex

	storeTask  *tasks.Store
	updateTask *tasks.Update
	db         *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{
		connections: make(map[string]connection), rwmutex: &sync.RWMutex{},

		db: db,

		storeTask:  tasks.NewStore(db),
		updateTask: tasks.NewUpdate(db),
	}
}

func (s *Server) Close() {
	for _, conn := range s.connections {
		conn.conn.Close()
	}
}

func (s *Server) TakeConnection(username string, conn *websocket.Conn) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	c := connection{conn: conn, user: username}

	oldConnection, found := s.connections[username]
	if found {
		oldConnection.conn.Close()
	}
	s.connections[username] = c
	go s.handleConnection(c)
}

func (s *Server) handleConnection(c connection) {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message for user '%s' with error '%#v' \n", c.user, err)
			s.removeConnection(c)
			break
		}
		if messageType == websocket.CloseMessage {
			log.Printf("received close for user '%s' \n", c.user)
			s.removeConnection(c)
			break
		}

		if messageType != websocket.TextMessage {
			log.Printf("received unexpected message type %d for user '%s' \n", messageType, c.user)
			continue
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

			task := tasks.Task{
				ID:        taskCreated.Id,
				Title:     taskCreated.Title,
				CreatedBy: taskCreated.CreatedBy,
			}

			if err := s.storeTask.Run(task); err != nil {
				log.Println("failed to store task ", err)
				if replyErr := s.reply(c, EventTypeTaskStoreFailure, err.Error()); replyErr != nil {
					log.Println("failed to reply with error ", replyErr)
				}
			}

			s.broadcast(event, c)

		case EventTypeTaskUpdated:
			log.Println("received task_updated event from user ", c.user)
			taskUpdated, err := event.AsEventTaskUpdated()
			if err != nil {
				log.Println("failed to parse task_updated event ", err, " ", string(message))
				continue
			}

			log.Println(taskUpdated)

			task := tasks.Task{
				ID:        taskUpdated.Id,
				Title:     taskUpdated.Title,
				Completed: taskUpdated.Completed,
			}

			// TODO: Change to have channles based on user id instead of username

			var userID uint
			if err := s.db.Get(&userID, "SELECT id FROM users WHERE username = ?", c.user); err != nil {
				log.Println("failed to get user id ", err)
				continue
			}

			if err := s.updateTask.Run(task, userID); err != nil {
				log.Println("failed to update task ", err)
				if replyErr := s.reply(c, EventTypeTaskStoreFailure, err.Error()); replyErr != nil {
					log.Println("failed to reply with error ", replyErr)
				}
			}

			s.broadcast(event, c)
		case EventTypePing:
			log.Println("received ping from user ", c.user)
			if err := s.reply(c, EventTypePing, nil); err != nil {
				log.Println("failed to reply with error ", err)
			}
		default:
			log.Println("unknown event type ", event.EventType)
			log.Println(string(message))
		}
	}
}

func (s *Server) broadcast(data any, broadcaster connection) {
	for _, conn := range s.connections {
		if conn.user == broadcaster.user {
			continue
		}
		log.Printf("broadcasting to user '%s' \n", conn.user)
		if err := conn.conn.WriteJSON(data); err != nil {
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
	case EventTypePing:
		return c.conn.WriteJSON(Event{EventType: EventTypePong, Data: nil})
	default:
		return fmt.Errorf("unknown replyEventType %s", replyEventType)
	}
}

func (s *Server) removeConnection(c connection) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()

	delete(s.connections, c.user)
}
