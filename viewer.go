package main

import (
	"context"
	"fmt"
	"liaozhai/viewer/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, lib.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	r := gin.Default()

	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/main.js", "./public/main.js")
	r.StaticFile("/main.js.map", "./public/main.js.map")
	r.StaticFile("/main.css", "./public/main.css")
	r.StaticFile("/tile.js", "./public/tile.js")
	r.StaticFile("/tile.js.map", "./public/tile.js.map")
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// stream := lib.NewServer()
	// r.Use(stream.ServeHTTP())

	// r.GET("/stream", func(c *gin.Context) {

	// 	c.Stream(func(w io.Writer) bool {
	// 		if msg, ok := <-stream.Message; ok {
	// 			c.Writer.Header().Set("Content-Type", "text/event-stream")
	// 			c.Writer.Header().Set("Cache-Control", "no-cache")
	// 			c.Writer.Header().Set("Connection", "keep-alive")
	// 			c.Writer.Header().Set("Transfer-Encoding", "chunked")
	// 			c.SSEvent("message", msg)
	// 			return true
	// 		}
	// 		return false
	// 	})
	// })

	r.GET("/tile/:layerId/:z/:x/:y", func(c *gin.Context) {
		var t lib.Tile
		if c.ShouldBindUri(&t) == nil {
			data, err := lib.Get(pool, ctx, t.LayerID, t.Z, t.X, t.Y)
			if err == nil {
				c.Data(http.StatusOK, "application/octet-stream", data)
			} else {
				c.Status(http.StatusNoContent)
			}
		}
	})

	// go func() {
	// 	conn, err := pool.Acquire(context.Background())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer conn.Release()

	// 	_, err = conn.Exec(context.Background(), "listen viewer")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for {
	// 		notification, err := conn.Conn().WaitForNotification(context.Background())
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		stream.Message <- fmt.Sprintf(`{ "channel": "%s", "payload": "%s" }`, notification.Channel, notification.Payload)
	// 	}
	// }()

	r.Run(fmt.Sprintf(":%d", lib.Port))

}
