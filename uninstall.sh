#!/bin/bash

# dbgod_workOrder_job 卸载脚本

set -e
APP_NAME="dbgod_workOrder_job"
INSTALL_DIR="/opt/dbgod_workOrder_job"

if [ "$EUID" -ne 0 ]; then
    echo "请使用 root 权限运行: sudo ./uninstall.sh"
    exit 1
fi

echo "停止服务..."
systemctl stop ${APP_NAME} 2>/dev/null || true
systemctl disable ${APP_NAME} 2>/dev/null || true

echo "删除服务文件..."
rm -f /etc/systemd/system/${APP_NAME}.service
systemctl daemon-reload

echo "删除安装目录..."
rm -rf ${INSTALL_DIR}

echo "卸载完成！"
