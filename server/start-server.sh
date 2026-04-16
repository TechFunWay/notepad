#!/bin/bash

# 停止8904端口的服务
echo "Stopping any existing service on port 8904..."
PID=$(lsof -t -i :8904)
if [ -n "$PID" ]; then
  kill $PID
  echo "Stopped existing service (PID: $PID)"
  sleep 1
fi

# 启动服务
echo "Starting Notepad server..."
go run main.go
