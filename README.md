# Todo
A simple golang todo list server.


## Make
The Makefile handles building, unit tests, and integration tests.

Standard *make* runs units tests and builds the server and cli binaries.

*make integration-tests* runs integration tests against postgres via docker compose.


## Services
The todo server exposes json services for getting, putting, and deleting tasks.

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
Puts a task. Accepts a json task object. Returns the either the provided task id, or a uid if none was provided.

### Delete
```
DELETE <host>/<id>
```
Deletes the task with the given id.


## Command Line Interface
A simple command line interface is included as an alternative to hitting the http services directly, and as a reference
usage of the client package.
```
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
./cli -X PUT <<< '<json task object>'
```
Puts a task. Accepts a json task object via stdin. Prints the task id.

### DEL
```
./cli -X DEL -id <id>
```
Deletes the task with the given id.


## Running locally
Docker compose can be used to run postgres and the todo server locally. The server runs on localhost:8080 and may be hit
directly or with the cli.

```
docker-compose up -d

./cli -X PUT <<< '{"id":"1", "title":"Shopping List", "description": "milk, eggs, bread"}'
> put task "1"

./cli -X GET -id 1
> &{1 Shopping List milk, eggs, bread}

./cli -X PUT <<< '{"title":"Call Mom", "description": "Call mom @5:00pm"}'
> put task "VkeEoUn1XQAB1bov"

./cli -X GET
> [{1 Shopping List milk, eggs, bread} {VkeEoUn1XQAB1bov Call Mom Call mom @5:00pm}]

./cli -X DEL -id 1
> deleted task "1"

docker-compose stop
```