<template>
  <div class="fnos-callback">
    <h1>获取当前登录信息</h1>
    <p>已获取当前飞牛 NAS 登录状态，正在返回记事本…</p>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { FNOS_AUTH_RESULT_TYPE } from '@/utils/fnosAuthPopup'

onMounted(() => {
  try {
    if (window.opener) {
      window.opener.postMessage({ type: FNOS_AUTH_RESULT_TYPE }, '*')
    }
  } catch (e) {
    // 忽略跨窗口通信异常
  }
  setTimeout(() => {
    window.close()
  }, 500)
})
</script>

<style scoped>
.fnos-callback {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px;
  text-align: center;
  background: var(--bg-primary, #fff);
  color: var(--text-primary, #1f2937);
}
.fnos-callback h1 {
  margin: 0;
  font-size: 20px;
}
.fnos-callback p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary, #6b7280);
}
</style>
