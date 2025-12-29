.PHONY: build run stop logs restart clean lint lint-fix

APP_NAME := template-srv
GOLANGCI_LINT_VERSION := v1.62.2
DOCKER_IMAGE := $(APP_NAME)
DOCKER_CONTAINER := $(APP_NAME)

build:
	docker build -t $(DOCKER_IMAGE) -f ./build/Dockerfile .

run: build
	docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		--env-file ./build/.env \
		$(DOCKER_IMAGE)
	@echo "Container started. Check: curl http://localhost:8080/ping"

stop:
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

logs:
	docker logs -f $(DOCKER_CONTAINER)

restart: stop run

clean:
	docker rmi $(DOCKER_IMAGE) || true

# Linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))
lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

