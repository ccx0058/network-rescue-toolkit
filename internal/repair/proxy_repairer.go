package repair

import (
	"context"
	"time"

	"golang.org/x/sys/windows/registry"
	"network-rescue-toolkit/pkg/types"
)

// ProxyRepairer 代理修复器
type ProxyRepairer struct{}

// NewProxyRepairer 创建代理修复器
func NewProxyRepairer() *ProxyRepairer {
	return &ProxyRepairer{}
}

// ID 返回修复器 ID
func (r *ProxyRepairer) ID() string {
	return "proxy"
}

// Name 返回修复器名称
func (r *ProxyRepairer) Name() string {
	return "清除代理设置"
}

// RequiresAdmin 是否需要管理员权限
func (r *ProxyRepairer) RequiresAdmin() bool {
	return false
}

// Repair 执行修复
func (r *ProxyRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.SET_VALUE)
	if err != nil {
		result.SetFailure("无法打开注册表: " + err.Error())
		return *result
	}
	defer key.Close()

	// 禁用代理
	err = key.SetDWordValue("ProxyEnable", 0)
	if err != nil {
		result.SetFailure("无法禁用代理: " + err.Error())
		return *result
	}

	// 验证修改
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil || proxyEnable != 0 {
		result.SetFailure("代理设置验证失败")
		return *result
	}

	result.SetSuccess("代理设置已清除")
	return *result
}
