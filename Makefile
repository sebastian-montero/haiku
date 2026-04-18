.PHONY: server client install dev clean

install:
	cd server && go mod download
	cd client && npm install

server:
	cd server && go run ./cmd/server

client:
	cd client && npm run dev

# Run both in parallel. Ctrl-C stops everything.
dev:
	@$(MAKE) -j 2 server client

clean:
	rm -rf server/data
