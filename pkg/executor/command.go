package executor

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// CommandResult 命令执行结果
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// IsSuccess 检查命令是否执行成功
func (r *CommandResult) IsSuccess() bool {
	return r.ExitCode == 0 && r.Error == nil
}

// CommandExecutor 命令执行器
type CommandExecutor struct {
	timeout time.Duration
}

// NewCommandExecutor 创建命令执行器
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		timeout: 30 * time.Second,
	}
}

// SetTimeout 设置超时时间
func (e *CommandExecutor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// Execute 执行命令
func (e *CommandExecutor) Execute(ctx context.Context, name string, args ...string) CommandResult {
	return e.executeInternal(ctx, name, args, false)
}


// ExecuteAsAdmin 以管理员权限执行命令
func (e *CommandExecutor) ExecuteAsAdmin(ctx context.Context, name string, args ...string) CommandResult {
	return e.executeInternal(ctx, name, args, true)
}

// executeInternal 内部执行命令
func (e *CommandExecutor) executeInternal(ctx context.Context, name string, args []string, asAdmin bool) CommandResult {
	result := CommandResult{}

	// 创建带超时的上下文
	execCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	var cmd *exec.Cmd
	if asAdmin {
		// Windows 下使用 runas 提升权限
		allArgs := append([]string{"/C", name}, args...)
		cmd = exec.CommandContext(execCtx, "cmd", allArgs...)
	} else {
		cmd = exec.CommandContext(execCtx, name, args...)
	}

	// 隐藏命令窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// 转换 GBK 编码到 UTF-8（Windows 命令行默认 GBK）
	result.Stdout = decodeGBK(stdout.Bytes())
	result.Stderr = decodeGBK(stderr.Bytes())

	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// decodeGBK 将 GBK 编码转换为 UTF-8
func decodeGBK(data []byte) string {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return strings.TrimSpace(buf.String())
}

// ExecuteNetsh 执行 netsh 命令
func (e *CommandExecutor) ExecuteNetsh(ctx context.Context, args ...string) CommandResult {
	return e.Execute(ctx, "netsh", args...)
}

// ExecuteIPConfig 执行 ipconfig 命令
func (e *CommandExecutor) ExecuteIPConfig(ctx context.Context, args ...string) CommandResult {
	return e.Execute(ctx, "ipconfig", args...)
}

// ExecutePing 执行 ping 命令
func (e *CommandExecutor) ExecutePing(ctx context.Context, target string, count int) CommandResult {
	return e.Execute(ctx, "ping", "-n", string(rune('0'+count)), target)
}
