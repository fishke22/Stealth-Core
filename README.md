# Chimera - Stealth-Core 滲透測試框架

Chimera 是一個基於 gRPC 的模組化滲透測試框架，專為 OpenClaw 安全生態系統設計。

## 🚀 功能特色

- **模組化架構**: 透過 gRPC 協調多個滲透工具適配器
- **分散式執行**: 各工具以 sidecar 模式運行，透過 orchestrator 協調
- **協議緩衝區**: 使用 Protocol Buffers 定義標準化通訊協議
- **OpenClaw 整合**: 無縫整合到 OpenClaw 安全代理生態系統

## 🛠️ 組件架構

### 核心組件
- **Orchestrator**: 主協調服務（port 50000）
- **Protocol Buffers**: 標準化 gRPC 服務定義
- **工具適配器**: 各類滲透工具整合

### 可用適配器
- `hexstrike` - HexStrike AI 攻擊工具
- `netrunners` - 網路滲透工具  
- `redteam_tools` - 紅隊工具整合
- `thesilencer` - 隱蔽行動工具
- `villager_ai` - 社交工程工具

## 📦 安裝部署

### 二進位安裝
```bash
# 下載最新版本
wget https://github.com/OpenClaw-Security/Stealth-Core/releases/latest/download/chimera-linux-amd64

# 安裝到系統路徑
sudo mv chimera-linux-amd64 /usr/local/bin/chimera
sudo chmod +x /usr/local/bin/chimera
```

### 源碼編譯
```bash
# 克隆倉庫
git clone https://github.com/OpenClaw-Security/Stealth-Core.git
cd Stealth-Core

# 安裝依賴
go mod tidy

# 編譯二進位
go build -o chimera ./cmd/orchestrator/
```

## 🚀 快速開始

### 啟動 Orchestrator
```bash
chimera
# 預設監聽 :50000
```

### 環境變數配置
```bash
export ORCHESTRATOR_LISTEN_ADDR=":50000"
chimera
```

### 啟動工具適配器
每個適配器需要單獨編譯和運行：
```bash
# 編譯 redteam_tools 適配器
cd adapters/redteam_tools/
go build -o redteam-tools-adapter .

# 運行適配器
./redteam-tools-adapter
```

## 📖 開發指南

### 專案結構
```
Stealth-Core/
├── cmd/orchestrator/     # 主協調服務
├── core/orchestrator/    # 協調器核心邏輯
├── pkg/proto/           # gRPC protobuf 定義
├── adapters/            # 工具適配器
│   ├── hexstrike/       # HexStrike 適配器
│   ├── netrunners/      # 網路滲透適配器
│   ├── redteam_tools/   # 紅隊工具適配器
│   ├── thesilencer/     # 隱蔽行動適配器
│   └── villager_ai/     # 社交工程適配器
└── go.mod              # Go 模組定義
```

### Protocol Buffers 開發
```bash
# 編譯 .proto 文件
protoc \
  --proto_path=./pkg/proto \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  pkg/proto/*.proto
```

## 🔧 配置說明

### Orchestrator 配置
透過環境變數配置：
- `ORCHESTRATOR_LISTEN_ADDR`: gRPC 監聽地址（預設: :50000）
- `LOG_LEVEL`: 日誌級別（debug, info, warn, error）

### 適配器配置
每個適配器有自己的配置檔案，通常位於：
- `adapters/[adapter-name]/config.yaml`
- 環境變數前綴：`[ADAPTER]_`

## 🤝 貢獻指南

1. Fork 本倉庫
2. 創建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 📜 許可協議

本專案採用 MIT 許可證 - 詳見 [LICENSE](LICENSE) 文件。

## 🐛 問題回報

請使用 GitHub Issues 回報問題：
https://github.com/OpenClaw-Security/Stealth-Core/issues

## 🙏 致謝

- OpenClaw 團隊 - 提供核心框架支援
- 所有貢獻者 - 感謝您的代碼貢獻和問題回報

---

**Chimera** - 讓滲透測試更智能、更協調