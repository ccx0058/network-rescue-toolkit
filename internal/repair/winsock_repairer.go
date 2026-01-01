package repair

import (
	"context"
	"time"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// WinsockRepairer Winsock 修复器
type WinsockRepairer struct {
	executor *executor.CommandExecutor
}

// NewWinsockRepairer 创建 Winsock 修复器
func NewWinsockRepairer() *WinsockRepairer {
	return &WinsockRepairer{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回修复器 ID
func (r *WinsockRepairer) ID() string {
	return "winsock"
}

// Name 返回修复器名称
func (r *WinsockRepairer) Name() string {
	return "重置 Winsock"
}

// RequiresAdmin 是否需要管理员权限
func (r *WinsockRepairer) RequiresAdmin() bool {
	return true
}

// Repair 执行修复
func (r *WinsockRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	// 执行 netsh winsock reset
	cmdResult := r.executor.ExecuteNetsh(ctx, "winsock", "reset")

	if cmdResult.IsSuccess() {
		result.SetSuccess("Winsock 重置成功，需要重启计算机以完成修复")
		result.SetRequireReboot()
	} else {
		result.SetFailure("Winsock 重置失败: " + cmdResult.Stderr)
	}

	return *result
}
