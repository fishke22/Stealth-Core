#!/bin/bash
# 部署 Redteam Tools 适配器

echo "🚀 Deploying Redteam Tools Adapter..."

# 1. 检查依赖
if ! command -v go &> /dev/null; then
    echo "❌ Go is required but not installed. Please install Go (e.g., sudo apt install golang)"
    exit 1
fi

# 2. 检查 Hexstrike-redteam 工具是否存在
REDTEAM_BASE_PATH="/tools/Hexstrike-redteam"
if [ ! -d "$REDTEAM_BASE_PATH" ]; then
    echo "❌ Hexstrike-redteam base directory not found at $REDTEAM_BASE_PATH"
    echo "Please clone it: git clone https://github.com/Yenn503/Hexstrike-redteam $REDTEAM_BASE_PATH"
    exit 1
fi

# 3. 构建 Go 适配器
echo "📦 Building Redteam Tools adapter..."
cd /tmp/chimera_project/adapters/redteam_tools # 确保在正确的目录下
go build -o redteam-tools-adapter . || { echo "❌ Go build failed for Redteam Tools adapter."; exit 1; }

# 4. 创建系统服务
echo "📋 Creating system service..."
cat > /etc/systemd/system/redteam-tools-adapter.service << EOF
[Unit]
Description=Redteam Tools Adapter Service
After=network.target

[Service]
Type=simple
ExecStart=$(pwd)/redteam-tools-adapter
WorkingDirectory=$(pwd)
Environment="REDTEAM_TOOLS_BASE_PATH=$REDTEAM_BASE_PATH"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 5. 启动服务
systemctl daemon-reload
systemctl enable redteam-tools-adapter
systemctl start redteam-tools-adapter

echo "✅ Redteam Tools Adapter deployed successfully!"
echo "📊 Check status: systemctl status redteam-tools-adapter"
echo "📝 Logs: journalctl -u redteam-tools-adapter -f"
