<script setup lang="ts">
import { ref } from 'vue'
import DiagnosticView from './DiagnosticView.vue'
import RepairView from './RepairView.vue'

const activeTab = ref<'diagnostic' | 'repair'>('diagnostic')
</script>

<template>
  <div class="bg-white rounded-lg shadow-md">
    <!-- 标签页切换 -->
    <div class="flex border-b">
      <button
        @click="activeTab = 'diagnostic'"
        :class="[
          'px-6 py-3 font-medium transition-colors',
          activeTab === 'diagnostic' 
            ? 'text-blue-600 border-b-2 border-blue-600 bg-blue-50' 
            : 'text-gray-600 hover:text-blue-600'
        ]"
      >
        🔍 网络诊断
      </button>
      <button
        @click="activeTab = 'repair'"
        :class="[
          'px-6 py-3 font-medium transition-colors',
          activeTab === 'repair' 
            ? 'text-blue-600 border-b-2 border-blue-600 bg-blue-50' 
            : 'text-gray-600 hover:text-blue-600'
        ]"
      >
        🛠️ 网络修复
      </button>
    </div>

    <!-- 内容区 -->
    <div class="p-6">
      <DiagnosticView v-if="activeTab === 'diagnostic'" />
      <RepairView v-else />
    </div>
  </div>
</template>
