#!/bin/bash
# 部署 Chimera Orchestrator

echo "🚀 Deploying Chimera Orchestrator..."

# 1. 检查 Go 依赖
if ! command -v go &> /dev/null; then
    echo "❌ Go is required but not installed. Please install Go (e.g., sudo apt install golang)"
    exit 1
fi

# 2. 检查所有依赖的适配器服务是否已部署并运行
echo "🔍 Checking dependent adapter services status..."
systemctl is-active --quiet hexstrike-adapter || { echo "❌ Hexstrike Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet thesilencer-adapter || { echo "❌ TheSilencer Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet netrunners-adapter || { echo "❌ NetRunners Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet villager-ai-adapter || { echo "❌ Villager-AI Adapter is not running. Please deploy and start it first."; exit 1; }
systemctl is-active --quiet redteam-tools-adapter || { echo "❌ Redteam Tools Adapter is not running. Please deploy and start it first."; exit 1; }
echo "✅ All dependent adapters are running."

# 3. 构建 Orchestrator 可执行文件
echo "📦 Building Orchestrator..."
cd /tmp/chimera_project/cmd/orchestrator # 确保在正确的目录下
go mod init github.com/OpenClaw-Security/Stealth-Core/cmd/orchestrator # 初始化 Go 模块
go mod tidy # 同步依赖
go build -o chimera-orchestrator . || { echo "❌ Go build failed for Orchestrator."; exit 1; }

# 4. 创建系统服务
echo "📋 Creating system service..."
cat > /etc/systemd/system/chimera-orchestrator.service << EOF
[Unit]
Description=Chimera Orchestrator Service
After=network.target hexstrike-adapter.service thesilencer-adapter.service netrunners-adapter.service villager-ai-adapter.service redteam-tools-adapter.service

[Service]
Type=simple
ExecStart=$(pwd)/chimera-orchestrator
WorkingDirectory=$(pwd)
Environment="ORCHESTRATOR_LISTEN_ADDR=:50000"
Environment="HEXSTRIKE_ADAPTER_ADDR=localhost:50051"
Environment="THESILENCER_ADAPTER_ADDR=localhost:50052"
Environment="NETRUNNERS_ADAPTER_ADDR=localhost:50053"
Environment="VILLAGER_AI_ADAPTER_ADDR=localhost:50054"
Environment="REDTEAM_TOOLS_ADAPTER_ADDR=localhost:50055"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 5. 启动服务
systemctl daemon-reload
systemctl enable chimera-orchestrator
systemctl start chimera-orchestrator

echo "✅ Chimera Orchestrator deployed successfully!"
echo "📊 Check status: systemctl status chimera-orchestrator"
echo "📝 Logs: journalctl -u chimera-orchestrator -f"
