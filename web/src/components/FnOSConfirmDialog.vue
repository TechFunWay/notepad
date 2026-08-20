<template>
  <teleport to="body">
    <div v-if="modelValue" class="fnos-confirm-mask" @click.self="handleCancel">
      <div class="fnos-confirm-card" role="dialog" aria-modal="true" aria-labelledby="fnos-confirm-title">
        <div class="fnos-confirm-badge">{{ initial }}</div>
        <h2 id="fnos-confirm-title">确认使用飞牛 NAS 登录</h2>
        <p class="fnos-confirm-desc">{{ appName }}正在请求使用下面的飞牛 NAS 账号登录</p>
        <div class="fnos-confirm-account">
          <span class="fnos-confirm-avatar">{{ initial }}</span>
          <span class="fnos-confirm-name">{{ username || '当前飞牛 NAS 用户' }}</span>
        </div>
        <el-button type="primary" size="large" class="fnos-confirm-btn" :loading="loading" @click="$emit('confirm')">
          {{ loading ? '正在登录…' : '确认登录' }}
        </el-button>
        <button type="button" class="fnos-confirm-link" :disabled="loading" @click="$emit('switch')">使用其他飞牛账号</button>
        <button type="button" class="fnos-confirm-link fnos-confirm-cancel" :disabled="loading" @click="handleCancel">取消</button>
      </div>
    </div>
  </teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  username: { type: String, default: '' },
  loading: { type: Boolean, default: false },
  appName: { type: String, default: '记事本' }
})
const emit = defineEmits(['update:modelValue', 'confirm', 'switch', 'cancel'])

const initial = computed(() => (props.username || '飞').slice(0, 1))

function handleCancel() {
  if (props.loading) return
  emit('update:modelValue', false)
  emit('cancel')
}
</script>

<style scoped>
.fnos-confirm-mask {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(3px);
}

.fnos-confirm-card {
  width: min(100%, 420px);
  padding: 32px;
  border-radius: 20px;
  background: var(--bg-primary, #fff);
  color: var(--text-primary, #1f2937);
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.fnos-confirm-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin: 0 auto;
  border-radius: 18px;
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  color: #fff;
  font-size: 26px;
  font-weight: 700;
  box-shadow: 0 8px 20px rgba(17, 153, 142, 0.35);
}

.fnos-confirm-card h2 {
  margin: 20px 0 8px;
  font-size: 22px;
  font-weight: 700;
}

.fnos-confirm-desc {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-secondary, #6b7280);
}

.fnos-confirm-account {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin: 20px 0;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-secondary, #f8fafc);
}

.fnos-confirm-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(17, 153, 142, 0.18);
  color: #11998e;
  font-size: 16px;
  font-weight: 600;
}

.fnos-confirm-name {
  font-size: 15px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fnos-confirm-btn {
  width: 100%;
  height: 48px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  border: none;
}

.fnos-confirm-link {
  display: block;
  width: 100%;
  margin-top: 12px;
  padding: 8px;
  border: none;
  background: transparent;
  color: #11998e;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.fnos-confirm-link:hover {
  opacity: 0.8;
}

.fnos-confirm-link:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.fnos-confirm-cancel {
  color: var(--text-secondary, #6b7280);
  font-weight: 500;
}
</style>
