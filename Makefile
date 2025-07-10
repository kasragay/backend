APP ?= all
LONG_VERSION ?= v1.0.0
ALL_APPS_NAMES = gateway user
ifeq ($(APP), all)
	APPS = $(ALL_APPS_NAMES)
else
	ifneq ($(APP), $(filter $(APP), $(ALL_APPS_NAMES)))
		$(error APP must be one of $(ALL_APPS_NAMES))
	endif
	APPS = $(APP)
endif

.PHONY: build
build: lint
	@echo "Building ${APPS}"
	@for APP in $(APPS); do \
		go build -o ./bin/$$APP cmd/$$APP/main.go; \
	done

.PHONY: build-settings
build-settings:
	@echo "Building settings..."
	@go build -o bin/settings cmd/settings/main.go

.PHONY: run
run: build
	@./bin/${APP}

.PHONY: createsuperuser
createsuperuser: build-settings
	@./bin/settings createsuperuser

.PHONY: deletesuperuser
deletesuperuserbyid: build-settings
	@./bin/settings deletesuperuser
	
.PHONY: docker-up
docker-up:
	@docker compose -p kg-back up -d

.PHONY: docker-bp
docker-bp: docker-build docker-push

.PHONY: docker-build
docker-build:
	@echo "Building docker images"
	@echo "Building ${APPS}"
	@for APP in $(APPS); do \
		docker build \
			-t ghcr.io/kasragay/backend/$$APP:${LONG_VERSION} \
			-t ghcr.io/kasragay/backend/$$APP:latest \
			-f Dockerfiles/$$APP \
			. ; \
	done

.PHONY: docker-down-app
docker-down-app:
	@echo "Stopping ${APPS}"
	@for APP in $(APPS); do \
		docker compose -p kg-back down $$APP; \
	done

.PHONY: docker-downv-app
docker-downv-app:
	@echo "Stopping and removing volumes ${APPS}"
	@for APP in $(APPS); do \
		docker compose -p kg-back down -v $$APP; \
	done

.PHONY: docker-down
docker-down:
	@echo "Stopping all apps"
	@docker compose -p kg-back down

.PHONY: docker-downv
docker-downv:
	@echo "Stopping all apps and removing volumes"
	@docker compose -p kg-back down -v

.PHONY: docker-push
docker-push:
	@echo "Pushing docker images"
	@echo "Pushing ${APPS}"
	@for APP in $(APPS); do \
		docker push ghcr.io/kasragay/backend/$$APP:${LONG_VERSION}; \
		docker push ghcr.io/kasragay/backend/$$APP:latest; \
	done

.PHONY: docker-logs
docker-logs:
	@docker compose -p kg-back logs ${APP}

.PHONY: test
test:
	@echo "Testing..."
	@go test ./... -v

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f bin/*

.PHONY: playground
playground:
	@echo "Running playground..."
	@go build -o bin/playground cmd/playground/main.go
	@./bin/playground
		
.PHONY: tidy
tidy: clean
	@go mod tidy

.PHONY: lint
lint: tidy
	@gofmt -s -w .
	@golangci-lint run

.PHONY: commit
commit: lint
	@git add .
	@printf "%s\n[app=%s] [version=%s]" "$(M)" "$(APP)" "$(LONG_VERSION)" | git commit -F -
	
.PHONY: push
push: commit
	@git push

.PHONY: compush 
compush: commit push


.PHONY: docker-all
docker-all:
	@echo "Starting docker-all for ${APPS}"
	@$(MAKE) docker-build-parallel
	@$(MAKE) docker-push-parallel
	@$(MAKE) docker-down-app-parallel
	@$(MAKE) docker-up

.PHONY: docker-build-parallel
docker-build-parallel:
	@echo "Building docker images in parallel"
	@for APP in $(APPS); do \
		( \
			echo "Building $$APP"; \
			docker build \
				-t ghcr.io/kasragay/backend/$$APP:${LONG_VERSION} \
				-t ghcr.io/kasragay/backend/$$APP:latest \
				-f Dockerfiles/$$APP \
				. \
		) & \
	done; \
	wait

.PHONY: docker-push-parallel
docker-push-parallel:
	@echo "Pushing docker images in parallel"
	@for APP in $(APPS); do \
		( \
			echo "Pushing $$APP"; \
			docker push ghcr.io/kasragay/backend/$$APP:${LONG_VERSION}; \
			docker push ghcr.io/kasragay/backend/$$APP:latest; \
		) & \
	done; \
	wait

.PHONY: docker-down-app-parallel
docker-down-app-parallel:
	@echo "Stopping ${APPS} in parallel"
	@for APP in $(APPS); do \
		( \
			echo "Stopping $$APP"; \
			docker compose -p kg-back down $$APP; \
		) & \
	done; \
	wait