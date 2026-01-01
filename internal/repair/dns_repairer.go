package repair

import (
	"context"
	"time"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// DNSRepairer DNS 修复器
type DNSRepairer struct {
	executor *executor.CommandExecutor
}

// NewDNSRepairer 创建 DNS 修复器
func NewDNSRepairer() *DNSRepairer {
	return &DNSRepairer{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回修复器 ID
func (r *DNSRepairer) ID() string {
	return "dns"
}

// Name 返回修复器名称
func (r *DNSRepairer) Name() string {
	return "刷新 DNS 缓存"
}

// RequiresAdmin 是否需要管理员权限
func (r *DNSRepairer) RequiresAdmin() bool {
	return false
}

// Repair 执行修复
func (r *DNSRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	// 执行 ipconfig /flushdns
	cmdResult := r.executor.ExecuteIPConfig(ctx, "/flushdns")

	if cmdResult.IsSuccess() {
		result.SetSuccess("DNS 缓存已刷新")
	} else {
		result.SetFailure("DNS 缓存刷新失败: " + cmdResult.Stderr)
	}

	return *result
}
