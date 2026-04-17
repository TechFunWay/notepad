<template>
  <div class="login-container">
    <div class="login-bg"></div>
    <div class="login-wrapper">
      <div class="login-card">
        <div class="login-left">
          <div class="brand-section">
            <div class="brand-icon">
              <el-icon :size="56"><Document /></el-icon>
            </div>
            <h1 class="brand-title">记事本</h1>
            <p class="brand-desc">记录每一刻的灵感与想法</p>
          </div>
          <div class="features">
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon><EditPen /></el-icon>
              </div>
              <div class="feature-text">
                <h4>富文本编辑</h4>
                <p>支持多种格式和样式</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon><Tickets /></el-icon>
              </div>
              <div class="feature-text">
                <h4>标签管理</h4>
                <p>轻松组织和查找笔记</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="feature-text">
                <h4>响应式设计</h4>
                <p>完美适配各种设备</p>
              </div>
            </div>
          </div>
        </div>
        <div class="login-right">
          <div class="form-header">
            <h2>欢迎回来</h2>
            <p>登录您的账号继续使用</p>
          </div>
          <div class="login-form">
            <div class="form-item" style="margin-bottom: 20px">
              <div class="input-wrapper">
                <el-icon class="input-icon"><User /></el-icon>
                <el-input 
                  v-model="form.username" 
                  placeholder="请输入用户名" 
                  size="large"
                  class="custom-input"
                />
              </div>
            </div>
            <div class="form-item" style="margin-bottom: 24px">
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.password" 
                  type="password" 
                  placeholder="请输入密码" 
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="handleLogin"
                />
              </div>
            </div>
            <div class="form-options">
              <label class="remember-me">
                <el-checkbox v-model="rememberMe">记住账号密码</el-checkbox>
              </label>
            </div>
            <div class="form-item">
              <el-button 
                type="primary" 
                size="large" 
                class="submit-btn" 
                :loading="loading" 
                @click="handleLogin"
              >
                <span>登录</span>
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </div>
          <div class="form-footer">
            <router-link to="/forgot-password" class="forgot-link">忘记密码？</router-link>
            <div class="register-prompt">
              <span>还没有账号？</span>
              <template v-if="allowRegister">
                <router-link to="/register" class="register-link">立即注册</router-link>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, User, Lock, ArrowRight, EditPen, Tickets, Monitor } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useConfigStore } from '@/stores/config'

const router = useRouter()
const auth = useAuthStore()
const config = useConfigStore()
const loading = ref(false)
const allowRegister = ref(true)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: ''
})

onMounted(async () => {
  await config.fetchPublicConfig()
  allowRegister.value = config.allowRegister

  const saved = localStorage.getItem('login_remember')
  if (saved) {
    try {
      const data = JSON.parse(saved)
      form.username = data.username || ''
      form.password = data.password || ''
      if (form.username && form.password) {
        rememberMe.value = true
      }
    } catch (e) {
      localStorage.removeItem('login_remember')
    }
  }
})

async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    if (rememberMe.value) {
      localStorage.setItem('login_remember', JSON.stringify({
        username: form.username,
        password: form.password
      }))
    } else {
      localStorage.removeItem('login_remember')
    }
    ElMessage.success('登录成功')
    router.push('/')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.login-bg::before {
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

.login-wrapper {
  position: relative;
  z-index: 1;
  padding: 20px;
  width: 100%;
  max-width: 1100px;
}

.login-card {
  background: white;
  border-radius: 24px;
  box-shadow: 0 30px 80px -15px rgba(0, 0, 0, 0.3);
  display: flex;
  overflow: hidden;
  min-height: 600px;
}

.login-left {
  flex: 1.2;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.login-left::before {
  content: '';
  position: absolute;
  top: -20%;
  right: -20%;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
}

.login-left::after {
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

.features {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.feature-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.feature-text h4 {
  font-size: 16px;
  font-weight: 600;
  color: white;
  margin: 0 0 4px;
}

.feature-text p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.75);
  margin: 0;
}

.login-right {
  flex: 1;
  min-width: 0;
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 100%;
}

.form-header {
  margin-bottom: 40px;
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

.login-form {
  width: 100%;
  margin-bottom: 24px;
}

.form-item {
  width: 100%;
}

.input-wrapper {
  position: relative;
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
  box-shadow: 0 0 0 1px #667eea;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #667eea;
}

.submit-btn {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.submit-btn:active {
  transform: translateY(0);
}

.form-options {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-start;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
  font-size: 14px;
  color: #6b7280;
}

.remember-me :deep(.el-checkbox__inner) {
  border-radius: 4px;
}

.remember-me :deep(.el-checkbox__label) {
  font-size: 14px;
  color: #6b7280;
}

.form-footer {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}

.forgot-link {
  color: #6b7280;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: #667eea;
  text-decoration: underline;
}

.register-prompt {
  display: flex;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
}

.register-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s;
}

.register-link:hover {
  color: #764ba2;
  text-decoration: underline;
}

@media (max-width: 968px) {
  .login-left {
    display: none;
  }

  .login-card {
    max-width: 480px;
    margin: 0 auto;
  }

  .login-right {
    padding: 48px 32px;
  }
}

@media (max-width: 480px) {
  .login-wrapper {
    padding: 0;
  }

  .login-card {
    border-radius: 0;
    min-height: 100vh;
  }

  .login-right {
    padding: 32px 20px 20px;
  }

  .form-header {
    margin-bottom: 24px;
  }

  .form-header h2 {
    font-size: 24px;
  }

  .login-form {
    margin-bottom: 0;
  }

  .form-options {
    margin-bottom: 16px;
  }

  .form-footer {
    margin-top: 12px;
  }
}
</style>
