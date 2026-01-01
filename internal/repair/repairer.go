package repair

import (
	"context"
	"sync"

	"network-rescue-toolkit/pkg/types"
)

// Repairer 修复器接口
type Repairer interface {
	// ID 返回修复器唯一标识
	ID() string
	// Name 返回修复器名称
	Name() string
	// RequiresAdmin 是否需要管理员权限
	RequiresAdmin() bool
	// Repair 执行修复并返回结果
	Repair(ctx context.Context) types.RepairResult
}

// Engine 修复引擎
type Engine struct {
	repairers []Repairer
	results   []types.RepairResult
	mu        sync.RWMutex
}

// NewEngine 创建新的修复引擎
func NewEngine() *Engine {
	e := &Engine{
		repairers: make([]Repairer, 0),
		results:   make([]types.RepairResult, 0),
	}
	// 注册所有修复器
	e.registerDefaultRepairers()
	return e
}

// registerDefaultRepairers 注册默认修复器
func (e *Engine) registerDefaultRepairers() {
	e.RegisterRepairer(NewWinsockRepairer())
	e.RegisterRepairer(NewTCPIPRepairer())
	e.RegisterRepairer(NewDNSRepairer())
	e.RegisterRepairer(NewIPRepairer())
	e.RegisterRepairer(NewHostsRepairer())
	e.RegisterRepairer(NewProxyRepairer())
	e.RegisterRepairer(NewAdapterRepairer())
}

// RegisterRepairer 注册修复器
func (e *Engine) RegisterRepairer(r Repairer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.repairers = append(e.repairers, r)
}

// Repair 执行单个修复操作
func (e *Engine) Repair(ctx context.Context, id string) types.RepairResult {
	for _, repairer := range e.repairers {
		if repairer.ID() == id {
			return repairer.Repair(ctx)
		}
	}

	// 未找到修复器
	result := types.NewRepairResult(id, "未知修复项")
	result.SetFailure("未找到指定的修复项")
	return *result
}

// RepairAll 执行所有修复操作（综合修复）
func (e *Engine) RepairAll(ctx context.Context) []types.RepairResult {
	e.mu.Lock()
	e.results = make([]types.RepairResult, 0, len(e.repairers))
	e.mu.Unlock()

	for _, repairer := range e.repairers {
		select {
		case <-ctx.Done():
			return e.results
		default:
			result := repairer.Repair(ctx)
			e.mu.Lock()
			e.results = append(e.results, result)
			e.mu.Unlock()
		}
	}

	return e.results
}

// GetResults 获取最近一次修复结果
func (e *Engine) GetResults() []types.RepairResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.results
}

// GetRepairerIDs 获取所有修复器 ID
func (e *Engine) GetRepairerIDs() []string {
	ids := make([]string, len(e.repairers))
	for i, r := range e.repairers {
		ids[i] = r.ID()
	}
	return ids
}

// GetRepairerInfo 获取修复器信息
func (e *Engine) GetRepairerInfo(id string) (name string, requiresAdmin bool, found bool) {
	for _, r := range e.repairers {
		if r.ID() == id {
			return r.Name(), r.RequiresAdmin(), true
		}
	}
	return "", false, false
}
