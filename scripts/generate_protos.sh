#!/bin/bash
# 生成 Go proto 代码

PROTO_DIR="/tmp/chimera_project/pkg/proto"

# 确保输出目录存在
mkdir -p "${PROTO_DIR}/core"
mkdir -p "${PROTO_DIR}/hexstrike"
mkdir -p "${PROTO_DIR}/thesilencer"
mkdir -p "${PROTO_DIR}/netrunners"
mkdir -p "${PROTO_DIR}/villager_ai"
mkdir -p "${PROTO_DIR}/redteam_tools"

# 生成 core.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/core" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/core" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/core.proto"

# 生成 hexstrike.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/hexstrike" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/hexstrike" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/hexstrike.proto"

# 生成 thesilencer.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/thesilencer" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/thesilencer" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/thesilencer.proto"

# 生成 netrunners.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/netrunners" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/netrunners" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/netrunners.proto"

# 生成 villager_ai.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/villager_ai" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/villager_ai" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/villager_ai.proto"

# 生成 redteam_tools.proto
protoc --proto_path="${PROTO_DIR}" \
       --go_out="${PROTO_DIR}/redteam_tools" --go_opt=paths=source_relative \
       --go-grpc_out="${PROTO_DIR}/redteam_tools" --go-grpc_opt=paths=source_relative \
       "${PROTO_DIR}/redteam_tools.proto"

echo "✅ Go proto code generated successfully!"
