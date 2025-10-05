package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
)

// JSONRPCRequest JSON-RPC 请求
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse JSON-RPC 响应
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
}

// Tool 工具信息
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolsListResult 工具列表结果
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

func main() {
	// 启动 MCP Server
	cmd := exec.Command("../yst-go-mcp")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// 启动错误日志监控
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Println("[Server Log]", scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	fmt.Println("=== 测试工具列表 ===\n")

	// 1. 初始化
	initReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]string{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := json.NewEncoder(stdin).Encode(initReq); err != nil {
		log.Fatalf("发送初始化请求失败: %v", err)
	}

	// 读取初始化响应
	var initResp JSONRPCResponse
	if err := json.NewDecoder(stdout).Decode(&initResp); err != nil {
		log.Fatalf("读取初始化响应失败: %v", err)
	}
	fmt.Println("✓ 初始化成功\n")

	// 2. 列出工具
	toolsReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	if err := json.NewEncoder(stdin).Encode(toolsReq); err != nil {
		log.Fatalf("发送工具列表请求失败: %v", err)
	}

	// 读取工具列表响应
	var toolsResp JSONRPCResponse
	if err := json.NewDecoder(stdout).Decode(&toolsResp); err != nil {
		log.Fatalf("读取工具列表响应失败: %v", err)
	}

	// 解析工具列表
	var result ToolsListResult
	if err := json.Unmarshal(toolsResp.Result, &result); err != nil {
		log.Fatalf("解析工具列表失败: %v", err)
	}

	fmt.Printf("共注册 %d 个工具：\n\n", len(result.Tools))

	for i, tool := range result.Tools {
		fmt.Printf("%d. %s\n", i+1, tool.Name)
		fmt.Printf("   描述: %s\n", tool.Description)

		// 打印参数
		if props, ok := tool.InputSchema["properties"].(map[string]interface{}); ok && len(props) > 0 {
			fmt.Println("   参数:")
			for paramName, paramInfo := range props {
				if param, ok := paramInfo.(map[string]interface{}); ok {
					desc := param["description"]
					fmt.Printf("     - %s: %v\n", paramName, desc)
				}
			}
		} else {
			fmt.Println("   参数: 无")
		}
		fmt.Println()
	}

	// 关闭
	stdin.(io.WriteCloser).Close()
	cmd.Wait()

	fmt.Println("✓ 测试完成！")
}
