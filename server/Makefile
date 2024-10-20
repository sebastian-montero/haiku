
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