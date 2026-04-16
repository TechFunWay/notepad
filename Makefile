VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)
DOCKER_REPO := wycto/notepad
OUTPUT_DIR := release/$(VERSION)

.PHONY: all clean web server build cross-compile docker fpk all-build dev-server dev-web dev

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

fpk:
	@echo "==> Building FPK package..."
	./deploy/scripts/pack-fpk.sh
	@echo "==> FPK package complete"

all-build: cross-compile fpk
	@echo "==> Full build complete"

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
