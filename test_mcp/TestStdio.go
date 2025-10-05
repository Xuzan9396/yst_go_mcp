package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("=== 简单 STDIO 测试 ===\n")

	// 直接通过 JSON-RPC 协议测试
	cmd := exec.Command("./yst-go-mcp")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// 发送初始化请求
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`
	fmt.Println("发送初始化请求...")
	fmt.Fprintf(stdin, "%s\n", initReq)

	// 读取响应
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				fmt.Printf("收到: %s\n", string(buf[:n]))
			}
		}
	}()

	<-ctx.Done()

	// 列出工具
	toolsReq := `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`
	fmt.Println("\n发送列出工具请求...")
	fmt.Fprintf(stdin, "%s\n", toolsReq)

	time.Sleep(2 * time.Second)

	// 调用清除 Cookie
	clearReq := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"clear_saved_cookies","arguments":{}}}`
	fmt.Println("\n发送调用工具请求...")
	fmt.Fprintf(stdin, "%s\n", clearReq)

	time.Sleep(5 * time.Second)

	stdin.Close()
	cmd.Wait()
}
