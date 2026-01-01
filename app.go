package main

import (
	"context"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"

	"network-rescue-toolkit/internal/diagnostic"
	"network-rescue-toolkit/internal/repair"
	"network-rescue-toolkit/pkg/backup"
	"network-rescue-toolkit/pkg/privilege"
	"network-rescue-toolkit/pkg/report"
	"network-rescue-toolkit/pkg/types"
)

// App 应用主结构
type App struct {
	ctx              context.Context
	diagnosticEngine *diagnostic.Engine
	repairEngine     *repair.Engine
	backupManager    *backup.Manager
	reportGenerator  *report.Generator
	privilegeHelper  *privilege.Helper
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{
		diagnosticEngine: diagnostic.NewEngine(),
		repairEngine:     repair.NewEngine(),
		backupManager:    backup.NewManager(),
		reportGenerator:  report.NewGenerator(),
		privilegeHelper:  privilege.NewHelper(),
	}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// IsAdmin 检查是否以管理员权限运行
func (a *App) IsAdmin() bool {
	return a.privilegeHelper.IsAdmin()
}

// RunDiagnostic 运行完整诊断
func (a *App) RunDiagnostic() []types.DiagnosticResult {
	return a.diagnosticEngine.RunAll(a.ctx)
}

// RunSingleDiagnostic 运行单个诊断项
func (a *App) RunSingleDiagnostic(id string) types.DiagnosticResult {
	return a.diagnosticEngine.RunSingle(a.ctx, id)
}

// RunRepair 执行修复操作
func (a *App) RunRepair(id string) types.RepairResult {
	return a.repairEngine.Repair(a.ctx, id)
}

// RunComprehensiveRepair 执行综合修复
func (a *App) RunComprehensiveRepair() []types.RepairResult {
	return a.repairEngine.RepairAll(a.ctx)
}

// CreateBackup 创建配置备份
func (a *App) CreateBackup() (string, error) {
	return a.backupManager.CreateBackup()
}

// RestoreBackup 还原配置备份
func (a *App) RestoreBackup(path string) error {
	return a.backupManager.RestoreBackup(path)
}

// ListBackups 列出所有备份
func (a *App) ListBackups() ([]backup.BackupInfo, error) {
	return a.backupManager.ListBackups()
}

// ExportReport 导出诊断报告
func (a *App) ExportReport(format string) (string, error) {
	results := a.diagnosticEngine.RunAll(a.ctx)
	return a.reportGenerator.Generate(results, format)
}

// RequestElevation 请求管理员权限
func (a *App) RequestElevation() error {
	return a.privilegeHelper.RequestElevation()
}

// ========== 网络工具 ==========

// RunPing 执行 Ping 测试
func (a *App) RunPing(target string) string {
	// 从 URL 中提取域名
	target = strings.TrimSpace(target)
	target = strings.TrimPrefix(target, "https://")
	target = strings.TrimPrefix(target, "http://")
	if idx := strings.Index(target, "/"); idx != -1 {
		target = target[:idx]
	}
	if idx := strings.Index(target, ":"); idx != -1 {
		target = target[:idx]
	}

	out, err := exec.Command("ping", "-n", "4", target).Output()
	if err != nil {
		return "Ping 失败: " + err.Error()
	}
	// 转换编码 (GBK -> UTF8)
	result, _ := simplifiedchinese.GBK.NewDecoder().String(string(out))
	return result
}

// SwitchDNS 切换 DNS 服务器
func (a *App) SwitchDNS(primary, secondary string) bool {
	// 获取活动网卡名称
	adapters := a.getActiveAdapters()
	if len(adapters) == 0 {
		return false
	}

	for _, adapter := range adapters {
		// 设置主 DNS
		exec.Command("netsh", "interface", "ip", "set", "dns", adapter, "static", primary).Run()
		// 设置备用 DNS
		exec.Command("netsh", "interface", "ip", "add", "dns", adapter, secondary, "index=2").Run()
	}
	return true
}

// FlushDNS 刷新 DNS 缓存
func (a *App) FlushDNS() error {
	return exec.Command("ipconfig", "/flushdns").Run()
}

// ResetNetworkStack 重置网络组件
func (a *App) ResetNetworkStack() error {
	exec.Command("netsh", "winsock", "reset").Run()
	exec.Command("netsh", "int", "ip", "reset").Run()
	return nil
}

// ReleaseRenewIP 释放并重新获取 IP
func (a *App) ReleaseRenewIP() error {
	exec.Command("ipconfig", "/release").Run()
	exec.Command("ipconfig", "/renew").Run()
	return nil
}

// getActiveAdapters 获取活动网卡名称
func (a *App) getActiveAdapters() []string {
	var adapters []string
	out, err := exec.Command("netsh", "interface", "show", "interface").Output()
	if err != nil {
		return adapters
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Connected") || strings.Contains(line, "已连接") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				adapters = append(adapters, fields[len(fields)-1])
			}
		}
	}
	return adapters
}

// RunTraceroute 路由追踪
func (a *App) RunTraceroute(target string) string {
	target = strings.TrimSpace(target)
	target = strings.TrimPrefix(target, "https://")
	target = strings.TrimPrefix(target, "http://")
	if idx := strings.Index(target, "/"); idx != -1 {
		target = target[:idx]
	}

	out, err := exec.Command("tracert", "-d", "-h", "15", target).Output()
	if err != nil {
		return "路由追踪失败: " + err.Error()
	}
	result, _ := simplifiedchinese.GBK.NewDecoder().String(string(out))
	return result
}

// CheckPort 端口检测
func (a *App) CheckPort(host string, port string) string {
	target := host + ":" + port
	out, err := exec.Command("powershell", "-Command",
		"$tcp = New-Object System.Net.Sockets.TcpClient; try { $tcp.Connect('"+host+"', "+port+"); 'open' } catch { 'closed' } finally { $tcp.Close() }").Output()
	if err != nil {
		return target + " - 检测失败"
	}
	result := strings.TrimSpace(string(out))
	if result == "open" {
		return target + " - ✓ 开放"
	}
	return target + " - ✗ 关闭/被封"
}

// GetNetworkInfo 获取网卡详细信息
func (a *App) GetNetworkInfo() string {
	out, err := exec.Command("ipconfig", "/all").Output()
	if err != nil {
		return "获取网络信息失败"
	}
	result, _ := simplifiedchinese.GBK.NewDecoder().String(string(out))
	return result
}

// RestartNetworkServices 重启网络服务
func (a *App) RestartNetworkServices() string {
	services := []string{"Dhcp", "Dnscache", "NlaSvc"}
	var results []string

	for _, svc := range services {
		exec.Command("net", "stop", svc).Run()
		err := exec.Command("net", "start", svc).Run()
		if err != nil {
			results = append(results, svc+": 重启失败")
		} else {
			results = append(results, svc+": 已重启")
		}
	}
	return strings.Join(results, "\n")
}

// GetFirewallStatus 获取防火墙状态
func (a *App) GetFirewallStatus() string {
	out, err := exec.Command("netsh", "advfirewall", "show", "allprofiles", "state").Output()
	if err != nil {
		return "获取防火墙状态失败"
	}
	result, _ := simplifiedchinese.GBK.NewDecoder().String(string(out))
	return result
}
