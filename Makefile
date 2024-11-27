run-server-local:
	@echo ">>> Setting up Docker..."
	docker compose up -d

	@echo ">>> Installing server dependencies..."
	cd server && go mod download

	@echo ">>> Running server..."
	cd server && go run cmd/server/main.go

docker-down:
	@echo ">>> Stopping docker compose..."
	docker compose down

docker-clean:
	@echo ">>> Cleaning up docker compose..."
	docker compose down --volumes --remove-orphans