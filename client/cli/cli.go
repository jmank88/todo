package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jmank88/todo/client"
	"github.com/jmank88/todo/task"
)

var (
	host   = flag.String("host", "http://localhost:8080", "http task host to connect to")
	method = flag.String("X", "", "method to execute. required. must be one of 'GET', 'PUT', or 'DEL'")
	id     = flag.String("id", "", "task id. required for delete, optional for get")
)

func main() {
	flag.Parse()

	if *method == "" {
		log.Fatal("no method (-X) specified. must be one of 'GET', 'PUT', or 'DEL'")
	}

	taskClient := client.NewClient(client.Host(*host))

	switch strings.ToUpper(*method) {
	case "GET":
		if *id == "" {
			tasks, err := taskClient.GetAll()
			if err != nil {
				log.Fatalf("failed to get all tasks: %s", err)
			}
			log.Println(tasks)
		} else {
			task, err := taskClient.Get(*id)
			if err != nil {
				log.Fatalf("failed to get task %q: %s", *id, err)
			}
			if task == nil {
				log.Fatalf("no task found for id %q", *id)
			}
			log.Println(task)
		}
	case "DEL":
		if *id == "" {
			log.Fatal("no id specified for delete")
		}
		if err := taskClient.Delete(*id); err != nil {
			log.Fatalf("failed to delete task %q: %s", *id, err)
		}
		log.Printf("deleted task %q\n", *id)
	case "PUT":
		var task task.Task
		if err := json.NewDecoder(os.Stdin).Decode(&task); err != nil {
			log.Fatalf("failed to deserialize json task: %s", err)
		}
		id, err := taskClient.Put(task)
		if err != nil {
			log.Fatalf("failed to put task: %s", err)
		}
		log.Printf("put task %q\n", id)
	default:
		log.Fatalf("unregonized method %q. must be one of 'GET', 'PUT', or 'DEL'", *method)
	}
}
