package lib

import (
	"context"
	"flag"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Tile struct {
	LayerID string `uri:"layerId" binding:"uuid"`
	X       int    `uri:"x"`
	Y       int    `uri:"y"`
	Z       int    `uri:"z"`
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
