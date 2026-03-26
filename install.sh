#!/bin/bash

# dbgod_workOrder_job 一键安装脚本
# 用法: sudo ./install.sh

set -e

APP_NAME="dbgod_workOrder_job"
INSTALL_DIR="/opt/dbgod_workOrder_job"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 root 权限
if [ "$EUID" -ne 0 ]; then
    log_error "请使用 root 权限运行此脚本: sudo ./install.sh"
    exit 1
fi

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

log_info "开始安装 ${APP_NAME}..."

# 停止已有服务
if systemctl is-active --quiet ${APP_NAME}; then
    log_info "停止现有服务..."
    systemctl stop ${APP_NAME}
fi

# 创建安装目录
log_info "创建安装目录: ${INSTALL_DIR}"
mkdir -p ${INSTALL_DIR}/logs

# 复制文件
log_info "复制程序文件..."
cp -f ${SCRIPT_DIR}/${APP_NAME} ${INSTALL_DIR}/
cp -f ${SCRIPT_DIR}/job_config.yml ${INSTALL_DIR}/
chmod +x ${INSTALL_DIR}/${APP_NAME}

# 安装 systemd 服务
log_info "安装 systemd 服务..."
cp -f ${SCRIPT_DIR}/${APP_NAME}.service /etc/systemd/system/
systemctl daemon-reload

# 设置开机自启并启动
log_info "启用开机自启..."
systemctl enable ${APP_NAME}

log_info "启动服务..."
systemctl start ${APP_NAME}

# 检查状态
sleep 2
if systemctl is-active --quiet ${APP_NAME}; then
    log_info "${GREEN}安装成功！${NC}"
    echo ""
    echo "管理命令:"
    echo "  启动: sudo systemctl start ${APP_NAME}"
    echo "  停止: sudo systemctl stop ${APP_NAME}"
    echo "  重启: sudo systemctl restart ${APP_NAME}"
    echo "  状态: sudo systemctl status ${APP_NAME}"
    echo "  日志: sudo journalctl -u ${APP_NAME} -f"
else
    log_error "服务启动失败，请检查日志: journalctl -u ${APP_NAME} -n 50"
    exit 1
fi
