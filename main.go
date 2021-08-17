package main

import (
	"log"

	"github.com/khanhvtn/netevent-go/api"
	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/services"
)

func main() {
	//Checking database connection
	if !database.ConnectionOK() {
		log.Fatal("Not connected to DB")
		return
	}

	//Create services
	di, err := services.New()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//start API
	api.Init(di)
}
