<script setup lang="ts">
import { ref } from 'vue'

// 当前Tab
const currentTab = ref<'rescue' | 'tools'>('rescue')

// ========== 断网急救 ==========
interface DiagnosticItem {
  id: string
  name: string
  desc: string
  status: 'pending' | 'checking' | 'ok' | 'warning' | 'error'
  message: string
  repairable: boolean
}

const items = ref<DiagnosticItem[]>([
  { id: 'adapter', name: '网络硬件配置', desc: '检查网线是否插好，网卡电源及驱动是否正常工作', status: 'pending', message: '', repairable: false },
  { id: 'ip', name: '网络连接配置', desc: '检查网卡相关设置是否正确，IP地址是否配置正确', status: 'pending', message: '', repairable: false },
  { id: 'dns', name: 'DNS服务', desc: '如果您能上QQ，但打不开网页，往往是DNS服务出现问题', status: 'pending', message: '', repairable: false },
  { id: 'hosts', name: 'HOSTS', desc: '如果有些网页无法打开，往往是HOSTS出现问题', status: 'pending', message: '', repairable: false },
  { id: 'proxy', name: '浏览器配置', desc: '检查浏览器代理、插件等配置问题', status: 'pending', message: '', repairable: false },
  { id: 'connectivity', name: '电脑能否上网', desc: '检查您的电脑是否可以访问网页，网络是否连通', status: 'pending', message: '', repairable: false },
])

const isRunning = ref(false)
const allDone = ref(false)
const hasError = ref(false)
const statusText = ref('点击下方按钮开始全面诊断网络')
const progress = ref(0)

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'pending': return '○'
    case 'checking': return '◐'
    case 'ok': return '✓'
    case 'warning': return '⚠'
    case 'error': return '✗'
    default: return '○'
  }
}

const getStatusText = (item: DiagnosticItem) => {
  switch (item.status) {
    case 'pending': return '未诊断'
    case 'checking': return '诊断中'
    case 'ok': return '正常'
    case 'warning': return '警告'
    case 'error': return '异常'
    default: return ''
  }
}

const startDiagnosis = async () => {
  if (isRunning.value) return
  isRunning.value = true
  allDone.value = false
  hasError.value = false
  statusText.value = '正在进行全面网络诊断，请稍候....'
  items.value.forEach(item => { item.status = 'pending'; item.message = ''; item.repairable = false })

  for (let i = 0; i < items.value.length; i++) {
    items.value[i].status = 'checking'
    progress.value = ((i + 0.5) / items.value.length) * 100
    try {
      // @ts-ignore
      const result = await window.go.main.App.RunSingleDiagnostic(items.value[i].id)
      items.value[i].status = result.status
      items.value[i].message = result.message
      items.value[i].repairable = result.repairable
      if (result.status === 'error') hasError.value = true
    } catch (e) {
      items.value[i].status = 'error'
      items.value[i].message = '诊断失败'
      hasError.value = true
    }
    progress.value = ((i + 1) / items.value.length) * 100
  }

  isRunning.value = false
  allDone.value = true
  const errorCount = items.value.filter(i => i.status === 'error').length
  const warningCount = items.value.filter(i => i.status === 'warning').length
  const problemCount = errorCount + warningCount
  if (problemCount > 0) {
    hasError.value = true
    statusText.value = `诊断完成，发现 ${problemCount} 个问题，点击"立即修复"按钮修复`
  } else {
    hasError.value = false
    statusText.value = '诊断完成，您的网络一切正常！'
  }
}

const repairAll = async () => {
  if (isRunning.value) return
  isRunning.value = true
  statusText.value = '正在修复网络问题，请稍候....'
  try {
    // @ts-ignore
    await window.go.main.App.RunComprehensiveRepair()
    statusText.value = '修复完成，正在重新检测...'
    isRunning.value = false
    await startDiagnosis()
  } catch (e) {
    statusText.value = '修复过程中出现错误'
    isRunning.value = false
  }
}

const repairSingle = async (id: string) => {
  try {
    // @ts-ignore
    await window.go.main.App.RunRepair(id)
    const item = items.value.find(i => i.id === id)
    if (item) {
      item.status = 'checking'
      // @ts-ignore
      const result = await window.go.main.App.RunSingleDiagnostic(id)
      item.status = result.status
      item.message = result.message
      item.repairable = result.repairable
    }
  } catch (e) { console.error('修复失败:', e) }
}

// ========== 网络工具箱 ==========
const toolRunning = ref('')
const toolResult = ref('')
const pingTarget = ref('www.baidu.com')
const selectedDns = ref('114')

const dnsOptions = [
  { id: '114', name: '114 DNS', primary: '114.114.114.114', secondary: '114.114.115.115' },
  { id: 'ali', name: '阿里 DNS', primary: '223.5.5.5', secondary: '223.6.6.6' },
  { id: 'tencent', name: '腾讯 DNS', primary: '119.29.29.29', secondary: '182.254.118.118' },
  { id: 'baidu', name: '百度 DNS', primary: '180.76.76.76', secondary: '180.76.76.76' },
  { id: '360', name: '360 DNS', primary: '101.226.4.6', secondary: '218.30.118.6' },
  { id: 'cnnic', name: 'CNNIC DNS', primary: '1.2.4.8', secondary: '210.2.4.8' },
  { id: 'onedns', name: 'OneDNS 纯净', primary: '117.50.10.10', secondary: '52.80.52.52' },
  { id: 'dnspod', name: 'DNSPod', primary: '119.28.28.28', secondary: '119.29.29.29' },
  { id: 'google', name: 'Google DNS', primary: '8.8.8.8', secondary: '8.8.4.4' },
  { id: 'cloudflare', name: 'Cloudflare', primary: '1.1.1.1', secondary: '1.0.0.1' },
]

const runPing = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'ping'
  toolResult.value = `正在 Ping ${pingTarget.value} ...`
  try {
    // @ts-ignore
    const result = await window.go.main.App.RunPing(pingTarget.value)
    toolResult.value = result
  } catch (e) {
    toolResult.value = 'Ping 执行失败'
  }
  toolRunning.value = ''
}

const switchDns = async () => {
  if (toolRunning.value) return
  const dns = dnsOptions.find(d => d.id === selectedDns.value)
  if (!dns) return
  toolRunning.value = 'dns'
  toolResult.value = `正在切换到 ${dns.name} ...`
  try {
    // @ts-ignore
    const result = await window.go.main.App.SwitchDNS(dns.primary, dns.secondary)
    toolResult.value = result ? `已切换到 ${dns.name} (${dns.primary})` : 'DNS 切换失败'
  } catch (e) {
    toolResult.value = 'DNS 切换失败'
  }
  toolRunning.value = ''
}

const flushDns = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'flush'
  toolResult.value = '正在刷新 DNS 缓存...'
  try {
    // @ts-ignore
    await window.go.main.App.FlushDNS()
    toolResult.value = 'DNS 缓存已刷新'
  } catch (e) {
    toolResult.value = 'DNS 刷新失败'
  }
  toolRunning.value = ''
}

const resetNetwork = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'reset'
  toolResult.value = '正在重置网络组件...'
  try {
    // @ts-ignore
    await window.go.main.App.ResetNetworkStack()
    toolResult.value = '网络组件已重置，建议重启电脑'
  } catch (e) {
    toolResult.value = '网络重置失败'
  }
  toolRunning.value = ''
}

const releaseRenewIP = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'ip'
  toolResult.value = '正在释放并重新获取 IP 地址...'
  try {
    // @ts-ignore
    await window.go.main.App.ReleaseRenewIP()
    toolResult.value = 'IP 地址已重新获取'
  } catch (e) {
    toolResult.value = 'IP 操作失败'
  }
  toolRunning.value = ''
}

// 新增工具
const traceTarget = ref('www.baidu.com')
const portHost = ref('www.baidu.com')
const portNumber = ref('443')

const runTraceroute = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'trace'
  toolResult.value = `正在追踪路由到 ${traceTarget.value} ...\n（可能需要1-2分钟）`
  try {
    // @ts-ignore
    const result = await window.go.main.App.RunTraceroute(traceTarget.value)
    toolResult.value = result
  } catch (e) {
    toolResult.value = '路由追踪失败'
  }
  toolRunning.value = ''
}

const checkPort = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'port'
  toolResult.value = `正在检测 ${portHost.value}:${portNumber.value} ...`
  try {
    // @ts-ignore
    const result = await window.go.main.App.CheckPort(portHost.value, portNumber.value)
    toolResult.value = result
  } catch (e) {
    toolResult.value = '端口检测失败'
  }
  toolRunning.value = ''
}

const getNetworkInfo = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'info'
  toolResult.value = '正在获取网卡信息...'
  try {
    // @ts-ignore
    const result = await window.go.main.App.GetNetworkInfo()
    toolResult.value = result
  } catch (e) {
    toolResult.value = '获取网卡信息失败'
  }
  toolRunning.value = ''
}

const restartServices = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'services'
  toolResult.value = '正在重启网络服务...'
  try {
    // @ts-ignore
    const result = await window.go.main.App.RestartNetworkServices()
    toolResult.value = result
  } catch (e) {
    toolResult.value = '重启服务失败'
  }
  toolRunning.value = ''
}

const getFirewallStatus = async () => {
  if (toolRunning.value) return
  toolRunning.value = 'firewall'
  toolResult.value = '正在获取防火墙状态...'
  try {
    // @ts-ignore
    const result = await window.go.main.App.GetFirewallStatus()
    toolResult.value = result
  } catch (e) {
    toolResult.value = '获取防火墙状态失败'
  }
  toolRunning.value = ''
}
</script>

<template>
  <div class="app">
    <header class="header">
      <div class="header-left">
        <span class="logo">🔧</span>
        <span class="title">网络急救工具箱</span>
      </div>
      <div class="tabs">
        <button :class="['tab', { active: currentTab === 'rescue' }]" @click="currentTab = 'rescue'">断网急救</button>
        <button :class="['tab', { active: currentTab === 'tools' }]" @click="currentTab = 'tools'">网络工具</button>
      </div>
    </header>

    <!-- 断网急救页面 -->
    <main v-if="currentTab === 'rescue'" class="main">
      <div class="status-area">
        <div class="status-icon"><span class="icon">🖥️</span></div>
        <div class="status-text">
          <div class="status-title">{{ statusText }}</div>
          <div v-if="isRunning" class="progress-bar"><div class="progress-fill" :style="{ width: progress + '%' }"></div></div>
        </div>
        <button v-if="!isRunning && !allDone" class="btn-action btn-primary" @click="startDiagnosis">全面诊断</button>
        <button v-else-if="!isRunning && allDone && hasError" class="btn-action btn-primary" @click="repairAll">立即修复</button>
        <button v-else-if="!isRunning && allDone && !hasError" class="btn-action btn-secondary" @click="startDiagnosis">重新诊断</button>
        <button v-else class="btn-action" style="background: #9e9e9e; color: white;" disabled>诊断中...</button>
      </div>
      <div class="items-list">
        <div v-for="item in items" :key="item.id" class="item">
          <div :class="['item-icon', item.status]">{{ getStatusIcon(item.status) }}</div>
          <div class="item-content">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-desc">{{ item.message || item.desc }}</div>
          </div>
          <span v-if="(item.status === 'error' || item.status === 'warning') && item.repairable">
            <button class="btn-repair" @click="repairSingle(item.id)">修复</button>
          </span>
          <span v-else :class="['item-status', item.status]">{{ getStatusText(item) }}</span>
        </div>
      </div>
    </main>

    <!-- 网络工具页面 -->
    <main v-else class="main tools-page">
      <div class="tools-grid">
        <!-- Ping 工具 -->
        <div class="tool-card">
          <div class="tool-header"><span>📡</span> 网络 Ping 测试</div>
          <div class="tool-body">
            <input v-model="pingTarget" placeholder="输入域名或IP" class="tool-input" />
            <button class="tool-btn" @click="runPing" :disabled="!!toolRunning">
              {{ toolRunning === 'ping' ? '测试中...' : 'Ping' }}
            </button>
          </div>
        </div>

        <!-- DNS 切换 -->
        <div class="tool-card">
          <div class="tool-header"><span>🌐</span> 一键切换 DNS</div>
          <div class="tool-body">
            <select v-model="selectedDns" class="tool-select">
              <option v-for="dns in dnsOptions" :key="dns.id" :value="dns.id">{{ dns.name }}</option>
            </select>
            <button class="tool-btn" @click="switchDns" :disabled="!!toolRunning">
              {{ toolRunning === 'dns' ? '切换中...' : '切换' }}
            </button>
          </div>
        </div>

        <!-- 刷新 DNS -->
        <div class="tool-card">
          <div class="tool-header"><span>🔄</span> 刷新 DNS 缓存</div>
          <div class="tool-body">
            <p class="tool-desc">清除本地 DNS 缓存，解决域名解析问题</p>
            <button class="tool-btn full" @click="flushDns" :disabled="!!toolRunning">
              {{ toolRunning === 'flush' ? '刷新中...' : '刷新 DNS' }}
            </button>
          </div>
        </div>

        <!-- 重置网络 -->
        <div class="tool-card">
          <div class="tool-header"><span>⚡</span> 重置网络组件</div>
          <div class="tool-body">
            <p class="tool-desc">重置 Winsock 和 TCP/IP 协议栈</p>
            <button class="tool-btn full warning" @click="resetNetwork" :disabled="!!toolRunning">
              {{ toolRunning === 'reset' ? '重置中...' : '重置网络' }}
            </button>
          </div>
        </div>

        <!-- 释放续约 IP -->
        <div class="tool-card">
          <div class="tool-header"><span>🔃</span> 释放/续约 IP</div>
          <div class="tool-body">
            <p class="tool-desc">重新从 DHCP 服务器获取 IP 地址</p>
            <button class="tool-btn full" @click="releaseRenewIP" :disabled="!!toolRunning">
              {{ toolRunning === 'ip' ? '执行中...' : '重新获取 IP' }}
            </button>
          </div>
        </div>

        <!-- 路由追踪 -->
        <div class="tool-card">
          <div class="tool-header"><span>🛤️</span> 路由追踪</div>
          <div class="tool-body">
            <input v-model="traceTarget" placeholder="输入域名或IP" class="tool-input" />
            <button class="tool-btn" @click="runTraceroute" :disabled="!!toolRunning">
              {{ toolRunning === 'trace' ? '追踪中...' : 'Tracert' }}
            </button>
          </div>
        </div>

        <!-- 端口检测 -->
        <div class="tool-card">
          <div class="tool-header"><span>🔌</span> 端口检测</div>
          <div class="tool-body">
            <input v-model="portHost" placeholder="域名/IP" class="tool-input" style="flex:2" />
            <input v-model="portNumber" placeholder="端口" class="tool-input" style="flex:1;min-width:60px" />
            <button class="tool-btn" @click="checkPort" :disabled="!!toolRunning">
              {{ toolRunning === 'port' ? '检测中...' : '检测' }}
            </button>
          </div>
        </div>

        <!-- 网卡信息 -->
        <div class="tool-card">
          <div class="tool-header"><span>📋</span> 网卡详细信息</div>
          <div class="tool-body">
            <p class="tool-desc">查看 IP、MAC、网关、DNS 等详细配置</p>
            <button class="tool-btn full" @click="getNetworkInfo" :disabled="!!toolRunning">
              {{ toolRunning === 'info' ? '获取中...' : '查看详情' }}
            </button>
          </div>
        </div>

        <!-- 重启网络服务 -->
        <div class="tool-card">
          <div class="tool-header"><span>🔧</span> 重启网络服务</div>
          <div class="tool-body">
            <p class="tool-desc">重启 DHCP、DNS 缓存等系统服务</p>
            <button class="tool-btn full warning" @click="restartServices" :disabled="!!toolRunning">
              {{ toolRunning === 'services' ? '重启中...' : '重启服务' }}
            </button>
          </div>
        </div>

        <!-- 防火墙状态 -->
        <div class="tool-card">
          <div class="tool-header"><span>🛡️</span> 防火墙状态</div>
          <div class="tool-body">
            <p class="tool-desc">查看 Windows 防火墙当前状态</p>
            <button class="tool-btn full" @click="getFirewallStatus" :disabled="!!toolRunning">
              {{ toolRunning === 'firewall' ? '获取中...' : '查看状态' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 结果显示 -->
      <div v-if="toolResult" class="tool-result">
        <pre>{{ toolResult }}</pre>
      </div>
    </main>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'Microsoft YaHei', sans-serif; }
.app { height: 100vh; display: flex; flex-direction: column; background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%); }
.header { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; background: #4caf50; color: white; }
.header-left { display: flex; align-items: center; gap: 8px; }
.logo { font-size: 20px; }
.title { font-size: 14px; font-weight: 500; }
.tabs { display: flex; gap: 4px; }
.tab { background: rgba(255,255,255,0.2); border: none; color: white; padding: 6px 16px; border-radius: 4px; cursor: pointer; font-size: 13px; }
.tab:hover { background: rgba(255,255,255,0.3); }
.tab.active { background: white; color: #4caf50; }

.main { flex: 1; padding: 20px; display: flex; flex-direction: column; overflow: hidden; }
.status-area { display: flex; align-items: center; gap: 20px; padding: 20px; background: white; border-radius: 12px; margin-bottom: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.status-icon { width: 80px; height: 80px; display: flex; align-items: center; justify-content: center; }
.status-icon .icon { font-size: 56px; }
.status-text { flex: 1; }
.status-title { font-size: 18px; color: #333; margin-bottom: 8px; }
.progress-bar { height: 4px; background: #e0e0e0; border-radius: 2px; margin-top: 12px; overflow: hidden; }
.progress-fill { height: 100%; background: #4caf50; transition: width 0.3s; }
.btn-action { padding: 12px 36px; font-size: 15px; border: none; border-radius: 6px; cursor: pointer; font-weight: 500; }
.btn-primary { background: #ff5722; color: white; }
.btn-primary:hover { background: #f4511e; }
.btn-secondary { background: #4caf50; color: white; }
.btn-secondary:hover { background: #43a047; }

.items-list { flex: 1; background: white; border-radius: 12px; padding: 8px 0; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow-y: auto; }
.item { display: flex; align-items: center; padding: 14px 20px; border-bottom: 1px solid #f0f0f0; }
.item:last-child { border-bottom: none; }
.item-icon { width: 28px; height: 28px; margin-right: 16px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.item-icon.pending { color: #9e9e9e; }
.item-icon.checking { color: #2196f3; animation: spin 1s linear infinite; }
.item-icon.ok { color: #4caf50; }
.item-icon.warning { color: #ff9800; }
.item-icon.error { color: #f44336; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.item-content { flex: 1; }
.item-name { font-size: 14px; color: #333; font-weight: 500; }
.item-desc { font-size: 12px; color: #999; margin-top: 2px; }
.item-status { font-size: 13px; padding: 4px 12px; border-radius: 4px; }
.item-status.pending { color: #9e9e9e; }
.item-status.ok { color: #4caf50; }
.item-status.warning { color: #ff9800; }
.item-status.error { color: #f44336; }
.btn-repair { background: #ff5722; color: white; border: none; padding: 6px 16px; border-radius: 4px; cursor: pointer; font-size: 12px; }
.btn-repair:hover { background: #f4511e; }

/* 工具页面样式 */
.tools-page { overflow-y: auto; }
.tools-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
.tool-card { background: white; border-radius: 12px; padding: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.tool-header { font-size: 14px; font-weight: 500; color: #333; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
.tool-header span { font-size: 18px; }
.tool-body { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.tool-input, .tool-select { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 13px; min-width: 120px; }
.tool-btn { padding: 8px 16px; background: #4caf50; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 13px; white-space: nowrap; }
.tool-btn:hover { background: #43a047; }
.tool-btn:disabled { background: #9e9e9e; cursor: not-allowed; }
.tool-btn.full { width: 100%; }
.tool-btn.warning { background: #ff9800; }
.tool-btn.warning:hover { background: #f57c00; }
.tool-desc { font-size: 12px; color: #666; margin-bottom: 8px; width: 100%; }
.tool-result { margin-top: 16px; background: #263238; border-radius: 8px; padding: 16px; }
.tool-result pre { color: #4caf50; font-family: Consolas, monospace; font-size: 12px; white-space: pre-wrap; word-break: break-all; }
</style>
