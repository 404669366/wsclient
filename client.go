package wsclient

import (
	"context"
	"github.com/gorilla/websocket"
)

type Client struct {
	id      int64
	conn    *websocket.Conn
	sends   chan []byte
	events  EventsInterface
	cancel  context.CancelFunc
	running bool
}

func New(conn *websocket.Conn, events EventsInterface) *Client {
	return &Client{
		conn:   conn,
		sends:  make(chan []byte, GetMaxWriterBuffer()),
		events: events,
	}
}

func (c *Client) reader() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.Close()
			if !websocket.IsCloseError(err, 1000) {
				go c.events.OnError(c, err)
			}
			go c.events.OnClose(c)
			return
		}
		go c.events.OnMessage(c, msg)
	}
}

func (c *Client) writer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.sends:
			if ok {
				if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					go c.events.OnError(c, err)
				}
			}
		}
	}
}

func (c *Client) Run() {
	if !c.running {
		var ctx context.Context
		ctx, c.cancel = context.WithCancel(context.Background())
		c.running = true
		go c.reader()
		go c.writer(ctx)
		go c.events.OnConnect(c)
	}
}

func (c *Client) Send(msg []byte) {
	c.sends <- msg
}

func (c *Client) Close() {
	c.cancel()
	_ = c.conn.Close()
}

func (c *Client) SetId(id int64) {
	c.id = id
}

func (c *Client) GetId() int64 {
	return c.id
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}
