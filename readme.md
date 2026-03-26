## 编译
GOOS=linux GOARCH=amd64 go build -o dbgod_workOrder_job dbgod_workOrder_job

## 部署文件

| 文件 | 说明 |
|------|------|
| `dbgod_workOrder_job` | 编译后的二进制文件 |
| `job_config.yml` | 配置文件 |
| `dbgod_workOrder_job.service` | systemd 服务文件 |
| `install.sh` | 一键安装脚本 |
| `uninstall.sh` | 卸载脚本 |

## 部署方式

将以上文件上传到服务器任意目录后执行：

```bash
# 添加执行权限
chmod +x install.sh

# 一键安装
./install.sh
```

## 管理命令

```bash
# 启动服务
systemctl start dbgod_workOrder_job

# 停止服务
systemctl stop dbgod_workOrder_job

# 重启服务
systemctl restart dbgod_workOrder_job

# 查看状态
systemctl status dbgod_workOrder_job

# 查看系统日志
journalctl -u dbgod_workOrder_job -f

# 查看应用日志
tail -f /opt/dbgod_workOrder_job/logs/run.log
```

## 卸载

```bash
./uninstall.sh
```