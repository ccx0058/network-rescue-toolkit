package types

import "time"

// DiagnosticStatus 诊断状态
type DiagnosticStatus string

const (
	StatusOK      DiagnosticStatus = "ok"
	StatusWarning DiagnosticStatus = "warning"
	StatusError   DiagnosticStatus = "error"
)

// DiagnosticResult 诊断结果
type DiagnosticResult struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Status     DiagnosticStatus  `json:"status"`
	Message    string            `json:"message"`
	Details    map[string]any    `json:"details"`
	Timestamp  time.Time         `json:"timestamp"`
	Repairable bool              `json:"repairable"`
}

// NewDiagnosticResult 创建新的诊断结果
func NewDiagnosticResult(id, name string) *DiagnosticResult {
	return &DiagnosticResult{
		ID:        id,
		Name:      name,
		Status:    StatusOK,
		Details:   make(map[string]any),
		Timestamp: time.Now(),
	}
}

// SetOK 设置为正常状态
func (r *DiagnosticResult) SetOK(message string) {
	r.Status = StatusOK
	r.Message = message
	r.Repairable = false
}

// SetWarning 设置为警告状态
func (r *DiagnosticResult) SetWarning(message string, repairable bool) {
	r.Status = StatusWarning
	r.Message = message
	r.Repairable = repairable
}

// SetError 设置为错误状态
func (r *DiagnosticResult) SetError(message string, repairable bool) {
	r.Status = StatusError
	r.Message = message
	r.Repairable = repairable
}

// AddDetail 添加详情
func (r *DiagnosticResult) AddDetail(key string, value any) {
	r.Details[key] = value
}

// DiagnosticReport 诊断报告
type DiagnosticReport struct {
	GeneratedAt time.Time          `json:"generatedAt"`
	SystemInfo  SystemInfo         `json:"systemInfo"`
	Results     []DiagnosticResult `json:"results"`
	Summary     ReportSummary      `json:"summary"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OSVersion    string `json:"osVersion"`
	ComputerName string `json:"computerName"`
	Username     string `json:"username"`
}

// ReportSummary 报告摘要
type ReportSummary struct {
	TotalChecks   int `json:"totalChecks"`
	PassedChecks  int `json:"passedChecks"`
	WarningChecks int `json:"warningChecks"`
	FailedChecks  int `json:"failedChecks"`
}

// CalculateSummary 计算报告摘要
func (r *DiagnosticReport) CalculateSummary() {
	r.Summary.TotalChecks = len(r.Results)
	r.Summary.PassedChecks = 0
	r.Summary.WarningChecks = 0
	r.Summary.FailedChecks = 0

	for _, result := range r.Results {
		switch result.Status {
		case StatusOK:
			r.Summary.PassedChecks++
		case StatusWarning:
			r.Summary.WarningChecks++
		case StatusError:
			r.Summary.FailedChecks++
		}
	}
}
