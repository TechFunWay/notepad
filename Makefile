VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)
DOCKER_REPO := techfunways/notepad
OUTPUT_DIR := release/$(VERSION)
DEV_DIR := dev
DEV_BINARY := $(DEV_DIR)/notepad
DEV_STAGE := $(DEV_DIR)/.next
PORT ?= 8904

.PHONY: all clean web server build cross-compile docker build-fpk dev-server dev-web dev serve stop

all: clean build

web:
	@echo "==> Building frontend..."
	cd web && npm ci --silent && npm run build

server: web
	@echo "==> Copying frontend assets..."
	rm -rf server/static/dist
	cp -r web/dist server/static/dist
	@echo "==> Building backend..."
	mkdir -p $(OUTPUT_DIR)
	cd server && CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(OUTPUT_DIR)/notepad ./main.go

build: server
	@echo "==> Build complete: $(OUTPUT_DIR)/notepad"

cross-compile: web
	@echo "==> Copying frontend assets..."
	rm -rf server/static/dist
	cp -r web/dist server/static/dist
	@echo "==> Cross-compiling..."
	mkdir -p $(OUTPUT_DIR)
	cd server && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ../$(OUTPUT_DIR)/notepad_linux_amd64 ./main.go && \
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ../$(OUTPUT_DIR)/notepad_linux_arm64 ./main.go && \
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ../$(OUTPUT_DIR)/notepad_darwin_amd64 ./main.go && \
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ../$(OUTPUT_DIR)/notepad_darwin_arm64 ./main.go
	@echo "==> Cross-compile complete"

docker:
	@echo "==> Building Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_REPO):$(VERSION) -t $(DOCKER_REPO):latest --load .
	@echo "==> Docker image built: $(DOCKER_REPO):$(VERSION)"

clean:
	rm -rf $(OUTPUT_DIR) server/static/dist web/dist

dev-server:
	@echo "==> Stopping any existing service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || true
	@echo "==> Starting backend server..."
	cd server && go run main.go

dev-web:
	cd web && npm run dev

dev:
	@echo "==> [1/5] Building frontend..."
	cd web && npm run build
	@echo "==> [2/5] Copying frontend assets..."
	rm -rf server/static/dist
	cp -r web/dist server/static/dist
	rm -rf $(DEV_STAGE)
	mkdir -p $(DEV_STAGE)
	cp -r web/dist $(DEV_STAGE)/www
	cp scripts/start-notepad.sh $(DEV_STAGE)/start.sh
	chmod +x $(DEV_STAGE)/start.sh
	@echo "==> [3/5] Building backend binary..."
	cd server && CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(DEV_STAGE)/notepad ./main.go
	@echo "==> [4/5] Stopping any existing service on port $(PORT)..."
	@PIDS="$$(lsof -t -iTCP:$(PORT) -sTCP:LISTEN 2>/dev/null)"; \
	if [ -n "$$PIDS" ]; then \
		echo "==> Stopping process(es): $$PIDS"; \
		kill $$PIDS 2>/dev/null || true; \
		for i in 1 2 3 4 5; do \
			if ! lsof -t -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1; then break; fi; \
			sleep 1; \
		done; \
		REMAINING="$$(lsof -t -iTCP:$(PORT) -sTCP:LISTEN 2>/dev/null)"; \
		if [ -n "$$REMAINING" ]; then \
			echo "==> Process did not stop gracefully; forcing stop: $$REMAINING"; \
			kill -9 $$REMAINING 2>/dev/null || true; \
			sleep 1; \
		fi; \
		if lsof -t -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1; then \
			echo "==> Unable to release port $(PORT); aborting restart"; \
			exit 1; \
		fi; \
	else \
		echo "==> Port $(PORT) is available"; \
	fi
	@echo "==> Activating packaged development files..."
	rm -rf $(DEV_DIR)/www
	mv $(DEV_STAGE)/www $(DEV_DIR)/www
	mv -f $(DEV_STAGE)/notepad $(DEV_BINARY)
	mv -f $(DEV_STAGE)/start.sh $(DEV_DIR)/start.sh
	rmdir $(DEV_STAGE)
	@echo "==> [5/5] Starting packaged development service..."
	@echo "==> Runtime directory: $(DEV_DIR)"
	@echo "==> Visit: http://localhost:$(PORT)"
	cd $(DEV_DIR) && PORT=$(PORT) ./start.sh

build-fpk:
	@echo "==> Building FPK package..."
	./scripts/build-all.sh
	./scripts/build-fnpack.sh
	@echo "==> FPK package complete"

serve:
	@echo "==> Building frontend..."
	cd web && npm run build
	@echo "==> Copying frontend assets to server/static/dist..."
	rm -rf server/static/dist
	cp -r web/dist server/static/dist
	@echo "==> Stopping any existing service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || true
	@echo "==> Building backend binary..."
	cd server && CGO_ENABLED=0 go build -o /tmp/notepad-server .
	@echo "==> Starting backend in background..."
	@nohup /tmp/notepad-server </dev/null >/tmp/notepad-server.log 2>&1 & disown
	@sleep 2
	@echo "==> Server started. PID: $$(lsof -t -i :8904 2>/dev/null | head -1)"
	@echo "==> Logs: /tmp/notepad-server.log"
	@echo "==> Visit http://localhost:8904"

stop:
	@echo "==> Stopping service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || echo "==> No service running on 8904"
