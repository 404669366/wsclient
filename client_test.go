package wsclient

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
)

func TestBind(t *testing.T) {
	server := gin.Default()

	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server.GET("/ws", func(ctx *gin.Context) {
		conn, _ := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		Mount(conn)
	})

	_ = server.Run(":8080")
}
