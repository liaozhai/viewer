package lib

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Tile struct {
	LayerID string `uri:"layerId" binding:"uuid"`
	X       int    `uri:"x"`
	Y       int    `uri:"y"`
	Z       int    `uri:"z"`
}

type Event struct {
	Message       chan string
	NewClients    chan chan string
	ClosedClients chan chan string
	TotalClients  map[chan string]bool
}

var (
	dbHost           string
	dbPort           uint
	dbUser           string
	dbPassword       string
	dbName           string
	ConnectionString string
	sql              string
	schema           string
	mvt              string
	Port             uint
)

func init() {
	flag.StringVar(&dbHost, "h", "0.0.0.0", "database host")
	flag.UintVar(&dbPort, "p", 5432, "database port")
	flag.StringVar(&dbUser, "U", "postgres", "database user")
	flag.StringVar(&dbPassword, "W", "postgres", "database password")
	flag.StringVar(&dbName, "d", "", "database name")
	flag.StringVar(&schema, "s", "geo", "schema")
	flag.StringVar(&mvt, "f", "get_mvt_ext", "mvt function")
	flag.UintVar(&Port, "P", 5000, "server port")
	flag.Parse()
	sql = fmt.Sprintf(`SELECT * FROM %s.%s($1, $2, $3, $4)`, schema, mvt)
	ConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}

func Get(conn *pgxpool.Pool, ctx context.Context, layerID string, z int, x int, y int) ([]byte, error) {
	row := conn.QueryRow(ctx, sql, layerID, z, x, y)
	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}
	go event.listen()
	return
}

func (stream *Event) listen() {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

type ClientChan chan string

func (stream *Event) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(ClientChan)
		stream.NewClients <- clientChan

		defer func() {
			stream.ClosedClients <- clientChan
		}()

		go func() {
			<-c.Done()
			stream.ClosedClients <- clientChan
		}()

		c.Next()
	}
}

// func HeadersMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Content-Type", "text/event-stream")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")
// 		c.Writer.Header().Set("Connection", "keep-alive")
// 		c.Writer.Header().Set("Transfer-Encoding", "chunked")
// 		c.Next()
// 	}
// }
