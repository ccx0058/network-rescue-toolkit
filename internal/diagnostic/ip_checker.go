package diagnostic

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"network-rescue-toolkit/pkg/executor"
	"network-rescue-toolkit/pkg/types"
)

// IPChecker IP 配置检查器
type IPChecker struct {
	executor *executor.CommandExecutor
}

// NewIPChecker 创建 IP 配置检查器
func NewIPChecker() *IPChecker {
	return &IPChecker{
		executor: executor.NewCommandExecutor(),
	}
}

// ID 返回检查器 ID
func (c *IPChecker) ID() string {
	return "ip"
}

// Name 返回检查器名称
func (c *IPChecker) Name() string {
	return "IP 配置检测"
}

// Check 执行检查
func (c *IPChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	// 执行 ipconfig /all
	cmdResult := c.executor.ExecuteIPConfig(ctx, "/all")
	if !cmdResult.IsSuccess() {
		result.SetError("无法获取 IP 配置信息: "+cmdResult.Stderr, false)
		return *result
	}

	// 解析输出
	configs := c.parseIPConfig(cmdResult.Stdout)
	result.AddDetail("configs", configs)

	// 检查是否有有效的 IP 配置
	hasValidIP := false
	hasDHCP := false
	for _, cfg := range configs {
		if len(cfg.IPAddresses) > 0 {
			hasValidIP = true
		}
		if cfg.DHCPEnabled {
			hasDHCP = true
		}
	}

	if !hasValidIP {
		result.SetError("未检测到有效的 IP 地址配置", true)
	} else if hasDHCP {
		result.SetOK("IP 配置正常 (DHCP)")
	} else {
		result.SetOK("IP 配置正常 (静态)")
	}

	return *result
}

// parseIPConfig 解析 ipconfig 输出
func (c *IPChecker) parseIPConfig(output string) []types.AdapterConfig {
	configs := make([]types.AdapterConfig, 0)
	
	// 按适配器分割
	sections := strings.Split(output, "适配器")
	if len(sections) <= 1 {
		sections = strings.Split(output, "adapter")
	}

	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	
	for _, section := range sections[1:] {
		lines := strings.Split(section, "\n")
		if len(lines) == 0 {
			continue
		}

		config := types.AdapterConfig{
			Name: strings.TrimSpace(strings.Split(lines[0], ":")[0]),
		}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			lowerLine := strings.ToLower(line)

			if strings.Contains(lowerLine, "dhcp") && strings.Contains(lowerLine, "是") {
				config.DHCPEnabled = true
			}
			if strings.Contains(lowerLine, "dhcp") && strings.Contains(lowerLine, "yes") {
				config.DHCPEnabled = true
			}

			if strings.Contains(lowerLine, "ipv4") || strings.Contains(lowerLine, "ip address") {
				if matches := ipRegex.FindString(line); matches != "" {
					config.IPAddresses = append(config.IPAddresses, matches)
				}
			}

			if strings.Contains(lowerLine, "子网掩码") || strings.Contains(lowerLine, "subnet mask") {
				if matches := ipRegex.FindString(line); matches != "" {
					config.SubnetMasks = append(config.SubnetMasks, matches)
				}
			}

			if strings.Contains(lowerLine, "默认网关") || strings.Contains(lowerLine, "default gateway") {
				if matches := ipRegex.FindString(line); matches != "" {
					config.Gateways = append(config.Gateways, matches)
				}
			}

			if strings.Contains(lowerLine, "dns") {
				if matches := ipRegex.FindString(line); matches != "" {
					config.DNSServers = append(config.DNSServers, matches)
				}
			}
		}

		if config.Name != "" {
			configs = append(configs, config)
		}
	}

	return configs
}

// formatAdapterCount 格式化适配器数量
func formatAdapterCount(count int) string {
	return fmt.Sprintf("%d", count)
}
