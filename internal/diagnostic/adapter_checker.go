package diagnostic

import (
	"context"
	"net"

	"network-rescue-toolkit/pkg/types"
)

// AdapterChecker 网络适配器检查器
type AdapterChecker struct{}

// NewAdapterChecker 创建网络适配器检查器
func NewAdapterChecker() *AdapterChecker {
	return &AdapterChecker{}
}

// ID 返回检查器 ID
func (c *AdapterChecker) ID() string {
	return "adapter"
}

// Name 返回检查器名称
func (c *AdapterChecker) Name() string {
	return "网络适配器检测"
}

// Check 执行检查
func (c *AdapterChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	interfaces, err := net.Interfaces()
	if err != nil {
		result.SetError("无法获取网络适配器列表: "+err.Error(), false)
		return *result
	}

	if len(interfaces) == 0 {
		result.SetError("未检测到任何网络适配器", false)
		return *result
	}

	adapters := make([]types.AdapterInfo, 0)
	activeCount := 0

	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		adapter := types.AdapterInfo{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr.String(),
		}

		// 检查接口状态
		if iface.Flags&net.FlagUp != 0 {
			adapter.Status = "Up"
			activeCount++
		} else {
			adapter.Status = "Down"
		}

		// 获取 IP 地址
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					adapter.IPAddresses = append(adapter.IPAddresses, ipnet.IP.String())
					mask := net.IP(ipnet.Mask).String()
					adapter.SubnetMasks = append(adapter.SubnetMasks, mask)
				}
			}
		}

		adapters = append(adapters, adapter)
	}

	result.AddDetail("adapters", adapters)
	result.AddDetail("totalCount", len(adapters))
	result.AddDetail("activeCount", activeCount)

	if activeCount == 0 {
		result.SetWarning("未检测到活动的网络适配器", true)
	} else {
		result.SetOK("检测到 " + string(rune('0'+activeCount)) + " 个活动网络适配器")
	}

	return *result
}
