package privilege

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: network-rescue-toolkit, Property 17: Privilege Handling
// 测试权限检查逻辑的一致性
// **Validates: Requirements 17.1, 17.3, 17.4**

func TestIsAdminConsistency(t *testing.T) {
	// 属性测试：多次调用 IsAdmin 应该返回一致的结果
	helper := NewHelper()

	properties := gopter.NewProperties(gopter.DefaultTestParameters())

	properties.Property("IsAdmin returns consistent results", prop.ForAll(
		func(_ int) bool {
			// 多次调用应该返回相同结果
			result1 := helper.IsAdmin()
			result2 := helper.IsAdmin()
			result3 := helper.IsAdmin()
			return result1 == result2 && result2 == result3
		},
		gen.Int(), // 生成随机数只是为了触发多次测试
	))

	properties.TestingRun(t)
}

func TestHelperCreation(t *testing.T) {
	// 测试 Helper 创建
	helper := NewHelper()
	if helper == nil {
		t.Error("NewHelper should not return nil")
	}
}

func TestIsAdminReturnsBool(t *testing.T) {
	// 测试 IsAdmin 返回布尔值
	helper := NewHelper()
	result := helper.IsAdmin()
	// 结果应该是 true 或 false，不会 panic
	_ = result
}
