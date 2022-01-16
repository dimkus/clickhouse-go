package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const query = `
SELECT
	number
	, randomString(25)
	, array(1, 2, 3, 4, 5)
	, now()
FROM system.numbers LIMIT 1000000
`

func benchmark(conn clickhouse.Conn) error {
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			col1 uint64
			col2 string
			col3 []uint8
			col4 time.Time
		)
		if err := rows.Scan(&col1, &col2, &col3, &col4); err != nil {
			return err
		}
	}
	return nil
}
func main() {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	if err := benchmark(conn); err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Since(start))
}