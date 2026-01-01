package diagnostic

import (
	"context"

	"golang.org/x/sys/windows/registry"
	"network-rescue-toolkit/pkg/types"
)

// ProxyChecker 代理设置检查器
type ProxyChecker struct{}

// NewProxyChecker 创建代理设置检查器
func NewProxyChecker() *ProxyChecker {
	return &ProxyChecker{}
}

// ID 返回检查器 ID
func (c *ProxyChecker) ID() string {
	return "proxy"
}

// Name 返回检查器名称
func (c *ProxyChecker) Name() string {
	return "代理设置检测"
}

// Check 执行检查
func (c *ProxyChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	proxyConfig := c.readProxySettings()
	result.AddDetail("proxyConfig", proxyConfig)

	if proxyConfig.Enabled {
		result.SetWarning("检测到已启用代理服务器: "+proxyConfig.Server, true)
	} else {
		result.SetOK("未启用代理服务器")
	}

	return *result
}

// readProxySettings 从注册表读取代理设置
func (c *ProxyChecker) readProxySettings() types.ProxyConfig {
	config := types.ProxyConfig{}

	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return config
	}
	defer key.Close()

	// 读取代理启用状态
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err == nil {
		config.Enabled = proxyEnable == 1
	}

	// 读取代理服务器地址
	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err == nil {
		config.Server = proxyServer
	}

	// 读取代理绕过列表
	proxyOverride, _, err := key.GetStringValue("ProxyOverride")
	if err == nil {
		config.BypassList = proxyOverride
	}

	// 读取自动配置 URL
	autoConfigURL, _, err := key.GetStringValue("AutoConfigURL")
	if err == nil {
		config.AutoConfigURL = autoConfigURL
	}

	return config
}
