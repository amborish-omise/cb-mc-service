# Application settings
APP_NAME := mastercom-service
APP_FILE := cmd/%/main.go
ENV_FILE := .env

# Docker related settings
BUILDKITE_DIR  := .buildkite
SCRIPT_DIR     := $(BUILDKITE_DIR)/scripts
DOCKER_COMPOSE := docker compose
DOCKER_BUILD   := docker build
DOCKER_TOOLS   := $(DOCKER_COMPOSE) -f $(BUILDKITE_DIR)/docker-compose.tools.yml
DOCKER_TEST    := $(DOCKER_COMPOSE) -f $(BUILDKITE_DIR)/docker-compose.yml
DOCKER_ORPHANS := --remove-orphans

.PHONY: env test

run-%: verify-% build-dev
	$(DOCKER_COMPOSE) up $(DOCKER_ORPHANS) app-$*

debug-%: build-debug env
	$(DOCKER_COMPOSE) run --rm --service-ports -v  ${SCRIPT_DIR}:/app/${SCRIPT_DIR} --entrypoint="${SCRIPT_DIR}/debug.sh cmd/$*/main.go" app-$*

check: build-test
	$(DOCKER_TEST) run --rm --entrypoint $(SCRIPT_DIR)/check.sh app-test

test: build-test
	$(DOCKER_TEST) run --rm --entrypoint $(SCRIPT_DIR)/test.sh app-test

# This is for internal use. Please use `run-*` target instead.
build-%: Dockerfile env
	$(DOCKER_BUILD) -t $(APP_NAME) --target $* .

release-%: verify-%
	$(DOCKER_BUILD) -t $(APP_NAME)-$* --build-arg $* --target release .

verify-%:
	[ ! -f $(subst %,$*,$(APP_FILE)) ] && echo "error: $(subst %,$*,$(APP_FILE)) not exists" && exit 1 || true

env:
	[ ! -f $(ENV_FILE) ] && cat .env.example | sed 's/APP_NAME=.*/APP_NAME=$(APP_NAME)/g' > .env || true

clean:
	$(DOCKER_COMPOSE) down $(DOCKER_ORPHANS)

gen-proto: build-tools
	$(DOCKER_TOOLS) run --rm tools \
	-I ./specs/grpc \
	-I /usr/local/include \
	--go_out ./specs/grpc --go_opt paths=source_relative \
	--go-grpc_out ./specs/grpc --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./specs/grpc --grpc-gateway_opt paths=source_relative \
	./specs/grpc/**/*.proto

gen-openapi: build-tools
	$(DOCKER_TOOLS) run --rm tools \
	-I ./specs/grpc \
	-I /usr/local/include \
	--openapiv2_out ./specs/grpc \
	--openapiv2_opt logtostderr=true \
	./specs/grpc/**/*.proto
