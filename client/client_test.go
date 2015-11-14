package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmank88/todo/task"
)

// Tests a get request.
func TestGet(t *testing.T) {
	expected := task.Task{
		ID:          "test task",
		Title:       "test title",
		Description: "test description",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		expectedPath := "/" + expected.ID
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %q but got %q", expectedPath, r.URL.Path)
		}

		if err := json.NewEncoder(w).Encode(expected); err != nil {
			t.Fatal("unexpected error encoding json: ", err)
		}
	}))
	defer ts.Close()

	ti := NewClient(Host(ts.URL))

	if got, err := ti.Get(expected.ID); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if got == nil {
		t.Fatalf("expected %v but got nil", expected)
	} else if *got != expected {
		t.Fatalf("expected %v but got %v", expected, got)
	}
}

// Tests a get all request.
func TestGetAll(t *testing.T) {
	expected := []task.Task{
		{
			ID:          "test id",
			Title:       "test title",
			Description: "test description",
		},
		{
			ID:          "test id 2",
			Title:       "test title",
			Description: "test description",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		const expectedPath = "/"
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %q but got %q", expectedPath, r.URL.Path)
		}

		if err := json.NewEncoder(w).Encode(expected); err != nil {
			t.Fatal("unexpected error encoding json: ", err)
		}
	}))
	defer ts.Close()

	ti := NewClient(Host(ts.URL))

	if got, err := ti.GetAll(); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if len(expected) != len(got) {
		t.Fatalf("epected %v but got %v", expected, got)
	} else {
		expectedMap := indexByID(expected)
		gotMap := indexByID(got)
		for id, task := range expectedMap {
			if gotTask, ok := gotMap[id]; !ok || gotTask != task {
				t.Fatalf("epected %v but got %v", expected, got)
			}
		}
	}
}

// Tests a put request.
func TestPut(t *testing.T) {
	testTask := task.Task{
		ID:          "test task",
		Title:       "test title",
		Description: "test description",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		const expectedPath = "/"
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %q but got %q", expectedPath, r.URL.Path)
		}

		var task task.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			t.Fatal("unexpected error decoding json: ", err)
		} else if task != testTask {
			t.Fatalf("expected %v but got %v", testTask, task)
		}

		if _, err := io.WriteString(w, task.ID); err != nil {
			t.Fatal("unexpected error writing response")
		}
	}))
	defer ts.Close()

	ti := NewClient(Host(ts.URL))

	if got, err := ti.Put(testTask); err != nil {
		t.Fatal("unexpected error: ", err)
	} else if got != testTask.ID {
		t.Fatalf("expected id %q but got %q", testTask.ID, got)
	}
}

// Tests a delete request.
func TestDelete(t *testing.T) {
	const testId = "id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		expectedPath := "/" + testId
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %q but got %q", expectedPath, r.URL.Path)
		}
	}))
	defer ts.Close()

	ti := NewClient(Host(ts.URL))

	if err := ti.Delete(testId); err != nil {
		t.Fatal("unexpected error: ", err)
	}
}

func indexByID(tasks []task.Task) map[string]task.Task {
	taskMap := make(map[string]task.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}
	return taskMap
}
