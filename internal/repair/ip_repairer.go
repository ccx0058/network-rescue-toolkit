package repair

import (
	"context"
	"time"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// IPRepairer IP 修复器
type IPRepairer struct {
	executor *executor.CommandExecutor
}

// NewIPRepairer 创建 IP 修复器
func NewIPRepairer() *IPRepairer {
	return &IPRepairer{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回修复器 ID
func (r *IPRepairer) ID() string {
	return "ip"
}

// Name 返回修复器名称
func (r *IPRepairer) Name() string {
	return "释放/续租 IP"
}

// RequiresAdmin 是否需要管理员权限
func (r *IPRepairer) RequiresAdmin() bool {
	return true
}

// Repair 执行修复
func (r *IPRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	// 执行 ipconfig /release
	releaseResult := r.executor.ExecuteIPConfig(ctx, "/release")
	if !releaseResult.IsSuccess() {
		result.SetFailure("IP 释放失败: " + releaseResult.Stderr)
		return *result
	}

	// 等待一秒
	time.Sleep(time.Second)

	// 执行 ipconfig /renew
	renewResult := r.executor.ExecuteIPConfig(ctx, "/renew")
	if !renewResult.IsSuccess() {
		result.SetFailure("IP 续租失败: " + renewResult.Stderr)
		return *result
	}

	result.SetSuccess("IP 地址已重新获取")
	return *result
}
