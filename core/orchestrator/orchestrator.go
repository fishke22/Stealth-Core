package orchestrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	corepb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core" // 核心 proto
	hexstrikepb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike"
	netrunnerspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
	redteam_toolspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools"
	thesilencerpb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	villager_aipb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai"
	"github.com/OpenClaw-Security/Stealth-Core/pkg/workflow" // 工作流解析器
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Orchestrator 结构体：现在包含所有适配器客户端
type Orchestrator struct {
	corepb.UnimplementedOrchestratorServer // Changed from proto.UnimplementedOrchestratorServer
	
	// 适配器客户端
	hexstrikeClient   hexstrikepb.HexstrikeAdapterClient
	thesilencerClient thesilencerpb.TheSilencerAdapterClient
	netrunnersClient  netrunnerspb.NetRunnersAdapterClient
	villagerAIClient  villager_aipb.VillagerAIAdapterClient
	redteamToolsClient redteam_toolspb.RedteamToolsAdapterClient

	workflowParser *workflow.WorkflowParser
	mu             sync.Mutex // 用于保护共享资源，例如连接状态
	connections    []*grpc.ClientConn // 存储所有 gRPC 连接，以便在关闭时清理
}

// NewOrchestrator 创建并初始化 Orchestrator 实例
func NewOrchestrator() (*Orchestrator, error) {
	orch := &Orchestrator{
		workflowParser: workflow.NewWorkflowParser(),
	}

	// 从环境变量或默认值获取适配器地址
	hexstrikeAddr := os.Getenv("HEXSTRIKE_ADAPTER_ADDR")
	if hexstrikeAddr == "" {
		hexstrikeAddr = "localhost:50051"
	}
	thesilencerAddr := os.Getenv("THESILENCER_ADAPTER_ADDR")
	if thesilencerAddr == "" {
		thesilencerAddr = "localhost:50052"
	}
	netrunnersAddr := os.Getenv("NETRUNNERS_ADAPTER_ADDR")
	if netrunnersAddr == "" {
		netrunnersAddr = "localhost:50053"
	}
	villagerAIAddr := os.Getenv("VILLAGER_AI_ADAPTER_ADDR")
	if villagerAIAddr == "" {
		villagerAIAddr = "localhost:50054"
	}
	redteamToolsAddr := os.Getenv("REDTEAM_TOOLS_ADAPTER_ADDR")
	if redteamToolsAddr == "" {
		redteamToolsAddr = "localhost:50055"
	}

	// 连接到 Hexstrike Adapter
	err := orch.connectAdapter("hexstrike", hexstrikeAddr, func(conn *grpc.ClientConn) {
		orch.hexstrikeClient = hexstrikepb.NewHexstrikeAdapterClient(conn)
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to Hexstrike adapter at %s: %v", hexstrikeAddr, err)
	}

	// 连接到 TheSilencer Adapter
	err = orch.connectAdapter("thesilencer", thesilencerAddr, func(conn *grpc.ClientConn) {
		orch.thesilencerClient = thesilencerpb.NewTheSilencerAdapterClient(conn)
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to TheSilencer adapter at %s: %v", thesilencerAddr, err)
	}

	// 连接到 NetRunners Adapter
	err = orch.connectAdapter("netrunners", netrunnersAddr, func(conn *grpc.ClientConn) {
		orch.netrunnersClient = netrunnerspb.NewNetRunnersAdapterClient(conn)
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to NetRunners adapter at %s: %v", netrunnersAddr, err)
	}

	// 连接到 Villager-AI Adapter
	err = orch.connectAdapter("villager_ai", villagerAIAddr, func(conn *grpc.ClientConn) {
		orch.villagerAIClient = villager_aipb.NewVillagerAIAdapterClient(conn)
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to Villager-AI adapter at %s: %v", villagerAIAddr, err)
	}

	// 连接到 RedteamTools Adapter
	err = orch.connectAdapter("redteam_tools", redteamToolsAddr, func(conn *grpc.ClientConn) {
		orch.redteamToolsClient = redteam_toolspb.NewRedteamToolsAdapterClient(conn)
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to RedteamTools adapter at %s: %v", redteamToolsAddr, err)
	}

	return orch, nil
}

// connectAdapter 辅助函数，用于连接到单个适配器
func (o *Orchestrator) connectAdapter(name, addr string, clientInitFunc func(*grpc.ClientConn)) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(5*time.Second)) // Changed grpc.WithInsecure() to insecure.NewCredentials()
	if err != nil {
		return fmt.Errorf("could not connect to %s adapter at %s: %w", name, addr, err)
	}
	o.connections = append(o.connections, conn)
	clientInitFunc(conn)
	log.Printf("Successfully connected to %s adapter at %s", name, addr)
	return nil
}

// Close gracefully closes all gRPC connections.
func (o *Orchestrator) Close() {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, conn := range o.connections {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}
	o.connections = nil
}

// ExecuteWorkflow 执行一个完整的工作流
func (o *Orchestrator) ExecuteWorkflow(req *corepb.WorkflowRequest, stream corepb.Orchestrator_ExecuteWorkflowServer) error { // Changed proto.WorkflowRequest, proto.Orchestrator_ExecuteWorkflowServer to corepb.WorkflowRequest, corepb.Orchestrator_ExecuteWorkflowServer
	log.Printf("Received workflow execution request: %s", req.Name)

	// 将 gRPC 请求转换为工作流定义
	workflowDef := &workflow.WorkflowDefinition{
		Name:        req.Name,
		Operations:  make([]workflow.OperationDefinition, len(req.Operations)),
		Variables:   req.Variables,
	}
	for i, opReq := range req.Operations {
		workflowDef.Operations[i] = workflow.OperationDefinition{
			ID:          opReq.OperationId,
			Tool:        opReq.Tool,
			Command:     opReq.Command,
			Parameters:  opReq.Params,
			// DependsOn, Timeout, Retry 等字段需要从 req.Operations 进一步映射，这里简化
		}
	}

	return o.executeWorkflowDefinition(workflowDef, stream)
}

// executeWorkflowDefinition 实际执行工作流的核心逻辑
func (o *Orchestrator) executeWorkflowDefinition(workflowDef *workflow.WorkflowDefinition, stream corepb.Orchestrator_ExecuteWorkflowServer) error { // Changed proto.OperationResponse to corepb.OperationResponse
	completedOps := make(map[string]corepb.OperationResponse) // 存储已完成操作的结果 // Changed proto.OperationResponse to corepb.OperationResponse
	opStatus := make(map[string]string)                   // 操作状态：pending, running, completed, failed
	var mu sync.Mutex                                     // 保护 completedOps 和 opStatus

	// 初始化所有操作状态为 pending
	for _, op := range workflowDef.Operations {
		opStatus[op.ID] = "pending"
	}

	// 使用一个 goroutine 池来执行并发操作
	var wg sync.WaitGroup
	operationQueue := make(chan workflow.OperationDefinition, len(workflowDef.Operations))

	// 启动 worker goroutines
	numWorkers := 5 // 可以配置工作协程数量
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for op := range operationQueue {
				ctx, cancel := context.WithCancel(stream.Context()) // 为每个操作创建新的上下文
				if op.Timeout != "" {
					duration, err := time.ParseDuration(op.Timeout)
					if err == nil {
						ctx, cancel = context.WithTimeout(ctx, duration)
					}
				}
				
				opResp := o.executeSingleOperation(ctx, op, workflowDef.Variables, completedOps)
				cancel()

				mu.Lock()
				completedOps[op.ID] = *opResp // 存储操作结果，供后续依赖使用
				if opResp.Success {
					opStatus[op.ID] = "completed"
				} else {
					opStatus[op.ID] = "failed"
				}
				mu.Unlock()

				// 流式发送结果
				stream.Send(&corepb.WorkflowResponse{ // Changed proto.WorkflowResponse to corepb.WorkflowResponse
					OperationId: op.ID,
					Response:    opResp,
					Status:      opStatus[op.ID],
				})
			}
		}()
	}

	// 调度操作
	for {
		mu.Lock()
		allCompleted := true
		progressMadeInIteration := false

		for _, op := range workflowDef.Operations {
			if opStatus[op.ID] == "pending" {
				allCompleted = false
				if o.dependenciesSatisfied(op, completedOps) {
					opStatus[op.ID] = "running"
					operationQueue <- op
					progressMadeInIteration = true
				}
			} else if opStatus[op.ID] == "running" {
				allCompleted = false
			}
		}
		mu.Unlock()

		if allCompleted {
			break
		}

		if !progressMadeInIteration && !allCompleted {
			log.Printf("Workflow deadlock or stuck operations detected. Operations status: %v", opStatus)
			return fmt.Errorf("workflow deadlock or stuck operations detected")
		}
		time.Sleep(100 * time.Millisecond)
	}

	close(operationQueue)
	wg.Wait()

	log.Printf("Workflow %s execution finished.", workflowDef.Name)
	return nil
}

// executeSingleOperation 执行单个操作，并处理重试
func (o *Orchestrator) executeSingleOperation(ctx context.Context, op workflow.OperationDefinition, globalVariables map[string]string, completedOps map[string]corepb.OperationResponse) *corepb.OperationResponse { // Changed proto.OperationResponse to corepb.OperationResponse
	var opResp *corepb.OperationResponse // Changed proto.OperationResponse to corepb.OperationResponse
	var err error

	// 替换操作参数中的变量和前序操作的 artifacts
	resolvedParams := o.resolveParameters(op.Parameters, globalVariables, completedOps)
	opReq := &corepb.OperationRequest{ // Changed proto.OperationRequest to corepb.OperationRequest
		OperationId: op.ID,
		Tool:        op.Tool,
		Command:     op.Command,
		Params:      resolvedParams,
	}

	for attempt := 0; attempt <= op.Retry; attempt++ {
		opResp, err = o.ExecuteOperation(ctx, opReq)
		if err == nil && opResp.Success {
			return opResp
		}
		if attempt < op.Retry {
			log.Printf("Operation %s failed (attempt %d/%d), retrying in 1 second: %v", op.ID, attempt+1, op.Retry, err)
			time.Sleep(1 * time.Second) // 重试间隔
		}
	}
	if err != nil {
		return &corepb.OperationResponse{Success: false, Error: fmt.Sprintf("Operation %s failed after %d retries: %v", op.ID, op.Retry+1, err)} // Changed proto.OperationResponse to corepb.OperationResponse
	}
	return opResp // 返回最后一次失败的响应
}


// ExecuteOperation 根据 req.Tool 路由到正确的适配器客户端
func (o *Orchestrator) ExecuteOperation(ctx context.Context, req *corepb.OperationRequest) (*corepb.OperationResponse, error) { // Changed proto.OperationRequest, proto.OperationResponse to corepb.OperationRequest, corepb.OperationResponse
	o.mu.Lock()
	defer o.mu.Unlock()

	switch req.Tool {
	case "hexstrike":
		if o.hexstrikeClient == nil {
			return nil, fmt.Errorf("hexstrike adapter not connected")
		}
		return o.hexstrikeClient.Execute(ctx, req)
	case "thesilencer":
		if o.thesilencerClient == nil {
			return nil, fmt.Errorf("thesilencer adapter not connected")
		}
		return o.thesilencerClient.Execute(ctx, req)
	case "netrunners":
		if o.netrunnersClient == nil {
			return nil, fmt.Errorf("netrunners adapter not connected")
		}
		return o.netrunnersClient.Execute(ctx, req)
	case "villager_ai":
		if o.villagerAIClient == nil {
			return nil, fmt.Errorf("villager_ai adapter not connected")
		}
		return o.villagerAIClient.Execute(ctx, req)
	case "redteam_tools":
		if o.redteamToolsClient == nil {
			return nil, fmt.Errorf("redteam_tools adapter not connected")
		}
		return o.redteamToolsClient.Execute(ctx, req)
	default:
		return nil, fmt.Errorf("unknown tool: %s", req.Tool)
	}
}

// GetStatus 返回所有适配器的健康状态
func (o *Orchestrator) GetStatus(ctx context.Context, req *corepb.StatusRequest) (*corepb.StatusResponse, error) { // Changed proto.StatusRequest, proto.StatusResponse to corepb.StatusRequest, corepb.StatusResponse
	o.mu.Lock()
	defer o.mu.Unlock()

	statusMap := make(map[string]string)

	// 检查 Hexstrike
	if o.hexstrikeClient != nil {
		health, err := o.hexstrikeClient.HealthCheck(ctx, &corepb.Empty{}) // Changed proto.Empty to corepb.Empty
		if err == nil {
			statusMap["hexstrike"] = health.Status
		} else {
			statusMap["hexstrike"] = fmt.Sprintf("error: %v", err)
		}
	} else {
		statusMap["hexstrike"] = "disconnected"
	}

	// 检查 TheSilencer
	if o.thesilencerClient != nil {
		health, err := o.thesilencerClient.HealthCheck(ctx, &corepb.Empty{}) // Changed proto.Empty to corepb.Empty
		if err == nil {
			statusMap["thesilencer"] = health.Status
		} else {
			statusMap["thesilencer"] = fmt.Sprintf("error: %v", err)
		}
	} else {
		statusMap["thesilencer"] = "disconnected"
	}
	
	// 检查 NetRunners
	if o.netrunnersClient != nil {
		health, err := o.netrunnersClient.HealthCheck(ctx, &corepb.Empty{}) // Changed proto.Empty to corepb.Empty
		if err == nil {
			statusMap["netrunners"] = health.Status
		} else {
			statusMap["netrunners"] = fmt.Sprintf("error: %v", err)
		}
	} else {
		statusMap["netrunners"] = "disconnected"
	}

	// 检查 Villager-AI
	if o.villagerAIClient != nil {
		health, err := o.villagerAIClient.HealthCheck(ctx, &corepb.Empty{}) // Changed proto.Empty to corepb.Empty
		if err == nil {
			statusMap["villager_ai"] = health.Status
		} else {
			statusMap["villager_ai"] = fmt.Sprintf("error: %v", err)
		}
	} else {
		statusMap["villager_ai"] = "disconnected"
	}

	// 检查 Redteam Tools
	if o.redteamToolsClient != nil {
		health, err := o.redteamToolsClient.HealthCheck(ctx, &corepb.Empty{}) // Changed proto.Empty to corepb.Empty
		if err == nil {
			statusMap["redteam_tools"] = health.Status
		} else {
			statusMap["redteam_tools"] = fmt.Sprintf("error: %v", err)
		}
	} else {
		statusMap["redteam_tools"] = "disconnected"
	}

	return &corepb.StatusResponse{ToolStatus: statusMap}, nil // Changed proto.StatusResponse to corepb.StatusResponse
}

// dependenciesSatisfied 检查一个操作的所有依赖是否已完成
func (o *Orchestrator) dependenciesSatisfied(op workflow.OperationDefinition, completedOps map[string]corepb.OperationResponse) bool { // Changed proto.OperationResponse to corepb.OperationResponse
	for _, depID := range op.DependsOn {
		if _, ok := completedOps[depID]; !ok {
			return false
		}
	}
	return true
}

// resolveParameters 替换操作参数中的变量和前序操作的 artifacts
func (o *Orchestrator) resolveParameters(params map[string]string, globalVariables map[string]string, completedOps map[string]corepb.OperationResponse) map[string]string { // Changed proto.OperationResponse to corepb.OperationResponse
	resolved := make(map[string]string)
	for k, v := range params {
		// 先替换全局变量
		resolvedValue := o.replaceVariables(v, globalVariables)

		// 然后替换前序操作的 artifacts (例如：${op_id.artifact_key})
		resolved[k] = o.replaceArtifacts(resolvedValue, completedOps)
	}
	return resolved
}

// replaceVariables 替换字符串中的 ${var_name} 格式的变量
func (o *Orchestrator) replaceVariables(input string, variables map[string]string) string {
	output := input
	for k, v := range variables {
		output = strings.ReplaceAll(output, fmt.Sprintf("${%s}", k), v)
	}
	return output
}

// replaceArtifacts 替换字符串中的 ${op_id.artifact_key} 格式的 artifacts
func (o *Orchestrator) replaceArtifacts(input string, completedOps map[string]corepb.OperationResponse) string { // Changed proto.OperationResponse to corepb.OperationResponse
	output := input
	// 简单的正则表达式或字符串解析来查找 ${op_id.artifact_key}
	// 这里做简化处理，实际需要更健壮的解析器
	for opID, opResp := range completedOps {
		for artifactKey, artifactValue := range opResp.Artifacts {
			placeholder := fmt.Sprintf("${%s.%s}", opID, artifactKey)
			output = strings.ReplaceAll(output, placeholder, artifactValue)
		}
	}
	return output
}

// main 函数现在在单独的文件中定义，例如 cmd/orchestrator/main.go
// 这里只是 Orchestrator 核心逻辑的定义
