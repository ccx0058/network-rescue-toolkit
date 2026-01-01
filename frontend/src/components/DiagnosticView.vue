<script setup lang="ts">
import { ref } from 'vue'

interface DiagnosticResult {
  id: string
  name: string
  status: 'ok' | 'warning' | 'error'
  message: string
  details: Record<string, any>
  timestamp: string
  repairable: boolean
}

const results = ref<DiagnosticResult[]>([])
const running = ref(false)

const runDiagnostic = async () => {
  running.value = true
  results.value = []
  try {
    // @ts-ignore - Wails 运行时注入
    results.value = await window.go.main.App.RunDiagnostic()
  } catch (e) {
    console.error('诊断失败:', e)
  }
  running.value = false
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'ok': return '✓'
    case 'warning': return '⚠'
    case 'error': return '✗'
    default: return '?'
  }
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'ok': return 'text-status-ok bg-green-50 border-green-200'
    case 'warning': return 'text-status-warning bg-orange-50 border-orange-200'
    case 'error': return 'text-status-error bg-red-50 border-red-200'
    default: return 'text-gray-500 bg-gray-50 border-gray-200'
  }
}

const exportReport = async (format: string) => {
  try {
    // @ts-ignore
    const path = await window.go.main.App.ExportReport(format)
    alert(`报告已导出到: ${path}`)
  } catch (e) {
    console.error('导出失败:', e)
  }
}
</script>

<template>
  <div>
    <!-- 操作按钮 -->
    <div class="flex gap-4 mb-6">
      <button
        @click="runDiagnostic"
        :disabled="running"
        class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white px-6 py-2 rounded-lg font-medium transition-colors"
      >
        {{ running ? '诊断中...' : '🔍 全面诊断' }}
      </button>
      <button
        v-if="results.length > 0"
        @click="exportReport('html')"
        class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg text-sm"
      >
        📄 导出 HTML 报告
      </button>
      <button
        v-if="results.length > 0"
        @click="exportReport('json')"
        class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg text-sm"
      >
        📋 导出 JSON
      </button>
    </div>

    <!-- 诊断结果列表 -->
    <div v-if="results.length > 0" class="space-y-3">
      <div
        v-for="result in results"
        :key="result.id"
        :class="['border rounded-lg p-4', getStatusClass(result.status)]"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span class="text-2xl">{{ getStatusIcon(result.status) }}</span>
            <div>
              <h3 class="font-medium">{{ result.name }}</h3>
              <p class="text-sm opacity-75">{{ result.message }}</p>
            </div>
          </div>
          <span v-if="result.repairable" class="text-xs bg-white px-2 py-1 rounded border">
            可修复
          </span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!running" class="text-center py-12 text-gray-500">
      <p class="text-4xl mb-4">🔍</p>
      <p>点击"全面诊断"按钮开始检测网络状态</p>
    </div>

    <!-- 加载状态 -->
    <div v-else class="text-center py-12 text-gray-500">
      <p class="text-4xl mb-4 animate-pulse">⏳</p>
      <p>正在诊断网络，请稍候...</p>
    </div>
  </div>
</template>
