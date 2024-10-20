
fmt:
	@echo "Applying format..."
	gofmt -w .

run:
	@echo "Running server..."
	go run cmd/server/main.go

build: fmt
	@echo "Building server..."
	go build -o bin/server cmd/server/main.go

docker-up:
	@echo "Starting docker compose..."
	docker compose up -d

query:
	@echo "Running query..."
	curl -X POST http://localhost:8080/users -d '{"username": "john_doe2", "email": "john3@example.com", "password": "password123"}' -H "Content-Type: application/json"

query-nbs:
	@echo "Running query..."
	curl -X POST http://localhost:8080/notebooks -d '{"title": "test", "owner_id": 2}' -H "Content-Type: application/json"

query-get:
	@echo "Running query..."
	curl -X GET http://localhost:8080/users/12

query-delete:
	@echo "Running query..."
	curl -X DELETE http://localhost:8080/users/12

query-update:
	@echo "Running query..."
	curl -X PUT http://localhost:8080/users -d '{"id":12, "username": "john_doe5", "email": "updateda@example", "password": "password123"}' -H "Content-Type: application/json"