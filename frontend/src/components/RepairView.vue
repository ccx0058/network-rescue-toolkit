<script setup lang="ts">
import { ref } from 'vue'

interface RepairOption {
  id: string
  name: string
  description: string
  requiresAdmin: boolean
  icon: string
}

const repairOptions: RepairOption[] = [
  { id: 'winsock', name: '重置 Winsock', description: '修复网络套接字和 LSP 问题', requiresAdmin: true, icon: '🔌' },
  { id: 'tcpip', name: '重置 TCP/IP', description: '重置 TCP/IP 协议栈', requiresAdmin: true, icon: '🌐' },
  { id: 'dns', name: '刷新 DNS 缓存', description: '清除 DNS 解析缓存', requiresAdmin: false, icon: '📡' },
  { id: 'ip', name: '释放/续租 IP', description: '重新获取 DHCP 分配的 IP 地址', requiresAdmin: true, icon: '🔄' },
  { id: 'hosts', name: '修复 HOSTS 文件', description: '恢复 HOSTS 文件为默认状态', requiresAdmin: true, icon: '📝' },
  { id: 'proxy', name: '清除代理设置', description: '禁用系统代理服务器', requiresAdmin: false, icon: '🚫' },
  { id: 'adapter', name: '重置网络适配器', description: '禁用后重新启用网卡', requiresAdmin: true, icon: '💻' },
]

const repairing = ref<string | null>(null)
const results = ref<Record<string, { success: boolean; message: string }>>({})

const runRepair = async (id: string) => {
  repairing.value = id
  try {
    // @ts-ignore - Wails 运行时注入
    const result = await window.go.main.App.RunRepair(id)
    results.value[id] = { success: result.success, message: result.message }
  } catch (e: any) {
    results.value[id] = { success: false, message: e.message || '修复失败' }
  }
  repairing.value = null
}

const runComprehensiveRepair = async () => {
  repairing.value = 'comprehensive'
  try {
    // @ts-ignore
    await window.go.main.App.RunComprehensiveRepair()
    alert('综合修复完成！建议重启电脑以使所有更改生效。')
  } catch (e) {
    console.error('综合修复失败:', e)
  }
  repairing.value = null
}
</script>

<template>
  <div>
    <!-- 综合修复按钮 -->
    <div class="mb-6 p-4 bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg text-white">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="font-bold text-lg">⚡ 一键综合修复</h3>
          <p class="text-sm opacity-90">执行所有修复操作，彻底解决网络问题</p>
        </div>
        <button
          @click="runComprehensiveRepair"
          :disabled="repairing !== null"
          class="bg-white text-blue-600 hover:bg-blue-50 disabled:opacity-50 px-6 py-2 rounded-lg font-medium"
        >
          {{ repairing === 'comprehensive' ? '修复中...' : '开始修复' }}
        </button>
      </div>
    </div>

    <!-- 单项修复选项 -->
    <h3 class="font-medium text-gray-700 mb-4">单项修复</h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="option in repairOptions"
        :key="option.id"
        class="border rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div class="flex gap-3">
            <span class="text-2xl">{{ option.icon }}</span>
            <div>
              <h4 class="font-medium">
                {{ option.name }}
                <span v-if="option.requiresAdmin" class="text-xs text-orange-500 ml-1">[需要管理员]</span>
              </h4>
              <p class="text-sm text-gray-500">{{ option.description }}</p>
              <!-- 修复结果 -->
              <p v-if="results[option.id]" :class="results[option.id].success ? 'text-green-600' : 'text-red-600'" class="text-sm mt-1">
                {{ results[option.id].success ? '✓' : '✗' }} {{ results[option.id].message }}
              </p>
            </div>
          </div>
          <button
            @click="runRepair(option.id)"
            :disabled="repairing !== null"
            class="bg-gray-100 hover:bg-gray-200 disabled:opacity-50 px-3 py-1 rounded text-sm"
          >
            {{ repairing === option.id ? '修复中...' : '修复' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
