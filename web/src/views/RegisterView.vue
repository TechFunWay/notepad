<template>
  <div class="register-container">
    <div class="register-bg"></div>
    <div class="register-wrapper">
      <div class="register-card">
        <div class="register-left">
          <div class="brand-section">
            <div class="brand-icon">
              <el-icon :size="56"><EditPen /></el-icon>
            </div>
            <h1 class="brand-title">{{ isFnOSBinding ? (isSetup ? '设置管理员' : '绑定飞牛 NAS') : (isSetup ? '设置管理员' : '创建账号') }}</h1>
            <p class="brand-desc">{{ isFnOSBinding ? (isSetup ? '创建管理员账号并绑定飞牛登录' : '关联应用账号，之后即可使用 NAS 登录') : (isSetup ? '首次使用，请设置管理员账号' : '开始记录您的每一个灵感') }}</p>
          </div>
          <div class="welcome-content">
            <div class="welcome-icon">
              <el-icon :size="80"><Star /></el-icon>
            </div>
            <h3>{{ isSetup ? '欢迎使用记事本' : '加入我们' }}</h3>
            <p>{{ isSetup ? '设置管理员后即可开始使用' : '创建一个账号，享受完整的笔记体验' }}</p>
            <div class="benefits">
              <div class="benefit-item">
                <el-icon><CircleCheck /></el-icon>
                <span>数据云端同步</span>
              </div>
              <div class="benefit-item">
                <el-icon><CircleCheck /></el-icon>
                <span>多设备支持</span>
              </div>
              <div class="benefit-item">
                <el-icon><CircleCheck /></el-icon>
                <span>安全加密存储</span>
              </div>
            </div>
          </div>
        </div>
        <div class="register-right">
          <div class="form-brand-mobile" aria-hidden="true">
            <span class="form-brand-mark"><el-icon><EditPen /></el-icon></span>
            <span>记事本</span>
          </div>
          <div class="form-header">
            <h2>{{ isFnOSBinding ? (isSetup ? '设置管理员并绑定飞牛' : '绑定飞牛 NAS 账号') : (isSetup ? '设置管理员' : '开始注册') }}</h2>
            <p>{{ isFnOSBinding ? (isSetup ? '请设置管理员账号和密码，完成后将自动与飞牛 NAS 登录绑定' : `当前飞牛用户 ${fnosUsername || '已登录用户'} 尚未绑定应用账号`) : (isSetup ? '请设置管理员账号和密码' : '只需几步即可完成注册') }}</p>
          </div>
          <el-form :model="form" class="register-form" label-width="0" @submit.prevent>
            <el-form-item>
              <label class="auth-label" for="register-username">用户名</label>
              <div class="input-wrapper">
                <el-icon class="input-icon"><User /></el-icon>
                <el-input 
                  id="register-username"
                  v-model="form.username" 
                  placeholder="请输入管理员用户名" 
                  aria-label="用户名"
                  autocomplete="username"
                  size="large"
                  class="custom-input"
                  @keyup.enter="focusPassword"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <label class="auth-label" for="register-password">密码</label>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  id="register-password"
                  v-model="form.password" 
                  type="password" 
                  placeholder="请输入管理员密码" 
                  aria-label="密码"
                  autocomplete="new-password"
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="focusConfirmPassword"
                />
              </div>
            </el-form-item>
            <el-form-item v-if="!isFnOSBinding || fnosMode === 'register'">
              <label class="auth-label" for="register-password-confirm">确认密码</label>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  id="register-password-confirm"
                  v-model="form.password_confirm" 
                  type="password" 
                  placeholder="请再次输入密码" 
                  aria-label="确认密码"
                  autocomplete="new-password"
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="focusSecurityQuestion"
                />
              </div>
            </el-form-item>
            <el-form-item v-if="!isFnOSBinding || fnosMode === 'register'">
              <label class="auth-label" for="register-question">安全问题</label>
              <div class="input-wrapper">
                <el-icon class="input-icon"><QuestionFilled /></el-icon>
                <el-input 
                  id="register-question"
                  v-model="form.security_question" 
                  placeholder="安全问题（用于找回密码）" 
                  aria-label="安全问题"
                  size="large"
                  class="custom-input"
                  @keyup.enter="focusSecurityAnswer"
                />
              </div>
            </el-form-item>
            <el-form-item v-if="!isFnOSBinding || fnosMode === 'register'">
              <label class="auth-label" for="register-answer">安全答案</label>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Key /></el-icon>
                <el-input 
                  id="register-answer"
                  v-model="form.security_answer" 
                  type="password" 
                  placeholder="安全答案" 
                  aria-label="安全答案"
                  autocomplete="off"
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="handleRegister"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary" 
                size="large" 
                class="submit-btn" 
                :loading="loading" 
                @click="handleRegister"
              >
                <span>{{ isFnOSBinding ? (fnosMode === 'register' ? '创建并绑定' : '验证并绑定') : (isSetup ? '设置管理员' : '创建账号') }}</span>
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </el-form-item>
            <el-button v-if="isFnOSBinding && !isSetup" text type="primary" class="switch-bind-btn" @click="fnosMode = fnosMode === 'register' ? 'bind' : 'register'">
              {{ fnosMode === 'register' ? '已有应用账号？验证并绑定' : '没有应用账号？创建并绑定' }}
            </el-button>
          </el-form>
          <div v-if="fnosEnabled && isSetup && !isFnOSBinding" class="fnos-setup">
            <div class="auth-divider"><span>或使用飞牛 NAS 账号</span></div>
            <el-button
              plain
              size="large"
              class="fnos-btn"
              :loading="fnosLoading"
              :disabled="loading || fnosLoading"
              @click="handleFnOSSetup"
            >
              使用飞牛 NAS 账号设置管理员
            </el-button>
            <p>授权后将进入设置管理员表单，输入管理员账号密码并与飞牛账号绑定。</p>
          </div>
          <div class="form-footer" v-if="!isSetup && !isFnOSBinding">
            <div class="login-prompt">
              <span>已有账号？</span>
              <router-link to="/login" class="login-link">立即登录</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
    <FnOSConfirmDialog
      v-model="fnosConfirmVisible"
      :username="fnosConfirmUsername"
      :loading="fnosLoading"
      @confirm="confirmFnOSSetup"
      @switch="switchFnOSAccount"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from '@/utils/message'
import { EditPen, User, Lock, QuestionFilled, Key, ArrowRight, Star, CircleCheck } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { md5 } from '@/utils/crypto'
import { markSetupComplete } from '@/utils/setupStatus'
import { openFnOSAuthPopup } from '@/utils/fnosAuthPopup'
import FnOSConfirmDialog from '@/components/FnOSConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const fnosLoading = ref(false)
const fnosEnabled = import.meta.env.VITE_FNOS_APP === 'true'
const fnosConfirmVisible = ref(false)
const fnosConfirmUsername = ref('')

const isSetup = computed(() => route.query.setup === 'true')
const isFnOSBinding = computed(() => import.meta.env.VITE_FNOS_APP === 'true' && route.query.fnos === 'bind')
const fnosUsername = computed(() => typeof route.query.fnos_username === 'string' ? route.query.fnos_username : '')
const fnosMode = ref('register')

const form = reactive({
  username: '',
  password: '',
  password_confirm: '',
  security_question: '',
  security_answer: ''
})

async function focusPassword() {
  await nextTick()
  const inputs = document.querySelectorAll('.custom-input')
  if (inputs[1]) {
    const input = inputs[1].querySelector('input')
    input?.focus()
  }
}

async function focusConfirmPassword() {
  await nextTick()
  const inputs = document.querySelectorAll('.custom-input')
  if (inputs[2]) {
    const input = inputs[2].querySelector('input')
    input?.focus()
  }
}

async function focusSecurityQuestion() {
  await nextTick()
  const inputs = document.querySelectorAll('.custom-input')
  if (inputs[3]) {
    const input = inputs[3].querySelector('input')
    input?.focus()
  }
}

async function focusSecurityAnswer() {
  await nextTick()
  const inputs = document.querySelectorAll('.custom-input')
  if (inputs[4]) {
    const input = inputs[4].querySelector('input')
    input?.focus()
  }
}

async function handleRegister() {
  if (!form.username || !form.password || ((!isFnOSBinding.value || fnosMode.value === 'register') && !form.password_confirm)) {
    message.warning('请填写必要信息')
    return
  }
  if ((!isFnOSBinding.value || fnosMode.value === 'register') && form.password !== form.password_confirm) {
    message.warning('两次输入的密码不一致')
    return
  }
  loading.value = true
  try {
    if (isFnOSBinding.value) {
      await auth.bindFnOS(fnosMode.value, form.username, md5(form.password))
      if (isSetup.value) markSetupComplete()
      message.success(isSetup.value ? '管理员设置成功' : '飞牛 NAS 账号绑定成功')
      router.push('/')
      return
    }
    await auth.register({
      username: form.username,
      password: md5(form.password),
      security_question: form.security_question,
      security_answer: md5(form.security_answer)
    })
    if (isSetup.value) markSetupComplete()
    message.success(isSetup.value ? '管理员设置成功' : '注册成功')
    router.push('/')
  } catch (e) {
    message.error(e.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}

async function handleFnOSSetup() {
  fnosLoading.value = true
  try {
    const data = await auth.getFnOSIdentity()
    fnosConfirmUsername.value = data.fnos_username || ''
    fnosConfirmVisible.value = true
  } catch (e) {
    message.error(e.response?.data?.error || e.message || '无法获取飞牛 NAS 账号')
  } finally {
    fnosLoading.value = false
  }
}

async function confirmFnOSSetup() {
  fnosLoading.value = true
  try {
    const result = await auth.fnosLogin()
    fnosConfirmVisible.value = false
    if (!result.binding_required) {
      // 该飞牛账号已绑定应用账号，直接用飞牛账号登录
      markSetupComplete()
      message.success('飞牛 NAS 登录成功')
      router.push('/')
      return
    }
    // 未绑定：进入创建管理员表单，输入应用管理员账号密码后与飞牛账号绑定
    router.push({
      path: '/register',
      query: { setup: 'true', fnos: 'bind', fnos_username: result.fnos_username || '' }
    })
  } catch (e) {
    fnosConfirmVisible.value = false
    message.error(e.response?.data?.error || e.message || '飞牛 NAS 授权失败')
  } finally {
    fnosLoading.value = false
  }
}

async function switchFnOSAccount() {
  fnosConfirmVisible.value = false
  const status = await openFnOSAuthPopup()
  if (status === 'blocked') {
    message.warning('浏览器拦截了登录弹窗，请允许本网站弹出窗口后重试')
    return
  }
  // 完成或取消后，重新拉取当前（可能已切换的）飞牛账号，并再次弹出确认窗口
  await handleFnOSSetup()
}

if (isFnOSBinding.value) {
  form.username = fnosUsername.value
}

</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.register-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.register-bg::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  animation: rotate 30s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.register-wrapper {
  position: relative;
  z-index: 1;
  padding: 20px;
  width: 100%;
  max-width: 1100px;
}

.register-card {
  background: var(--bg-primary);
  border-radius: 24px;
  box-shadow: 0 30px 80px -15px rgba(0, 0, 0, 0.3);
  display: flex;
  overflow: hidden;
  min-height: 650px;
}

.register-left {
  flex: 1.2;
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.register-left::before {
  content: '';
  position: absolute;
  top: -20%;
  right: -20%;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
}

.register-left::after {
  content: '';
  position: absolute;
  bottom: -10%;
  left: -10%;
  width: 200px;
  height: 200px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 50%;
}

.brand-section {
  position: relative;
  z-index: 1;
}

.brand-icon {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
}

.brand-title {
  font-size: 36px;
  font-weight: 700;
  color: white;
  margin: 0 0 12px;
  letter-spacing: -0.5px;
}

.brand-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
  line-height: 1.6;
}

.welcome-content {
  position: relative;
  z-index: 1;
  text-align: center;
}

.welcome-icon {
  color: white;
  margin-bottom: 20px;
  opacity: 0.9;
}

.welcome-content h3 {
  font-size: 24px;
  font-weight: 600;
  color: white;
  margin: 0 0 8px;
}

.welcome-content p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 32px;
}

.benefits {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-start;
}

.benefit-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 14px;
}

.register-right {
  flex: 1;
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-header {
  margin-bottom: 32px;
}

.form-header h2 {
  font-size: 30px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.form-header p {
  font-size: 15px;
  color: #6b7280;
  margin: 0;
}

.register-form {
  margin-bottom: 24px;
}

.fnos-setup {
  margin-top: 24px;
}

.fnos-btn {
  width: 100%;
  min-height: 48px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
}

.fnos-setup > p {
  margin: 8px 0 20px;
  color: #6b7280;
  font-size: 13px;
  line-height: 1.6;
  text-align: center;
}

.auth-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #9ca3af;
  font-size: 13px;
}

.auth-divider::before,
.auth-divider::after {
  content: '';
  height: 1px;
  flex: 1;
  background: #e5e7eb;
}

.input-wrapper {
  position: relative;
  width: 100%;
}

.custom-input {
  width: 100%;
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 20px;
  z-index: 1;
}

.custom-input :deep(.el-input__wrapper) {
  padding-left: 48px !important;
  border-radius: 12px;
  box-shadow: 0 0 0 1px #e5e7eb;
  transition: all 0.3s;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #11998e;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #11998e;
}

.submit-btn {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(17, 153, 142, 0.4);
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(17, 153, 142, 0.5);
}

.submit-btn:active {
  transform: translateY(0);
}

.form-footer {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}

.login-prompt {
  display: flex;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
}

.login-link {
  color: #11998e;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s;
}

.login-link:hover {
  color: #0e7e75;
  text-decoration: underline;
}

@media (max-width: 968px) {
  .register-left {
    display: none;
  }

  .register-card {
    max-width: 480px;
    margin: 0 auto;
    min-height: auto;
  }

  .register-right {
    padding: 48px 32px;
  }
}
</style>
