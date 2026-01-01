package diagnostic

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"network-rescue-toolkit/pkg/types"
)

// ConnectivityChecker 网络连通性检查器
type ConnectivityChecker struct{}

// NewConnectivityChecker 创建网络连通性检查器
func NewConnectivityChecker() *ConnectivityChecker {
	return &ConnectivityChecker{}
}

// ID 返回检查器 ID
func (c *ConnectivityChecker) ID() string {
	return "connectivity"
}

// Name 返回检查器名称
func (c *ConnectivityChecker) Name() string {
	return "电脑能否上网"
}

// Check 执行检查
func (c *ConnectivityChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	// 测试多个目标，任一成功即可
	targets := []struct {
		name string
		test func() (bool, int64)
	}{
		{"HTTP-百度", func() (bool, int64) { return c.testHTTP("https://www.baidu.com") }},
		{"HTTP-腾讯", func() (bool, int64) { return c.testHTTP("https://www.qq.com") }},
		{"HTTP-阿里", func() (bool, int64) { return c.testHTTP("https://www.taobao.com") }},
		{"TCP-114DNS", func() (bool, int64) { return c.testTCP("114.114.114.114:53") }},
		{"TCP-阿里DNS", func() (bool, int64) { return c.testTCP("223.5.5.5:53") }},
	}

	var successCount int
	var totalLatency int64
	var lastError string

	for _, t := range targets {
		ok, latency := t.test()
		if ok {
			successCount++
			totalLatency += latency
		} else {
			lastError = t.name + " 连接失败"
		}
	}

	result.AddDetail("successCount", successCount)
	result.AddDetail("totalTargets", len(targets))

	if successCount >= 2 {
		avgLatency := totalLatency / int64(successCount)
		result.SetOK(fmt.Sprintf("网络连通正常，平均延迟 %dms", avgLatency))
	} else if successCount >= 1 {
		result.SetWarning("网络连通不稳定", true)
	} else {
		result.SetError("无法连接到互联网: "+lastError, true)
	}

	return *result
}

// testHTTP 测试 HTTP 连通性
func (c *ConnectivityChecker) testHTTP(url string) (bool, int64) {
	client := &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不跟随重定向
		},
	}

	start := time.Now()
	resp, err := client.Head(url) // 使用 HEAD 请求更快
	latency := time.Since(start).Milliseconds()

	if err != nil {
		// 尝试 GET 请求
		resp, err = client.Get(url)
		if err != nil {
			return false, 0
		}
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 500, latency
}

// testTCP 测试 TCP 连通性
func (c *ConnectivityChecker) testTCP(addr string) (bool, int64) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return false, 0
	}
	conn.Close()
	return true, latency
}
