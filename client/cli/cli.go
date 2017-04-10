package main

import (
	"flag"
	"log"
	"strings"

	"github.com/jmank88/todo/client"
	"github.com/jmank88/todo/task"
)

var (
	host        = flag.String("host", "http://localhost:8080", "http task host to connect to")
	method      = flag.String("X", "", "method to execute. required. must be one of 'GET', 'PUT', or 'DEL'")
	id          = flag.String("id", "", "task id. required for delete, optional for get and put")
	title       = flag.String("title", "", "task title. only used for put")
	description = flag.String("description", "", "task description. only used for put")
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
			log.Printf("%#v\n", tasks)
		} else {
			task, err := taskClient.Get(*id)
			if err != nil {
				log.Fatalf("failed to get task %q: %s", *id, err)
			}
			if task == nil {
				log.Fatalf("no task found for id %q", *id)
			}
			log.Printf("%#v\n", task)
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
		id, err := taskClient.Put(task.Task{
			ID:          *id,
			Title:       *title,
			Description: *description,
		})
		if err != nil {
			log.Fatalf("failed to put task: %s", err)
		}
		log.Printf("put task %q\n", id)
	default:
		log.Fatalf("unregonized method %q. must be one of 'GET', 'PUT', or 'DEL'", *method)
	}
}
