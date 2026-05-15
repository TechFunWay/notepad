<template>
  <div class="admin-container">
    <div class="admin-header">
      <h2>系统配置</h2>
    </div>

    <!-- Desktop table -->
    <el-table v-if="!isMobile" :data="configs" stripe style="width: 100%">
      <el-table-column prop="key" label="配置项" width="250" />
      <el-table-column prop="description" label="说明" />
      <el-table-column prop="value" label="当前值" width="200" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="showEditDialog(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Mobile card list -->
    <div v-else class="mobile-card-list">
      <div v-for="config in configs" :key="config.key" class="mobile-card">
        <div class="mobile-card-header">
          <span class="mobile-card-title">{{ config.description }}</span>
        </div>
        <div class="mobile-card-info">
          <span class="config-key">{{ config.key }}</span>
          <span class="config-value">{{ formatValue(config) }}</span>
        </div>
        <div class="mobile-card-actions">
          <el-button size="small" type="primary" @click="showEditDialog(config)">编辑</el-button>
        </div>
      </div>
    </div>

    <el-dialog v-model="editDialogVisible" title="编辑配置" :width="isMobile ? '95%' : '400px'">
      <el-form :label-width="isMobile ? 'auto' : '80px'">
        <el-form-item label="配置项">
          <el-input :model-value="editingConfig.key" disabled />
        </el-form-item>
        <el-form-item label="说明">
          <el-input :model-value="editingConfig.description" disabled />
        </el-form-item>
        <el-form-item label="值">
          <el-select v-if="editingConfig.key === 'allow_register'" v-model="editValue">
            <el-option label="允许" value="true" />
            <el-option label="禁止" value="false" />
          </el-select>
          <el-input v-else v-model="editValue" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getConfigs, updateConfig } from '../../api/config'
import { message } from '@/utils/message'

const configs = ref([])
const editDialogVisible = ref(false)
const editingConfig = ref({})
const editValue = ref('')
const isMobile = ref(false)

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadConfigs()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

async function loadConfigs() {
  const { data } = await getConfigs()
  configs.value = data.configs
}

function formatValue(config) {
  if (config.key === 'allow_register') {
    return config.value === 'true' ? '允许' : '禁止'
  }
  return config.value
}

function showEditDialog(config) {
  editingConfig.value = config
  editValue.value = config.value
  editDialogVisible.value = true
}

async function handleEdit() {
  try {
    await updateConfig(editingConfig.value.key, editValue.value)
    message.success('配置更新成功')
    editDialogVisible.value = false
    loadConfigs()
  } catch (e) {
    message.error(e.response?.data?.error || '更新失败')
  }
}
</script>

<style scoped>
.admin-container {
  padding: 24px;
  overflow-y: auto;
  height: 100%;
  -webkit-overflow-scrolling: touch;
  background: var(--bg-primary);
  border-radius: 16px;
  box-shadow: var(--card-shadow);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.admin-header h2 {
  margin: 0;
}

.mobile-card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mobile-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 14px 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.mobile-card-header {
  margin-bottom: 8px;
}

.mobile-card-title {
  font-weight: 500;
  font-size: 15px;
}

.mobile-card-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
  margin-bottom: 10px;
}

.config-key {
  color: var(--text-secondary);
  font-family: monospace;
  font-size: 12px;
}

.config-value {
  color: var(--text-primary);
  font-weight: 500;
}

.mobile-card-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 768px) {
  .admin-container {
    padding: 16px;
  }

  .admin-header h2 {
    font-size: 18px;
  }
}
</style>
