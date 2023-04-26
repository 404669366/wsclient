package wsclient

type EventsInterface interface {
	OnConnect(client *Client)
	OnMessage(client *Client, msg []byte)
	OnClose(client *Client, err error)
}

type Events struct{}

func (e *Events) OnConnect(client *Client) {
}

func (e *Events) OnMessage(client *Client, msg []byte) {
}

func (e *Events) OnClose(client *Client, err error) {
}
