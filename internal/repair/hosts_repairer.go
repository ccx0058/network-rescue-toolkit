package repair

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"network-rescue-toolkit/pkg/types"
)

// HostsRepairer HOSTS 修复器
type HostsRepairer struct{}

// NewHostsRepairer 创建 HOSTS 修复器
func NewHostsRepairer() *HostsRepairer {
	return &HostsRepairer{}
}

// ID 返回修复器 ID
func (r *HostsRepairer) ID() string {
	return "hosts"
}

// Name 返回修复器名称
func (r *HostsRepairer) Name() string {
	return "修复 HOSTS 文件"
}

// RequiresAdmin 是否需要管理员权限
func (r *HostsRepairer) RequiresAdmin() bool {
	return true
}

// Repair 执行修复
func (r *HostsRepairer) Repair(ctx context.Context) types.RepairResult {
	result := types.NewRepairResult(r.ID(), r.Name())
	result.Timestamp = time.Now()

	hostsPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
	backupPath := hostsPath + ".backup"

	// 备份当前文件
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		result.SetFailure("无法读取 HOSTS 文件: " + err.Error())
		return *result
	}

	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		result.SetFailure("无法备份 HOSTS 文件: " + err.Error())
		return *result
	}
	result.SetBackupPath(backupPath)

	// 写入默认内容
	defaultHosts := `# Copyright (c) 1993-2009 Microsoft Corp.
#
# This is a sample HOSTS file used by Microsoft TCP/IP for Windows.
#
# This file contains the mappings of IP addresses to host names. Each
# entry should be kept on an individual line. The IP address should
# be placed in the first column followed by the corresponding host name.
# The IP address and the host name should be separated by at least one
# space.
#
# Additionally, comments (such as these) may be inserted on individual
# lines or following the machine name denoted by a '#' symbol.
#
# For example:
#
#      102.54.94.97     rhino.acme.com          # source server
#       38.25.63.10     x.acme.com              # x client host

# localhost name resolution is handled within DNS itself.
#	127.0.0.1       localhost
#	::1             localhost
`

	err = os.WriteFile(hostsPath, []byte(defaultHosts), 0644)
	if err != nil {
		result.SetFailure("无法写入 HOSTS 文件: " + err.Error())
		return *result
	}

	result.SetSuccess("HOSTS 文件已恢复为默认状态")
	return *result
}
