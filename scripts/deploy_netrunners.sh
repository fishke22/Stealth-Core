#!/bin/bash
# 部署 Net-Runners 适配器

echo "🚀 Deploying Net-Runners Adapter..."

# 1. 检查依赖
if ! command -v python3 &> /dev/null; then
    echo "❌ Python3 is required but not installed"
    exit 1
fi

# 2. 检查 Net-Runners 工具是否存在
NETRUNNERS_PATH="/tools/Net-Runners/main.py"
if [ ! -f "$NETRUNNERS_PATH" ]; then
    echo "❌ Net-Runners not found at $NETRUNNERS_PATH"
    echo "Please clone it: git clone https://github.com/fishke22/Net-Runners /tools/Net-Runners"
    exit 1
fi

# 3. 构建 Go 适配器
echo "📦 Building Net-Runners adapter..."
cd adapters/netrunners
go build -o netrunners-adapter .

# 4. 创建系统服务
echo "📋 Creating system service..."
cat > /etc/systemd/system/netrunners-adapter.service << EOF
[Unit]
Description=Net-Runners Adapter Service
After=network.target

[Service]
Type=simple
ExecStart=$(pwd)/netrunners-adapter
WorkingDirectory=$(pwd)
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 5. 启动服务
systemctl daemon-reload
systemctl enable netrunners-adapter
systemctl start netrunners-adapter

echo "✅ Net-Runners Adapter deployed successfully!"
echo "📊 Check status: systemctl status netrunners-adapter"
echo "📝 Logs: journalctl -u netrunners-adapter -f"
