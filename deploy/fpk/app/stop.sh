#!/bin/bash
set -e

# 停止8904端口的服务
echo "Stopping Notepad server..."
PID=$(lsof -t -i :8904 2>/dev/null)
if [ -n "$PID" ]; then
  kill $PID
  echo "Stopped Notepad server (PID: $PID)"
  rm -f /tmp/notepad.pid
else
  echo "No Notepad server running on port 8904"
fi