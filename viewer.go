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
	conn, err := pgxpool.Connect(ctx, lib.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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

	r.GET("/tile/:layerId/:z/:x/:y", func(c *gin.Context) {
		var t lib.Tile
		if c.ShouldBindUri(&t) == nil {
			data, err := lib.Get(conn, ctx, t.LayerID, t.Z, t.X, t.Y)
			if err == nil {
				c.Data(http.StatusOK, "application/octet-stream", data)
			} else {
				c.Status(http.StatusNoContent)
			}
		}
	})

	r.Run(fmt.Sprintf(":%d", lib.Port))
}
