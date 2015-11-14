// Package datastore provides a task.TaskInterface backed by an sql database.
package datastore

import (
	"database/sql"
	"fmt"
	"github.com/rs/xid"

	_ "github.com/lib/pq"

	"github.com/jmank88/todo/task"
)

// The NewDatastore function creates a new task.TaskInterface backed by the postgres database host.
func NewDatastore(host string) (task.TaskInterface, error) {
	return &dataStore{host: host}, nil
}

// A datastore implements the task.TaskInterface, and executes commands against a sql.DB.
type dataStore struct {
	host string

	// Cached database handle. Not to be used directly, use db() instead.
	cachedDB *sql.DB
}

// The db function returns a cached sql.DB, or create and instantiates a new one.
func (d *dataStore) db() (*sql.DB, error) {
	if d.cachedDB != nil {
		return d.cachedDB, nil
	}

	db, err := sql.Open("postgres", d.host)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection for host (%s): %s", d.host, err)
	}
	if err := initDB(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %s", err)
	}

	d.cachedDB = db

	return db, nil
}

// The initDB method initializes the database.
func initDB(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS tasks (id TEXT PRIMARY KEY, title TEXT, content TEXT)"); err != nil {
		return fmt.Errorf("failed to create tasks table: %s", err)
	}
	return nil
}

// The Get method queries the tasks table for a single task with the given id.
func (d *dataStore) Get(id string) (*task.Task, error) {
	db, err := d.db()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT title, content FROM tasks WHERE id = $1", id)
	task := &task.Task{ID: id}
	if err := row.Scan(&(task.Title), &(task.Description)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task %q: %s", id, err)
	}
	return task, nil
}

// The GetAll method queries the tasks table for all tasks.
func (d *dataStore) GetAll() ([]task.Task, error) {
	db, err := d.db()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, title, content FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []task.Task
	for rows.Next() {
		var task task.Task
		err = rows.Scan(&(task.ID), &(task.Title), &(task.Description))
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// The Put method inserts task into the tasks table.
func (d *dataStore) Put(task task.Task) (string, error) {
	db, err := d.db()
	if err != nil {
		return "", err
	}
	if task.ID == "" {
		// No id, so generate a random id.
		task.ID = xid.New().String()
	}
	_, err = db.Exec("INSERT INTO tasks (id, title, content) VALUES ($1, $2, $3)", task.ID, task.Title, task.Description)
	return task.ID, err
}

// The Delete method deletes the task with the given id from the tasks table.
func (d *dataStore) Delete(id string) error {
	db, err := d.db()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
