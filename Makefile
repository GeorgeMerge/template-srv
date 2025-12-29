.PHONY: build run stop logs restart clean

APP_NAME := template-srv
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

