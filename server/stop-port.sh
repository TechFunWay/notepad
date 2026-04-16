#!/bin/bash

# 停止8904端口的服务
echo "Stopping service on port 8904..."

# 查找占用8904端口的进程ID
PID=$(lsof -t -i :8904)

if [ -z "$PID" ]; then
  echo "No service found on port 8904"
else
  echo "Found process $PID on port 8904"
  kill $PID
  echo "Service stopped successfully"
fi
