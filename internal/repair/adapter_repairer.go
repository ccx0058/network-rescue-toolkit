package repair

import (
	"context"
	"time"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// AdapterRepairer 网络适配器修复器
type AdapterRepairer struct {
	executor *executor.CommandExecutor
}

// NewAdapterRepairer 创建网络适配器修复器
func NewAdapterRepairer() *AdapterRepairer {
	return &AdapterRepairer{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回修复器 ID
func (r *AdapterRepairer) ID() string {
	return "adapter"
}

// Name 返回修复器名称
func (r *AdapterRepairer) Name() string {
	return "重置网络适配器"
}

// RequiresAdmin 是否需要管理员权限
func (r *AdapterRepairer) RequiresAdmin() bool {
	return true
}

// Repair 执行修复
func (r *AdapterRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	// 获取所有网络适配器并重置
	// 使用 netsh 禁用再启用
	disableResult := r.executor.ExecuteNetsh(ctx, "interface", "set", "interface", "以太网", "disable")
	if !disableResult.IsSuccess() {
		// 尝试英文名称
		disableResult = r.executor.ExecuteNetsh(ctx, "interface", "set", "interface", "Ethernet", "disable")
	}

	// 等待 2 秒
	time.Sleep(2 * time.Second)

	// 重新启用
	enableResult := r.executor.ExecuteNetsh(ctx, "interface", "set", "interface", "以太网", "enable")
	if !enableResult.IsSuccess() {
		enableResult = r.executor.ExecuteNetsh(ctx, "interface", "set", "interface", "Ethernet", "enable")
	}

	if enableResult.IsSuccess() {
		result.SetSuccess("网络适配器已重置")
	} else {
		result.SetFailure("网络适配器重置失败，请手动在设备管理器中操作")
	}

	return *result
}
