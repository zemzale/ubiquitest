package ws

import (
	"github.com/gorilla/websocket"
	"github.com/zemzale/ubiquitest/domain/users"
)

type Client struct {
	conn *websocket.Conn
	user users.User
}

func NewClient(conn *websocket.Conn, user users.User) *Client {
	return &Client{conn: conn, user: user}
}

func (c *Client) Close() {
	c.conn.Close()
}
