# WriteOut Server 

The WriteOut backend server powers the core functionality of the WriteOut application. It allows users to do real time streaming while writing projects. Users can create, open, and write in notebooks, with live updates available to other users via WebSockets. The server is built in Go and handles user sessions, notebook management, and real-time communication.

## Prerequisites
- Go 1.16 or higher
- Docker (for containerized environments)
- PostgreSQL (configured in the dev.yaml file for the database connection)

## Makefile Commands

The project uses a Makefile to simplify common tasks. Here’s a list of available commands:


### fmt
This command applies Go formatting to the source files using gofmt.
Usage: `make fmt`

### run
Runs the Go server from the main entry point at cmd/server/main.go.
Usage: `make run`

### build
Formats the code and builds the server binary into the bin/ directory.
Usage: `make build`

### build-docs
Generates Go documentation and writes the index file to ./docs/index.idx.
Usage: `make build-docs`

### run-docs
Runs a local godoc server that serves the Go documentation at http://localhost:6060/.
Usage: `make run-docs`

### docker-up
Starts the Docker containers as defined in your docker-compose configuration.
Usage: `make docker-up`