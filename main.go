package main

import (
	"log"

	"github.com/khanhvtn/netevent-go/api"
	"github.com/khanhvtn/netevent-go/database"
)

func main() {
	//Checking database connection
	if !database.ConnectionOK() {
		log.Fatal("Not connected to DB")
		return
	}
	//start API
	api.Init()
}
