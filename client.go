package wsclient

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	id           int64
	conn         *websocket.Conn
	lock         *sync.RWMutex
	sends        chan []byte
	events       EventsInterface
	cancel       context.CancelFunc
	isClosed     bool
	writerBuffer uint
	messageType  int
}

type Set func(client *Client)

func Mount(conn *websocket.Conn, sets ...Set) *Client {
	client := &Client{
		conn:         conn,
		lock:         &sync.RWMutex{},
		events:       &Events{},
		writerBuffer: 4,
		messageType:  websocket.TextMessage,
	}
	for _, set := range sets {
		set(client)
	}
	client.sends = make(chan []byte, client.writerBuffer)
	client.run()
	return client
}

func (c *Client) run() {
	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	go c.writer(ctx)
	go c.reader(ctx)
	go c.events.OnConnect(c)
}

func (c *Client) reader(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				c.closer(err)
				return
			}
			c.events.OnMessage(c, msg)
		}
	}
}

func (c *Client) writer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, on := <-c.sends:
			if on {
				if err := c.conn.WriteMessage(c.messageType, msg); err != nil {
					c.closer(err)
					return
				}
			}
		}
	}
}

func (c *Client) closer(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.isClosed {
		c.cancel()
		close(c.sends)
		c.isClosed = true
		c.events.OnClose(c, err)
	}
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) Send(msg []byte) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.isClosed {
		c.sends <- msg
	}
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

func (c *Client) IsClosed() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.isClosed
}

func SetEvents(events EventsInterface) Set {
	return func(client *Client) {
		client.events = events
	}
}

func SetWriterBuffer(num uint) Set {
	return func(client *Client) {
		client.writerBuffer = num
	}
}

func SetMessageType(mt int) Set {
	return func(client *Client) {
		client.messageType = mt
	}
}
