package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jmank88/todo/datastore"
	"github.com/jmank88/todo/server"
)

var host = flag.String("host", "postgres://postgres:postgres@localhost:5432?sslmode=disable", "postgres host")
var port = flag.String("port", "8080", "port to serve")

func main() {
	flag.Parse()

	log.Println("using host: ", *host)

	taskInterface, err := datastore.NewDatastore(*host)
	if err != nil {
		log.Fatal("failed to create datastore: ", err)
	}

	log.Fatal(http.ListenAndServe(":" + *port, server.NewServer(taskInterface)))
}
