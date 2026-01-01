package diagnostic

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"network-rescue-toolkit/pkg/types"
)

// HostsChecker HOSTS 文件检查器
type HostsChecker struct{}

// NewHostsChecker 创建 HOSTS 文件检查器
func NewHostsChecker() *HostsChecker {
	return &HostsChecker{}
}

// ID 返回检查器 ID
func (c *HostsChecker) ID() string {
	return "hosts"
}

// Name 返回检查器名称
func (c *HostsChecker) Name() string {
	return "HOSTS 文件检测"
}

// Check 执行检查
func (c *HostsChecker) Check(ctx context.Context) types.DiagnosticResult {
	result := types.NewDiagnosticResult(c.ID(), c.Name())

	hostsPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
	
	file, err := os.Open(hostsPath)
	if err != nil {
		result.SetError("无法读取 HOSTS 文件: "+err.Error(), true)
		return *result
	}
	defer file.Close()

	entries := make([]types.HostsEntry, 0)
	suspiciousCount := 0
	lineNum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry := c.parseLine(line, lineNum)
		if entry != nil {
			entry.Suspicious = entry.IsSuspicious()
			if entry.Suspicious {
				suspiciousCount++
			}
			entries = append(entries, *entry)
		}
	}

	result.AddDetail("entries", entries)
	result.AddDetail("totalEntries", len(entries))
	result.AddDetail("suspiciousCount", suspiciousCount)

	if suspiciousCount > 0 {
		result.SetWarning("检测到可疑的 HOSTS 条目", true)
	} else if len(entries) > 10 {
		result.SetWarning("HOSTS 文件包含较多自定义条目", false)
	} else {
		result.SetOK("HOSTS 文件正常")
	}

	return *result
}

// parseLine 解析 HOSTS 文件行
func (c *HostsChecker) parseLine(line string, lineNum int) *types.HostsEntry {
	// 移除行内注释
	if idx := strings.Index(line, "#"); idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}

	if line == "" {
		return nil
	}

	// 分割 IP 和主机名
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}

	return &types.HostsEntry{
		IP:       fields[0],
		Hostname: fields[1],
		LineNum:  lineNum,
	}
}
