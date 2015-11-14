// +build integration

package datastore

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	"github.com/jmank88/todo/task"
)

var host = flag.String("host", "postgres://postgres:postgres@localhost:5432?sslmode=disable", "database host")

func init() {
	flag.Parse()
}

// The fixture function creates a task.TaskInterface backed by a clean database.
func fixture(t *testing.T) task.TaskInterface {
	d := &dataStore{host: *host}
	db, err := d.db()
	if err != nil {
		t.Fatal("failed to connect to database: ", err)
	}
	if err := clear(db); err != nil {
		t.Fatal("failed to clear database: ", err)
	}
	return d
}

// The clear function clears the database by truncating the tasks table.
func clear(db *sql.DB) error {
	if _, err := db.Exec("TRUNCATE TABLE tasks"); err != nil {
		return fmt.Errorf("failed to truncate tasks table: ", err)
	}
	return nil
}

// Tests getting a task after putting it.
func TestPutGet(t *testing.T) {
	taskInterface := fixture(t)

	task := task.Task{
		ID:          "testId",
		Title:       "testTitle",
		Description: "testDescription",
	}
	if id, err := taskInterface.Put(task); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if id != task.ID {
		t.Fatal("expected %q but got %q", task.ID, id)
	}

	if got, err := taskInterface.Get(task.ID); err != nil {
		t.Fatal("unexpected error getting task: ", err)
	} else if *got != task {
		t.Fatal("expected %v got %v", task, got)
	}
}

// Tests deleting a task after putting it.
func TestPutDelete(t *testing.T) {
	taskInterface := fixture(t)

	task := task.Task{
		ID:          "testId",
		Title:       "testTitle",
		Description: "testDescription",
	}
	if id, err := taskInterface.Put(task); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if id != task.ID {
		t.Fatal("expected %q but got %q", task.ID, id)
	}

	if err := taskInterface.Delete(task.ID); err != nil {
		t.Fatal("unexpected error deleting task: ", err)
	}

	if got, err := taskInterface.Get(task.ID); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if got != nil {
		t.Fatalf("expected no result but got %v", got)
	}
}

// Tests getting all tasks after putting them.
func TestPutGetAll(t *testing.T) {
	taskInterface := fixture(t)

	tasks := indexByID([]task.Task{
		task.Task{
			ID:          "1",
			Title:       "task 1",
			Description: "description 1",
		},
		task.Task{
			ID:          "2",
			Title:       "task 2",
			Description: "description 2",
		},
		task.Task{
			ID:          "3",
			Title:       "task 3",
			Description: "description 3",
		},
	})

	for _, task := range tasks {
		if id, err := taskInterface.Put(task); err != nil {
			t.Fatal("unexpected error: ", err)
		} else if id != task.ID {
			t.Fatal("expected %q but got %q", task.ID, id)
		}
	}

	if gotSlice, err := taskInterface.GetAll(); err != nil {
		t.Fatal("unexpected error: ", err)
	} else {
		gotMap := indexByID(gotSlice)
		if len(gotMap) != len(tasks) {
			t.Fatal("expected equal maps\n expected %v\n but got %v", tasks, gotMap)
		}
		for id, task := range tasks {
			if gotTask, ok := gotMap[id]; !ok {
				t.Fatal("expected returned map to contain %v", task)
			} else if task != gotTask {
				t.Fatal("expected %v for id %s but got %v", task, id, gotTask)
			}
		}
	}
}

// Tests getting all from an empty database.
func TestGetAllNone(t *testing.T) {
	taskInterface := fixture(t)

	if tasks, err := taskInterface.GetAll(); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if len(tasks) > 0 {
		t.Fatalf("expected empty task set, but got %v", tasks)
	}
}

// Tests getting a non existent task.
func TestGetNonExistent(t *testing.T) {
	taskInterface := fixture(t)

	if task, err := taskInterface.Get("testId"); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if task != nil {
		t.Fatalf("expected nil task but got %v", task)
	}
}

func indexByID(tasks []task.Task) map[string]task.Task {
	taskMap := make(map[string]task.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}
	return taskMap
}
