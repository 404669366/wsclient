# websocket manager

Websocket management center based on gorilla/websocket.

```
go get github.com/404669366/wsManager@latest
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

	server.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	})

	upGrader := websocket.Upgrader{
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
		wsManager.New(conn, &wsManager.Event{}).Run()
	})

	if err := server.Run(":8080"); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
}

```