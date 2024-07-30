package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Hana-ame/neo-moonchan/api"
	"github.com/Hana-ame/neo-moonchan/psql"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// connect to database
	connStr := os.Getenv("DATABASE_URL")
	psql.Connect(connStr)
	// loop
	for {
		func() {
			// atoshimatsu
			defer func() {
				e := recover()
				if e != nil {
					fmt.Println(e)
				}
			}()

			err := api.Main()
			fmt.Printf("%v", err)
		}()

		time.Sleep(5 * time.Second)
	}
}
