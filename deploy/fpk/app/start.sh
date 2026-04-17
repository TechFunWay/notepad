#!/bin/bash
set -e

# 创建数据目录
APP_DATA_DIR="/vol*/@appdata/notepad/data"
mkdir -p "$APP_DATA_DIR"

# 停止8904端口的服务
echo "Stopping any existing service on port 8904..."
PID=$(lsof -t -i :8904 2>/dev/null)
if [ -n "$PID" ]; then
  kill $PID
  echo "Stopped existing service (PID: $PID)"
  sleep 1
fi

# 设置环境变量
export PORT=8904
export DB_PATH="$APP_DATA_DIR/notepad.db"
export JWT_SECRET=""
export TZ=Asia/Shanghai

# 启动服务
echo "Starting Notepad server..."
cd "$(dirname "$0")/.."
./notepad-server > /dev/null 2>&1 &

# 等待服务启动
sleep 3

# 检查服务是否启动成功
PID=$(lsof -t -i :8904 2>/dev/null)
if [ -n "$PID" ]; then
  echo "Notepad server started successfully (PID: $PID)"
  echo $PID > /tmp/notepad.pid
else
  echo "Failed to start Notepad server"
  exit 1
fi