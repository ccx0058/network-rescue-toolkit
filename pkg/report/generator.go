package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"network-rescue-toolkit/pkg/types"
)

// Generator 报告生成器
type Generator struct {
	outputDir string
}

// NewGenerator 创建报告生成器
func NewGenerator() *Generator {
	homeDir, _ := os.UserHomeDir()
	outputDir := filepath.Join(homeDir, ".network-rescue-toolkit", "reports")
	os.MkdirAll(outputDir, 0755)

	return &Generator{
		outputDir: outputDir,
	}
}

// Generate 生成报告
func (g *Generator) Generate(results []types.DiagnosticResult, format string) (string, error) {
	report := types.DiagnosticReport{
		GeneratedAt: time.Now(),
		SystemInfo:  g.getSystemInfo(),
		Results:     results,
	}
	report.CalculateSummary()

	switch format {
	case "json":
		return g.generateJSON(report)
	case "html":
		return g.generateHTML(report)
	default:
		return "", fmt.Errorf("不支持的格式: %s", format)
	}
}

// getSystemInfo 获取系统信息
func (g *Generator) getSystemInfo() types.SystemInfo {
	hostname, _ := os.Hostname()
	username := os.Getenv("USERNAME")

	return types.SystemInfo{
		OSVersion:    "Windows",
		ComputerName: hostname,
		Username:     username,
	}
}


// generateJSON 生成 JSON 格式报告
func (g *Generator) generateJSON(report types.DiagnosticReport) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("report_%s.json", timestamp)
	filepath := filepath.Join(g.outputDir, filename)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化报告失败: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("保存报告失败: %w", err)
	}

	return filepath, nil
}

// generateHTML 生成 HTML 格式报告
func (g *Generator) generateHTML(report types.DiagnosticReport) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("report_%s.html", timestamp)
	filepath := filepath.Join(g.outputDir, filename)

	tmpl := template.Must(template.New("report").Parse(htmlTemplate))

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("创建报告文件失败: %w", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, report)
	if err != nil {
		return "", fmt.Errorf("生成报告失败: %w", err)
	}

	return filepath, nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>网络诊断报告</title>
    <style>
        body { font-family: 'Microsoft YaHei', sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #2563eb; padding-bottom: 10px; }
        .summary { display: flex; gap: 20px; margin: 20px 0; }
        .summary-item { padding: 15px; border-radius: 8px; flex: 1; text-align: center; }
        .summary-ok { background: #dcfce7; color: #166534; }
        .summary-warning { background: #fef3c7; color: #92400e; }
        .summary-error { background: #fee2e2; color: #991b1b; }
        .result { padding: 15px; margin: 10px 0; border-radius: 8px; border-left: 4px solid; }
        .result-ok { background: #f0fdf4; border-color: #22c55e; }
        .result-warning { background: #fffbeb; border-color: #f97316; }
        .result-error { background: #fef2f2; border-color: #ef4444; }
        .result h3 { margin: 0 0 5px 0; }
        .result p { margin: 5px 0; color: #666; }
        .info { color: #666; font-size: 14px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔧 网络诊断报告</h1>
        <p class="info">生成时间: {{.GeneratedAt.Format "2006-01-02 15:04:05"}} | 计算机: {{.SystemInfo.ComputerName}} | 用户: {{.SystemInfo.Username}}</p>
        
        <div class="summary">
            <div class="summary-item summary-ok">
                <div style="font-size: 24px; font-weight: bold;">{{.Summary.PassedChecks}}</div>
                <div>正常</div>
            </div>
            <div class="summary-item summary-warning">
                <div style="font-size: 24px; font-weight: bold;">{{.Summary.WarningChecks}}</div>
                <div>警告</div>
            </div>
            <div class="summary-item summary-error">
                <div style="font-size: 24px; font-weight: bold;">{{.Summary.FailedChecks}}</div>
                <div>错误</div>
            </div>
        </div>

        <h2>诊断详情</h2>
        {{range .Results}}
        <div class="result result-{{.Status}}">
            <h3>{{.Name}}</h3>
            <p>{{.Message}}</p>
        </div>
        {{end}}
    </div>
</body>
</html>`
