#!/bin/bash

APP_NAME="f02-admin"
APP_PATH="./bin"
CONFIG_PATH="./configs"

start() {
    # 检查应用是否已运行
    PID=$(pgrep -f "$APP_PATH/$APP_NAME")
    if [ -n "$PID" ]; then
        echo "$APP_NAME is already running with PID $PID."
        exit 1
    fi

    echo "Starting $APP_NAME..."
    nohup "$APP_PATH/$APP_NAME" --conf="$CONFIG_PATH" --env=pro >/dev/null 2>&1 &
    echo "$APP_NAME started."
}

stop() {
    # 获取运行中的应用 PID
    PID=$(pgrep -f "$APP_PATH/$APP_NAME")
    if [ -z "$PID" ]; then
        echo "$APP_NAME is not running."
        exit 1
    fi

    echo "Stopping $APP_NAME with PID $PID..."
    kill "$PID"
    echo "$APP_NAME stopped."
}

restart() {
    stop
    start
}

status() {
    # 检查应用是否正在运行
    PID=$(pgrep -f "$APP_PATH/$APP_NAME")
    if [ -n "$PID" ]; then
        echo "$APP_NAME is running with PID $PID."
    else
        echo "$APP_NAME is not running."
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

exit 0
