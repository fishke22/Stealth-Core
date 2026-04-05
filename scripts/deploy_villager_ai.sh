#!/bin/bash
# 部署 Villager-AI 适配器

echo "🚀 Deploying Villager-AI Adapter..."

# 1. 检查依赖
if ! command -v go &> /dev/null; then
    echo "❌ Go is required but not installed. Please install Go (e.g., sudo apt install golang)"
    exit 1
fi

# 2. 检查依赖的适配器是否已部署并运行
echo "🔍 Checking dependent adapters status..."
systemctl is-active --quiet hexstrike-adapter || { echo "❌ Hexstrike Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet thesilencer-adapter || { echo "❌ TheSilencer Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet netrunners-adapter || { echo "❌ NetRunners Adapter is not running. Please deploy and start it first."; exit 1; }
echo "✅ All dependent adapters are running."

# 3. 构建 Go 适配器
echo "📦 Building Villager-AI adapter..."
cd /tmp/chimera_project/adapters/villager_ai # 确保在正确的目录下
go build -o villager-ai-adapter . || { echo "❌ Go build failed for Villager-AI adapter."; exit 1; }

# 4. 创建系统服务
echo "📋 Creating system service..."
cat > /etc/systemd/system/villager-ai-adapter.service << EOF
[Unit]
Description=Villager-AI Adapter Service
After=network.target hexstrike-adapter.service thesilencer-adapter.service netrunners-adapter.service

[Service]
Type=simple
ExecStart=$(pwd)/villager-ai-adapter
WorkingDirectory=$(pwd)
Environment="HEXSTRIKE_ADAPTER_ADDR=localhost:50051"
Environment="THESILENCER_ADAPTER_ADDR=localhost:50052"
Environment="NETRUNNERS_ADAPTER_ADDR=localhost:50053"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 5. 启动服务
systemctl daemon-reload
systemctl enable villager-ai-adapter
systemctl start villager-ai-adapter

echo "✅ Villager-AI Adapter deployed successfully!"
echo "📊 Check status: systemctl status villager-ai-adapter"
echo "📝 Logs: journalctl -u villager-ai-adapter -f"
