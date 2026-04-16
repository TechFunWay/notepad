# Makefile for Notepad server

.PHONY: dev stop

dev:
	@echo "Stopping any existing service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || true
	@echo "Starting Notepad server..."
	@go run main.go

stop:
	@echo "Stopping service on port 8904..."
	@lsof -t -i :8904 | xargs kill -9 2>/dev/null || echo "No service found on port 8904"
