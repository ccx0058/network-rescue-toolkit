package registry

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// RegistryHelper 注册表操作助手
type RegistryHelper struct{}

// NewRegistryHelper 创建注册表助手
func NewRegistryHelper() *RegistryHelper {
	return &RegistryHelper{}
}

// ReadString 读取字符串值
func (h *RegistryHelper) ReadString(root registry.Key, path, name string) (string, error) {
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	if err != nil {
		return "", fmt.Errorf("读取注册表值失败: %w", err)
	}

	return value, nil
}

// ReadDWORD 读取 DWORD 值
func (h *RegistryHelper) ReadDWORD(root registry.Key, path, name string) (uint32, error) {
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue(name)
	if err != nil {
		return 0, fmt.Errorf("读取注册表值失败: %w", err)
	}

	return uint32(value), nil
}


// WriteString 写入字符串值
func (h *RegistryHelper) WriteString(root registry.Key, path, name, value string) error {
	key, err := registry.OpenKey(root, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	err = key.SetStringValue(name, value)
	if err != nil {
		return fmt.Errorf("写入注册表值失败: %w", err)
	}

	return nil
}

// WriteDWORD 写入 DWORD 值
func (h *RegistryHelper) WriteDWORD(root registry.Key, path, name string, value uint32) error {
	key, err := registry.OpenKey(root, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	err = key.SetDWordValue(name, value)
	if err != nil {
		return fmt.Errorf("写入注册表值失败: %w", err)
	}

	return nil
}

// DeleteValue 删除注册表值
func (h *RegistryHelper) DeleteValue(root registry.Key, path, name string) error {
	key, err := registry.OpenKey(root, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(name)
	if err != nil {
		return fmt.Errorf("删除注册表值失败: %w", err)
	}

	return nil
}

// KeyExists 检查注册表键是否存在
func (h *RegistryHelper) KeyExists(root registry.Key, path string) bool {
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	key.Close()
	return true
}

// ValueExists 检查注册表值是否存在
func (h *RegistryHelper) ValueExists(root registry.Key, path, name string) bool {
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	_, _, err = key.GetStringValue(name)
	return err == nil
}

// 常用注册表路径
const (
	// ProxySettingsPath IE 代理设置路径
	ProxySettingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`
)
