package diagnostic

import (
	"context"
	"fmt"
	"net"
	"time"

	"network-rescue-toolkit/pkg/types"
)

// DNSChecker DNS 检查器
type DNSChecker struct{}

// NewDNSChecker 创建 DNS 检查器
func NewDNSChecker() *DNSChecker {
	return &DNSChecker{}
}

// ID 返回检查器 ID
func (c *DNSChecker) ID() string {
	return "dns"
}

// Name 返回检查器名称
func (c *DNSChecker) Name() string {
	return "DNS 服务检测"
}

// Check 执行检查
func (c *DNSChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	// 测试 DNS 解析
	testDomains := []string{"www.baidu.com", "www.qq.com"}
	var successCount int
	var totalLatency int64

	for _, domain := range testDomains {
		start := time.Now()
		_, err := net.LookupHost(domain)
		latency := time.Since(start).Milliseconds()

		if err == nil {
			successCount++
			totalLatency += latency
		}
	}

	result.AddDetail("testedDomains", testDomains)
	result.AddDetail("successCount", successCount)

	if successCount == 0 {
		result.SetError("DNS 解析失败，无法解析任何域名", true)
	} else if successCount < len(testDomains) {
		avgLatency := totalLatency / int64(successCount)
		result.AddDetail("avgLatencyMs", avgLatency)
		result.SetWarning(fmt.Sprintf("DNS 部分正常，平均延迟 %dms", avgLatency), true)
	} else {
		avgLatency := totalLatency / int64(successCount)
		result.AddDetail("avgLatencyMs", avgLatency)
		result.SetOK(fmt.Sprintf("DNS 服务正常，平均延迟 %dms", avgLatency))
	}

	return *result
}
