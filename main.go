package main

import (
	"flag"
	"fmt"
	"os"

	"hireme-api/config"
	"hireme-api/db"
	"hireme-api/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.DbConnect()
	db.CacheConnect()

	// dataseeding
	// dataseed()

	// start server
	server.Init()
}

func dataseed() {

}
