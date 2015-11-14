// Package client provides a remote http client implementation of task.TaskInterface.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmank88/todo/task"
)

const defaultHost = "localhost:8080"

// The NewClient function creates a new remote client implementing task.TaskInterface.
// The client will use "localhost:8080" and http.DefaultClient, unless configured different with options.
func NewClient(options ...Option) task.TaskInterface {
	c := &client{
		httpClient: http.DefaultClient,
		host:       defaultHost,
	}
	for _, o := range options {
		o(c)
	}
	return c
}

// An Option is a functional option for configuring a client.
type Option func(*client)

// The Host function returns an Option for configuring a client's host.
func Host(host string) Option {
	return func(c *client) {
		c.host = host
	}
}

// The HTTPClient function returns an Option for configuring a client's http.Client.
func HTTPClient(hc *http.Client) Option {
	return func(c *client) {
		c.httpClient = hc
	}
}

// A client implements task.TaskInterface, and executes commands against a remote host with json over http.
type client struct {
	httpClient *http.Client
	host       string
}

func (c *client) Get(id string) (*task.Task, error) {
	if id == "" {
		return nil, errors.New("no id specified")
	}
	resp, err := c.httpClient.Get(c.host + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task %q: %s", id, err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var task task.Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			return nil, fmt.Errorf("failed to deserialize task %q: %s", id, err)
		}
		return &task, nil
	default:
		errStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response from get task %q attempt: %s", id, err)
		}
		return nil, fmt.Errorf("failed to get task %q: %s", id, errStr)
	}

}

func (c *client) GetAll() ([]task.Task, error) {
	resp, err := c.httpClient.Get(c.host)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response from get all attempt: %s", err)
		}
		return nil, fmt.Errorf("failed to get all tasks: %s", errStr)
	}

	var tasks []task.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to deserialize tasks: %s", err)
	}
	return tasks, nil
}

func (c *client) Put(task task.Task) (string, error) {
	bs, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to serialize task %v: %s", task, err)
	}
	req, err := http.NewRequest("PUT", c.host, bytes.NewReader(bs))
	if err != nil {
		return "", fmt.Errorf("failed to create http request for task %v: %s", task, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute put request for task %v: %s", task, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read error response from task %v put attempt: %s", task, err)
		}
		return "", fmt.Errorf("failed to put task %v: %s", task, errStr)
	}

	id, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response after putting task %v: %s", task, err)
	}
	return string(id), nil
}

func (c *client) Delete(id string) error {
	req, err := http.NewRequest("DELETE", c.host+"/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create http request for task %q: %s", id, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute delete request for task %s: %s", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response from task %q delete attempt: %s", id, err)
		}
		return fmt.Errorf("failed to delete task %q: %s", id, errStr)
	}

	return nil
}
