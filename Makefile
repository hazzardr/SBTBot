include .env
export

PROJECT_NAME := smart-bitches-trashy-bot
EXEC_NAME := sbtbot
SSH_USER := ansible
DEPLOY_TARGET_IP := 100.100.77.57
DB_URL := data/sbtb.db

.PHONY: help ## print this
help:
	@echo ""
	@echo "$(PROJECT_NAME) Development CLI"
	@echo ""
	@echo "Usage:"
	@echo "  make <command>"
	@echo ""
	@echo "Commands:"
	@grep '^.PHONY: ' Makefile | sed 's/.PHONY: //' | awk '{split($$0,a," ## "); printf "  \033[34m%0-10s\033[0m %s\n", a[1], a[2]}'

.PHONY: init ## initialize the project
init:
	go run . init

.PHONY: run ## Run the project
run:
	@go run . bot serve --token=$$SBTB_DISC_TOKEN

.PHONY: sync ## Sync discord commands with gateway
sync:
	@go run . bot sync --token=$$SBTB_DISC_TOKEN

.PHONY: clean ## delete generated code
clean:
	rm -rf generated

.PHONY: build ## builds the project
build:
	go build -ldflags='-s' -o=./bin/${PROJECT_NAME} .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/${PROJECT_NAME} .

.PHONY: lint ## run golangci-lint
lint:
	golangci-lint run

.PHONY: test ## run tests
test:
	go test -v ./...

.PHONY: fmt ## format the project
fmt:
	go fmt ./...

.PHONY: generate ## generate database stubs
generate:
	sqlc generate -f sqlc.yaml

.PHONY: db/migration/status ## get the status of the db migrations
db/migration/status:
	goose sqlite3 $(DB_URL) -dir sql/migrations status

.PHONY: db/migrate ## run database migrations
db/migrate:
	goose sqlite3 $(DB_URL) -dir sql/migrations up

.PHONY: production/connect ## connects to production deployment server
production/connect:
	ssh ${SSH_USER}@${DEPLOY_TARGET_IP}

