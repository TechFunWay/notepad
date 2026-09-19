<template>
  <!-- 主弹窗：赞赏说明 + 收款码 -->
  <el-dialog
    :model-value="support.show"
    width="400px"
    align-center
    append-to-body
    :close-on-click-modal="true"
    @update:model-value="onMainShowChange"
  >
    <template #header>
      <div class="support-header">
        <div class="support-emoji">☕</div>
        <h3 class="support-title">请作者喝杯咖啡</h3>
      </div>
    </template>

    <div class="support-body">
      <p class="support-text">
        这个应用免费、无广告，数据完全保存在你自己的设备上。
      </p>
      <p class="support-text">
        如果它帮到了你，欢迎请作者喝杯咖啡——<strong>金额随意，1 元也是心意</strong>。
        <br />
        <span class="support-sub">不赞赏也完全没有问题，<strong>不支付不影响任何功能</strong>。</span>
      </p>

      <div class="support-qr-wrap">
        <img
          :src="donateQr"
          alt="微信赞赏码"
          width="352"
          height="480"
          decoding="async"
          class="support-qr"
        />
        <span class="support-qr-tip">微信扫码赞赏</span>
      </div>

      <p v-if="support.errorText" class="support-error">{{ support.errorText }}</p>
    </div>

    <template #footer>
      <div class="support-actions">
        <el-button :disabled="support.sending" @click="support.dismiss()">暂不支持</el-button>
        <el-button type="primary" :disabled="support.sending" :loading="support.sending" @click="openAmount">
          ❤ {{ support.sending ? '发送中…' : '已支持' }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 金额输入弹窗（点击「已支持」后出现） -->
  <el-dialog
    v-model="amountDialog"
    title="填写赞赏金额"
    width="340px"
    align-center
    append-to-body
    :close-on-click-modal="!support.sending"
  >
    <div class="support-amount-body">
      <p class="support-amount-sub">金额随意，1 元也是心意；仅用于接收端统计，不做支付核验</p>
      <div class="support-presets">
        <el-button
          v-for="p in presets"
          :key="p"
          size="small"
          round
          :type="amount === p ? 'primary' : 'default'"
          :plain="amount !== p"
          @click="selectPreset(p)"
        >
          {{ p }} 元
        </el-button>
      </div>
      <div class="support-custom-row">
        <span class="support-custom-label">或输入</span>
        <el-input-number
          v-model="customAmount"
          :min="0"
          :step="0.01"
          :controls="false"
          placeholder="金额（元）"
          class="support-input"
          @update:model-value="onCustomInput"
        />
      </div>
      <p v-if="amount === 0" class="support-zero">未填写金额将以 0 元上报</p>
      <p v-if="support.errorText" class="support-error">{{ support.errorText }}</p>
    </div>
    <template #footer>
      <div class="support-actions">
        <el-button :disabled="support.sending" @click="closeAmount">返回</el-button>
        <el-button type="primary" :disabled="support.sending" :loading="support.sending" @click="confirmAmount">
          {{ support.sending ? '发送中…' : '确定' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElDialog, ElButton, ElInputNumber } from 'element-plus'
import { useSupportStore } from '@/stores/support'
import donateQr from '@/assets/donate-wechat.png'

const support = useSupportStore()

// 金额输入弹窗状态
const amountDialog = ref(false)
const presets = [1, 5, 10, 50]
const customAmount = ref(null)
const amount = ref(0)

// 主弹窗关闭时（含关闭按钮/遮罩/发送成功）一并复位金额弹窗
watch(
  () => support.show,
  visible => {
    if (!visible) resetAmount()
  }
)

function onMainShowChange(visible) {
  if (!visible) support.dismiss()
}

function resetAmount() {
  amountDialog.value = false
  customAmount.value = null
  amount.value = 0
}

function openAmount() {
  resetAmount()
  support.errorText = ''
  amountDialog.value = true
}

function selectPreset(p) {
  amount.value = p
  customAmount.value = p
}

function onCustomInput(v) {
  amount.value = Number(v) || 0
}

function closeAmount() {
  if (support.sending) return
  resetAmount()
}

function confirmAmount() {
  if (support.sending) return
  void support.confirmSupported(amount.value || 0)
}
</script>

<style scoped>
.support-header {
  text-align: center;
}

.support-emoji {
  width: 54px;
  height: 54px;
  margin: 0 auto 10px;
  display: grid;
  place-items: center;
  font-size: 26px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.18), rgba(249, 115, 22, 0.18));
}

.support-title {
  margin: 0;
  font-size: 17px;
  font-weight: 800;
  color: var(--text-primary, #1f2937);
}

.support-body {
  text-align: center;
}

.support-text {
  margin: 0 0 10px;
  color: var(--text-secondary, #6b7280);
  font-size: 13.5px;
  line-height: 1.7;
}

.support-text strong {
  color: var(--text-primary, #1f2937);
}

.support-sub {
  font-size: 12px;
}

.support-qr-wrap {
  margin: 14px 0 4px;
}

.support-qr {
  display: block;
  margin: 0 auto;
  width: 190px;
  height: auto;
  max-height: 32vh;
  object-fit: contain;
  padding: 4px;
  background: #fff;
  border: 1px solid var(--border-light, #e5e7eb);
  border-radius: 12px;
}

.support-qr-tip {
  display: block;
  margin-top: 8px;
  color: var(--text-muted, #9ca3af);
  font-size: 12px;
}

.support-error {
  margin-top: 10px;
  color: var(--el-color-danger, #ef4444);
  font-size: 12px;
}

.support-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.support-amount-body {
  text-align: center;
}

.support-amount-sub {
  margin: 0;
  color: var(--text-muted, #9ca3af);
  font-size: 12px;
}

.support-presets {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
}

.support-custom-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.support-custom-label {
  color: var(--text-secondary, #6b7280);
  font-size: 13px;
}

.support-input {
  width: 120px;
}

.support-zero {
  margin-top: 8px;
  color: var(--el-color-warning, #f59e0b);
  font-size: 12px;
}
</style>
