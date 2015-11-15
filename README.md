# Todo
A simple golang todo list server.

The *todo* binary runs a todo list server backed by postgres, exposing a json http api. The *cli* binary provides a
command line interface for running commands against a todo server.


## Make
The Makefile handles building, unit tests, and integration tests.

The default target runs units tests and builds the *server* and *cli* binaries.
```
make
```

The *integration-tests* target runs integration tests against postgres via docker compose.
```
make integration-tests
```


## Services
The *todo* server exposes json services for getting, putting, and deleting tasks.
```
./todo --help

Usage of ./todo:
  -host string
    	postgres host (default "postgres://postgres:postgres@localhost:5432?sslmode=disable")
  -port string
    	port to serve (default "8080")
```
Serves port 8080 by default, and assumes localhost postgres.

### Get All
```
GET <host>/
```
Gets all tasks. Returns a json list of task objects.

### Get
```
GET <host>/<id>
```
Gets a single task with the given id. Returns a json task object.

### Put
```
PUT <host>/
```
Puts a task. Accepts a json task object. Returns either the provided task id, or a uid if none was provided.

### Delete
```
DELETE <host>/<id>
```
Deletes the task with the given id.


## Command Line Interface
A simple command line interface is included as an alternative to hitting the http services directly, and also serves as
an example usage of the client package.
```
./cli --help

Usage of ./cli:
  -X string
    	method to execute. required. must be one of 'GET', 'PUT', or 'DEL'
  -host string
    	http task host to connect to (default "http://localhost:8080")
  -id string
    	task id. required for delete, optional for get
```

### Get All
```
./cli -X GET
```
Gets all tasks. Prints a go slice of task structs.

### Get
```
./cli -X GET -id <id>
```
Gets a single task by id. Prints the go task struct.

### PUT
```
./cli -X PUT -id <id> -title <title> -description <description>
```
Puts a task. Accepts optional id, title, and description flags. Prints the task id.

### DEL
```
./cli -X DEL -id <id>
```
Deletes the task with the given id.


## Running locally
Docker compose can be used to run postgres and the todo server locally. The server runs on localhost:8080 and may be hit
directly or with the cli.

Both postgres and the todo server can be started in the background by running:
```
docker-compose up -d
```
At this point, the todo server is available at localhost:8080.

Shut them down with:
```
docker-compose stop
```

Here is an example session using the cli:
```
docker-compose up -d

./cli -X PUT -id 1 -title "Shopping List" -description "milk, eggs, bread"
> put task "1"

./cli -X GET -id 1
> &task.Task{ID:"1", Title:"Shopping List", Description:"milk, eggs, bread"}

./cli -X PUT -title "Call Mom" -description "Call mom @5:00pm"
> put task "VkeEoUn1XQAB1bov"

./cli -X GET
> []task.Task{task.Task{ID:"1", Title:"Shopping List", Description:"milk, eggs, bread"}, task.Task{ID:"VkgI6xJGrAABnEXc", Title:"Call Mom", Description:"Call mom @5:00pm"}}

./cli -X DEL -id 1
> deleted task "1"

docker-compose stop
```