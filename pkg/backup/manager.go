package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"network-rescue-toolkit/pkg/types"
)

// BackupInfo 备份信息
type BackupInfo struct {
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
	Size      int64     `json:"size"`
}

// Manager 备份管理器
type Manager struct {
	backupDir string
}

// NewManager 创建备份管理器
func NewManager() *Manager {
	// 默认备份目录在用户目录下
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".network-rescue-toolkit", "backups")
	os.MkdirAll(backupDir, 0755)

	return &Manager{
		backupDir: backupDir,
	}
}

// CreateBackup 创建配置备份
func (m *Manager) CreateBackup() (string, error) {
	config := types.NetworkConfig{
		// TODO: 收集当前网络配置
		Adapters:      []types.AdapterConfig{},
		DNSServers:    []string{},
		ProxySettings: types.ProxyConfig{},
		HostsContent:  "",
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("backup_%s.json", timestamp)
	filepath := filepath.Join(m.backupDir, filename)

	// 序列化并保存
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("保存备份文件失败: %w", err)
	}

	return filepath, nil
}


// RestoreBackup 还原配置备份
func (m *Manager) RestoreBackup(path string) error {
	// 读取备份文件
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %w", err)
	}

	// 解析配置
	var config types.NetworkConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析备份文件失败: %w", err)
	}

	// TODO: 应用配置
	_ = config

	return nil
}

// ListBackups 列出所有备份
func (m *Manager) ListBackups() ([]BackupInfo, error) {
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %w", err)
	}

	backups := make([]BackupInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			Path:      filepath.Join(m.backupDir, entry.Name()),
			Timestamp: info.ModTime(),
			Size:      info.Size(),
		})
	}

	return backups, nil
}

// DeleteBackup 删除备份
func (m *Manager) DeleteBackup(path string) error {
	return os.Remove(path)
}

// ValidateBackup 验证备份文件完整性
func (m *Manager) ValidateBackup(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %w", err)
	}

	var config types.NetworkConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("备份文件格式无效: %w", err)
	}

	return nil
}
