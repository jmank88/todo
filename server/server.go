// Package server contains an http.Handler for routing http requests to a task.TaskInterface.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jmank88/todo/task"
)

// The NewServer function returns a new server as an http.Handler which routes requests to taskInterface.
func NewServer(taskInterface task.TaskInterface) http.Handler {
	return &server{taskInterface}
}

// Server implements http.Handler, and routes requests to a task.TaskInterface.
type server struct {
	task.TaskInterface
}

// Routes requests based on Method.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	if strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		if id == "" {
			s.getAll(w, r)
		} else {
			s.get(id, w, r)
		}
	case "PUT":
		s.put(w, r)
	case "DELETE":
		if id == "" {
			http.NotFound(w, r)
		} else {
			s.delete(id, w, r)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s not supported", r.Method), http.StatusMethodNotAllowed)
	}
}

// Gets all tasks.
func (s *server) getAll(w http.ResponseWriter, r *http.Request) {
	if tasks, err := s.GetAll(); err != nil {
		http.Error(w, fmt.Sprintf("failed to get all tasks: %s", err), http.StatusInternalServerError)
	} else if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize tasks: %v", tasks), http.StatusInternalServerError)
	}
}

// Gets a single task.
func (s *server) get(id string, w http.ResponseWriter, r *http.Request) {
	if task, err := s.Get(id); err != nil {
		http.Error(w, fmt.Sprintf("failed to get task %s: %s", id, err), http.StatusInternalServerError)
	} else if task == nil {
		http.Error(w, fmt.Sprintf("no task found for id %q", id), http.StatusFound)
	} else if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize task: %v", task), http.StatusInternalServerError)
	}
}

// Puts a task.
func (s *server) put(w http.ResponseWriter, r *http.Request) {
	var task task.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "failed to deserialize task", http.StatusBadRequest)
	} else if id, err := s.Put(task); err != nil {
		http.Error(w, fmt.Sprintf("failed to store task: %v", err), http.StatusInternalServerError)
	} else {
		if _, err := io.WriteString(w, id); err != nil {
			http.Error(w, fmt.Sprintf("failed writing response id %q: ", id, err), http.StatusInternalServerError)
		}
	}
}

// Deletes a task.
func (s *server) delete(id string, w http.ResponseWriter, r *http.Request) {
	if err := s.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete task %s", id), http.StatusInternalServerError)
	}
}
