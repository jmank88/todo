package server

import (
	"net/http/httptest"
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/jmank88/todo/task"
)

// Tests a get request.
func TestGet(t *testing.T) {
	expected := task.Task{
		ID:          "test task",
		Title:       "test title",
		Description: "test description",
	}

	ts := httptest.NewServer(NewServer(&mockTaskInterface{
		get: func(id string) (*task.Task, error) {
			if id != expected.ID {
				t.Fatalf("expected %q but got %q", expected.ID, id)
			}
			return &expected, nil
		},
	}))
	defer ts.Close()

	if resp, err := http.Get(ts.URL + "/" + expected.ID); err != nil {
		t.Fatal("unexpected error sending request: ", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		var got task.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatal("unexpected error decoding response: ", err)
		} else if got != expected {
			t.Fatalf("expected %v but got %v", expected, got)
		}
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

	ts := httptest.NewServer(NewServer(&mockTaskInterface{
		getAll: func() ([]task.Task, error) {
			return expected, nil
		},
	}))
	defer ts.Close()

	if resp, err := http.Get(ts.URL); err != nil {
		t.Fatal("unexpected error sending request: ", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		var got []task.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatal("unexpected error decoding response: ", err)
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
}

// Tests a put request.
func TestPut(t *testing.T) {
	testTask := task.Task{
		ID:          "test task",
		Title:       "test title",
		Description: "test description",
	}

	ts := httptest.NewServer(NewServer(&mockTaskInterface{
		put: func(task task.Task) (string, error) {
			if task != testTask {
				t.Fatalf("expected %v but got %v", testTask, task)
			}
			return testTask.ID, nil
		},
	}))
	defer ts.Close()

	bs, err := json.Marshal(testTask)
	if err != nil {
		t.Fatal("unexpected error marshaling json: ", err)
	}
	req, err := http.NewRequest("PUT", ts.URL, bytes.NewReader(bs))
	if err != nil {
		t.Fatal("unexpected error building request: ", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, resp.StatusCode)
	}

	if got, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("failed to read response after putting task %v: %s", testTask, err)
	} else if string(got) != testTask.ID {
		t.Fatalf("expected id %q but got %q", testTask.ID, got)
	}
}

// Tests a delete request.
func TestDelete(t *testing.T) {
	const testId = "id"
	ts := httptest.NewServer(NewServer(&mockTaskInterface{
		delete: func(id string) error {
			if id != testId {
				t.Fatalf("expected %q but got %q", testId, id)
			}
			return nil
		},
	}))
	defer ts.Close()

	req, err := http.NewRequest("DELETE", ts.URL+"/"+testId, nil)
	if err != nil {
		t.Fatal("unexpected error building request: ", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("unexpcted error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

type mockTaskInterface struct {
	get    func(string) (*task.Task, error)
	getAll func() ([]task.Task, error)
	put    func(task.Task) (string, error)
	delete func(string) error
}

func (m *mockTaskInterface) Get(id string) (*task.Task, error) {
	return m.get(id)
}

func (m *mockTaskInterface) GetAll() ([]task.Task, error) {
	return m.getAll()
}

func (m *mockTaskInterface) Put(task task.Task) (string, error) {
	return m.put(task)
}

func (m *mockTaskInterface) Delete(id string) error {
	return m.delete(id)
}

func indexByID(tasks []task.Task) map[string]task.Task {
	taskMap := make(map[string]task.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}
	return taskMap
}
