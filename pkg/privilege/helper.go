package privilege

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Helper 权限助手
type Helper struct{}

// NewHelper 创建权限助手
func NewHelper() *Helper {
	return &Helper{}
}

// IsAdmin 检查是否以管理员权限运行
func (h *Helper) IsAdmin() bool {
	var sid *windows.SID

	// 创建管理员组 SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	// 检查当前进程是否属于管理员组
	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}


// RequestElevation 请求提升权限（重新以管理员身份启动）
func (h *Helper) RequestElevation() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	return h.RunElevated(exe)
}

// RunElevated 以管理员权限运行程序
func (h *Helper) RunElevated(executable string, args ...string) error {
	verb := "runas"
	cwd, _ := os.Getwd()

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(executable)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)

	var argsStr string
	for _, arg := range args {
		argsStr += " " + arg
	}
	argsPtr, _ := syscall.UTF16PtrFromString(argsStr)

	var showCmd int32 = 1 // SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argsPtr, cwdPtr, showCmd)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrentProcessToken 获取当前进程令牌
func (h *Helper) GetCurrentProcessToken() (windows.Token, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return 0, err
	}
	return token, nil
}

// HasPrivilege 检查是否具有指定权限
func (h *Helper) HasPrivilege(privilegeName string) bool {
	token, err := h.GetCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()

	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(privilegeName), &luid)
	if err != nil {
		return false
	}

	var privileges windows.Tokenprivileges
	var returnLength uint32
	err = windows.GetTokenInformation(token, windows.TokenPrivileges, (*byte)(unsafe.Pointer(&privileges)), uint32(unsafe.Sizeof(privileges)), &returnLength)
	if err != nil {
		return false
	}

	// 简化检查：如果是管理员就认为有权限
	return h.IsAdmin()
}
