package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
)

type Server struct {
	connections map[string]*Client

	// We are using a single channel for client changes to avoid locking
	// the map and also to avoid race conditions with register/unregister channels
	clientChangeChan chan *clientChange

	storeTask  *tasks.Store
	updateTask *tasks.Update
	userFind   *users.FindByUsername
}

type clientChange struct {
	client *Client
	action clientchangeAction
}

type clientchangeAction int

const (
	add clientchangeAction = iota
	remove
)

func NewServer(storeTask *tasks.Store, updateTask *tasks.Update, findUserByUsername *users.FindByUsername) *Server {
	return &Server{
		connections:      make(map[string]*Client),
		clientChangeChan: make(chan *clientChange),

		storeTask:  storeTask,
		updateTask: updateTask,
		userFind:   findUserByUsername,
	}
}

func (s *Server) Close() {
	close(s.clientChangeChan)

	for _, client := range s.connections {
		client.Close()
	}
}

func (s *Server) Run(ctx context.Context) {
	go s.handleClients(ctx)
}

func (s *Server) handleClients(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-s.clientChangeChan:
			switch client.action {
			case add:
				s.registerClient(client.client)
			case remove:
				s.unregisterClient(client.client)
			default:
				log.Println("unknown client action ", client.action)
			}
		}
	}
}

func (s *Server) registerClient(client *Client) {
	log.Println("registering client ", client.user.Username)

	if oldClient, ok := s.connections[client.user.Username]; ok {
		oldClient.Close()
	}

	s.connections[client.user.Username] = client

	go s.handleConnection(client)
}

func (s *Server) unregisterClient(client *Client) {
	log.Println("unregistering client ", client.user.Username)

	if _, ok := s.connections[client.user.Username]; !ok {
		return
	}

	delete(s.connections, client.user.Username)
}

func (s *Server) TakeConnection(username string, conn *websocket.Conn) {
	user, err := s.userFind.Run(username)
	if err != nil {
		log.Println("failed to find user ", err)
		return
	}

	c := NewClient(conn, user)
	s.clientChangeChan <- &clientChange{client: c, action: add}

	go s.handleConnection(c)
}

func (s *Server) handleConnection(c *Client) {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message for user '%s' with error '%#v' \n", c.user.Username, err)
			s.clientChangeChan <- &clientChange{client: c, action: remove}
			return
		}
		if messageType == websocket.CloseMessage {
			log.Printf("received close for user '%s' \n", c.user.Username)
			s.clientChangeChan <- &clientChange{client: c, action: remove}
			return
		}

		if messageType != websocket.TextMessage {
			log.Printf("received unexpected message type %d for user '%s' \n", messageType, c.user.Username)
			continue
		}

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("failed to unmarshal event %s \n", err.Error())
			continue
		}

		log.Println("got a new event : ", event.EventType)

		switch event.EventType {
		case EventTypeTaskCreated:
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
				ParentID:  taskCreated.ParentId,
				Cost:      taskCreated.Cost,
			}

			if err := s.storeTask.Run(task); err != nil {
				log.Println("failed to store task ", err)
				if replyErr := s.reply(c, EventTypeTaskStoreFailure, err.Error()); replyErr != nil {
					log.Println("failed to reply with error ", replyErr)
				}
				continue
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
				Cost:      taskUpdated.Cost,
			}

			if err := s.updateTask.Run(task, c.user.ID); err != nil {
				log.Println("failed to update task ", err)
				if replyErr := s.reply(c, EventTypeTaskStoreFailure, err.Error()); replyErr != nil {
					log.Println("failed to reply with error ", replyErr)
				}
				continue
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

func (s *Server) broadcast(data any, broadcaster *Client) {
	for _, conn := range s.connections {
		if conn.user.ID == broadcaster.user.ID {
			continue
		}
		log.Printf("broadcasting to user '%#v' \n", conn.user)
		if err := conn.conn.WriteJSON(data); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) reply(c *Client, replyEventType EventType, replyEventData any) error {
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
