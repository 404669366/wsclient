# websocket client

websocket client for [gorilla/websocket](https://github.com/gorilla/websocket).

```
go get github.com/404669366/wsclient@latest
```

```golang
package main

import (
	"fmt"
	"github.com/404669366/wsclient"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	server := gin.Default()

	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server.GET("/ws", func(ctx *gin.Context) {
		conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
			return
		}
		wsclient.New(conn, &wsclient.Events{}).Run()
	})

	if err := server.Run(":8080"); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
}

```