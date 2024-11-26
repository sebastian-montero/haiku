
fmt:
	@echo ">>> Applying format..."
	gofmt -w .

run: fmt
	@echo ">>> Running server..."
	go run cmd/server/main.go

build: fmt
	@echo ">>> Building server..."
	go build -o bin/server cmd/server/main.go

build-docs: 
	@echo ">>> Generating docs..."
	GO111MODULE=on
	godoc -write_index -index_files=./docs/index.idx

run-docs:
	@echo ">>> Running docs..."
	godoc -http=:6060 -index_files=./docs/index.idx