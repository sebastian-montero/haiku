
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
	curl -X POST http://localhost:8080/users -d '{"username": "john_doe2", "email": "john2@example.com", "password": "password123"}' -H "Content-Type: application/json"