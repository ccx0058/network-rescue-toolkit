package types

// AdapterInfo 网络适配器信息
type AdapterInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MACAddress  string   `json:"macAddress"`
	Status      string   `json:"status"` // Up, Down, Unknown
	IPAddresses []string `json:"ipAddresses"`
	SubnetMasks []string `json:"subnetMasks"`
	Gateways    []string `json:"gateways"`
	DNSServers  []string `json:"dnsServers"`
	DHCPEnabled bool     `json:"dhcpEnabled"`
	DHCPServer  string   `json:"dhcpServer,omitempty"`
}

// IsValid 检查适配器信息是否有效
func (a *AdapterInfo) IsValid() bool {
	return a.Name != "" && a.MACAddress != ""
}

// HasIPAddress 检查是否有 IP 地址
func (a *AdapterInfo) HasIPAddress() bool {
	return len(a.IPAddresses) > 0
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled       bool   `json:"enabled"`
	Server        string `json:"server"`
	Port          int    `json:"port"`
	BypassList    string `json:"bypassList"`
	AutoConfigURL string `json:"autoConfigUrl,omitempty"`
}

// IsConfigured 检查代理是否已配置
func (p *ProxyConfig) IsConfigured() bool {
	return p.Enabled && p.Server != ""
}

// HostsEntry HOSTS 文件条目
type HostsEntry struct {
	IP         string `json:"ip"`
	Hostname   string `json:"hostname"`
	Comment    string `json:"comment,omitempty"`
	LineNum    int    `json:"lineNum"`
	Suspicious bool   `json:"suspicious"`
}

// SuspiciousDomains 可疑域名列表（常见被劫持的域名）
var SuspiciousDomains = []string{
	"www.baidu.com",
	"www.google.com",
	"www.qq.com",
	"www.taobao.com",
	"www.jd.com",
	"www.163.com",
	"www.sina.com.cn",
	"www.weibo.com",
	"www.alipay.com",
	"www.tmall.com",
}

// IsSuspicious 检查条目是否可疑
func (h *HostsEntry) IsSuspicious() bool {
	// localhost 相关的不算可疑
	if h.IP == "127.0.0.1" || h.IP == "::1" {
		if h.Hostname == "localhost" {
			return false
		}
	}

	// 检查是否是常见被劫持的域名
	for _, domain := range SuspiciousDomains {
		if h.Hostname == domain {
			// 如果不是指向 localhost，则可疑
			if h.IP != "127.0.0.1" && h.IP != "::1" {
				return true
			}
		}
	}

	return false
}

// NetworkConfig 网络配置快照（用于备份）
type NetworkConfig struct {
	Adapters      []AdapterConfig `json:"adapters"`
	DNSServers    []string        `json:"dnsServers"`
	ProxySettings ProxyConfig     `json:"proxySettings"`
	HostsContent  string          `json:"hostsContent"`
}

// AdapterConfig 适配器配置
type AdapterConfig struct {
	Name        string   `json:"name"`
	DHCPEnabled bool     `json:"dhcpEnabled"`
	IPAddresses []string `json:"ipAddresses,omitempty"`
	SubnetMasks []string `json:"subnetMasks,omitempty"`
	Gateways    []string `json:"gateways,omitempty"`
	DNSServers  []string `json:"dnsServers,omitempty"`
}

// ConnectivityResult 连通性测试结果
type ConnectivityResult struct {
	Target     string `json:"target"`
	Success    bool   `json:"success"`
	LatencyMs  int64  `json:"latencyMs"`
	Error      string `json:"error,omitempty"`
}

// PublicDNSServers 公共 DNS 服务器列表
var PublicDNSServers = map[string][]string{
	"114DNS":     {"114.114.114.114", "114.114.115.115"},
	"阿里DNS":    {"223.5.5.5", "223.6.6.6"},
	"腾讯DNS":    {"119.29.29.29", "182.254.116.116"},
	"百度DNS":    {"180.76.76.76"},
	"GoogleDNS": {"8.8.8.8", "8.8.4.4"},
	"Cloudflare": {"1.1.1.1", "1.0.0.1"},
}
