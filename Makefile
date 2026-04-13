test:
	go test ./... -coverprofile=coverage.out

install:
	docker-compose up -d
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/vektra/mockery/v2@v2.38.0
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy

build:
	swag init --dir ./src/infrastructure/http/routes/dashboard --parseDependency true 
	go build -o bin/app

build-private:
	swag init --dir ./src/infrastructure/http/routes/private --parseDependency true 
	go build -o bin/app


run:
	$(MAKE) build 
	./bin/app

run-private:
	$(MAKE) build-private 
	./bin/app private

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
