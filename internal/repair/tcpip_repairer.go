package repair

import (
	"context"
	"time"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// TCPIPRepairer TCP/IP 修复器
type TCPIPRepairer struct {
	executor *executor.CommandExecutor
}

// NewTCPIPRepairer 创建 TCP/IP 修复器
func NewTCPIPRepairer() *TCPIPRepairer {
	return &TCPIPRepairer{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回修复器 ID
func (r *TCPIPRepairer) ID() string {
	return "tcpip"
}

// Name 返回修复器名称
func (r *TCPIPRepairer) Name() string {
	return "重置 TCP/IP"
}

// RequiresAdmin 是否需要管理员权限
func (r *TCPIPRepairer) RequiresAdmin() bool {
	return true
}

// Repair 执行修复
func (r *TCPIPRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	// 执行 netsh int ip reset
	cmdResult := r.executor.ExecuteNetsh(ctx, "int", "ip", "reset")

	if cmdResult.IsSuccess() {
		result.SetSuccess("TCP/IP 协议栈重置成功，需要重启计算机以完成修复")
		result.SetRequireReboot()
	} else {
		result.SetFailure("TCP/IP 重置失败: " + cmdResult.Stderr)
	}

	return *result
}
