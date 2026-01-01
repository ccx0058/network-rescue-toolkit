package diagnostic

import (
	"context"
	"sync"

	"network-rescue-toolkit/pkg/types"
)

// Checker 诊断检查器接口
type Checker interface {
	// ID 返回检查器唯一标识
	ID() string
	// Name 返回检查器名称
	Name() string
	// Check 执行检查并返回结果
	Check(ctx context.Context) types.DiagnosticResult
}

// Engine 诊断引擎
type Engine struct {
	checkers []Checker
	results  []types.DiagnosticResult
	mu       sync.RWMutex
}

// NewEngine 创建新的诊断引擎
func NewEngine() *Engine {
	e := &Engine{
		checkers: make([]Checker, 0),
		results:  make([]types.DiagnosticResult, 0),
	}
	// 注册所有检查器
	e.registerDefaultCheckers()
	return e
}

// registerDefaultCheckers 注册默认检查器
func (e *Engine) registerDefaultCheckers() {
	e.RegisterChecker(NewAdapterChecker())
	e.RegisterChecker(NewIPChecker())
	e.RegisterChecker(NewDNSChecker())
	e.RegisterChecker(NewHostsChecker())
	e.RegisterChecker(NewProxyChecker())
	e.RegisterChecker(NewConnectivityChecker())
}

// RegisterChecker 注册检查器
func (e *Engine) RegisterChecker(c Checker) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.checkers = append(e.checkers, c)
}

// RunAll 运行所有检查器
func (e *Engine) RunAll(ctx context.Context) []types.DiagnosticResult {
	e.mu.Lock()
	e.results = make([]types.DiagnosticResult, 0, len(e.checkers))
	e.mu.Unlock()

	for _, checker := range e.checkers {
		select {
		case <-ctx.Done():
			return e.results
		default:
			result := checker.Check(ctx)
			e.mu.Lock()
			e.results = append(e.results, result)
			e.mu.Unlock()
		}
	}

	return e.results
}

// RunSingle 运行单个检查器
func (e *Engine) RunSingle(ctx context.Context, id string) types.DiagnosticResult {
	for _, checker := range e.checkers {
		if checker.ID() == id {
			return checker.Check(ctx)
		}
	}

	// 未找到检查器
	result := types.NewDiagnosticResult(id, "未知检查项")
	result.SetError("未找到指定的检查项", false)
	return *result
}

// GetResults 获取最近一次诊断结果
func (e *Engine) GetResults() []types.DiagnosticResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.results
}

// GetCheckerIDs 获取所有检查器 ID
func (e *Engine) GetCheckerIDs() []string {
	ids := make([]string, len(e.checkers))
	for i, c := range e.checkers {
		ids[i] = c.ID()
	}
	return ids
}
