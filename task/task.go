// Package task provides the Task data structure, and defines TaskInterface for interacting with Tasks.
package task

// A Task is a single todo item.
type Task struct {

	// ID is the unique id of this task.
	ID string `json:"id"`

	// Title is a short description of this task.
	Title string `json:"title"`

	// Description is the main body of this task.
	Description string `json:"description"`
}

// The TaskInterface provides an interface for getting, putting, and deleting tasks.
type TaskInterface interface {

	// The Get method looks up a single task by id.
	Get(id string) (*Task, error)

	// The GetAll method lists all tasks.
	GetAll() ([]Task, error)

	// The Put method adds a single task, and returns the task's id.
	Put(Task) (id string, err error)

	// The Delete method deletes a single task by id.
	Delete(id string) error
}
