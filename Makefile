VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)
DOCKER_REPO := techfunways/notepad
OUTPUT_DIR := release/$(VERSION)

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
	@echo "==> Building frontend..."
	cd web && npm run build
	@echo "==> Copying frontend assets..."
	rm -rf server/static/dist
	cp -r web/dist server/static/dist
	@echo "==> Stopping any existing service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || true
	@echo "==> Starting backend server..."
	cd server && go run main.go

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
