package types

import "time"

// RepairResult 修复结果
type RepairResult struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	RequireReboot bool      `json:"requireReboot"`
	Timestamp     time.Time `json:"timestamp"`
	BackupPath    string    `json:"backupPath,omitempty"`
}

// NewRepairResult 创建新的修复结果
func NewRepairResult(id, name string) *RepairResult {
	return &RepairResult{
		ID:        id,
		Name:      name,
		Timestamp: time.Now(),
	}
}

// SetSuccess 设置为成功
func (r *RepairResult) SetSuccess(message string) {
	r.Success = true
	r.Message = message
}

// SetFailure 设置为失败
func (r *RepairResult) SetFailure(message string) {
	r.Success = false
	r.Message = message
}

// SetRequireReboot 设置需要重启
func (r *RepairResult) SetRequireReboot() {
	r.RequireReboot = true
}

// SetBackupPath 设置备份路径
func (r *RepairResult) SetBackupPath(path string) {
	r.BackupPath = path
}

// RepairReport 修复报告
type RepairReport struct {
	GeneratedAt time.Time      `json:"generatedAt"`
	Results     []RepairResult `json:"results"`
	Summary     RepairSummary  `json:"summary"`
}

// RepairSummary 修复摘要
type RepairSummary struct {
	TotalRepairs     int  `json:"totalRepairs"`
	SuccessfulRepairs int  `json:"successfulRepairs"`
	FailedRepairs    int  `json:"failedRepairs"`
	RequiresReboot   bool `json:"requiresReboot"`
}

// CalculateSummary 计算修复摘要
func (r *RepairReport) CalculateSummary() {
	r.Summary.TotalRepairs = len(r.Results)
	r.Summary.SuccessfulRepairs = 0
	r.Summary.FailedRepairs = 0
	r.Summary.RequiresReboot = false

	for _, result := range r.Results {
		if result.Success {
			r.Summary.SuccessfulRepairs++
		} else {
			r.Summary.FailedRepairs++
		}
		if result.RequireReboot {
			r.Summary.RequiresReboot = true
		}
	}
}
