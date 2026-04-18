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
            <h1 class="brand-title">创建账号</h1>
            <p class="brand-desc">开始记录您的每一个灵感</p>
          </div>
          <div class="welcome-content">
            <div class="welcome-icon">
              <el-icon :size="80"><Star /></el-icon>
            </div>
            <h3>加入我们</h3>
            <p>创建一个账号，享受完整的笔记体验</p>
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
          <div class="form-header">
            <h2>开始注册</h2>
            <p>只需几步即可完成注册</p>
          </div>
          <el-form :model="form" class="register-form" label-width="0">
            <el-form-item>
              <div class="input-wrapper">
                <el-icon class="input-icon"><User /></el-icon>
                <el-input 
                  v-model="form.username" 
                  placeholder="请输入用户名" 
                  size="large"
                  class="custom-input"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.password" 
                  type="password" 
                  placeholder="请输入密码" 
                  size="large"
                  show-password
                  class="custom-input"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.password_confirm" 
                  type="password" 
                  placeholder="请再次输入密码" 
                  size="large"
                  show-password
                  class="custom-input"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <div class="input-wrapper">
                <el-icon class="input-icon"><QuestionFilled /></el-icon>
                <el-input 
                  v-model="form.security_question" 
                  placeholder="安全问题（用于找回密码）" 
                  size="large"
                  class="custom-input"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Key /></el-icon>
                <el-input 
                  v-model="form.security_answer" 
                  type="password" 
                  placeholder="安全答案" 
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
                <span>创建账号</span>
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </el-form-item>
          </el-form>
          <div class="form-footer">
            <div class="login-prompt">
              <span>已有账号？</span>
              <router-link to="/login" class="login-link">立即登录</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message } from '@/utils/message'
import { EditPen, User, Lock, QuestionFilled, Key, ArrowRight, Star, CircleCheck } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { md5 } from '@/utils/crypto'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  password_confirm: '',
  security_question: '',
  security_answer: ''
})

async function handleRegister() {
  if (!form.username || !form.password || !form.password_confirm) {
    message.warning('请填写必要信息')
    return
  }
  if (form.password !== form.password_confirm) {
    message.warning('两次输入的密码不一致')
    return
  }
  loading.value = true
  try {
    await auth.register({
      username: form.username,
      password: md5(form.password),
      security_question: form.security_question,
      security_answer: md5(form.security_answer)
    })
    message.success('注册成功')
    router.push('/')
  } catch (e) {
    message.error(e.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
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
  background: white;
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
    padding: 32px 24px;
    justify-content: flex-start;
    gap: 16px;
  }

  .form-header {
    margin-bottom: 24px;
  }

  .register-form {
    margin-bottom: 0;
  }
}

@media (max-width: 480px) {
  .register-wrapper {
    padding: 0;
    max-width: 100%;
  }

  .register-card {
    border-radius: 0;
    min-height: 100vh;
    max-width: none;
    width: 100%;
  }

  .register-container {
    align-items: stretch;
  }

  .register-right {
    padding: 24px 16px 16px;
    gap: 0;
    width: 100%;
    box-sizing: border-box;
    flex: 1;
  }

  .form-header {
    margin-bottom: 20px;
    width: 100%;
  }

  .form-header h2 {
    font-size: 24px;
  }

  .form-header p {
    font-size: 13px;
  }

  .register-form {
    margin-bottom: 0;
    width: 100%;
  }

  /* Force form items to be full-width block elements */
  .register-form :deep(.el-form-item) {
    margin-bottom: 10px;
    display: block !important;
    width: 100% !important;
  }

  .register-form :deep(.el-form-item:last-child) {
    margin-bottom: 16px;
  }

  /* Remove label space completely */
  .register-form :deep(.el-form-item__content) {
    margin-left: 0 !important;
    display: block !important;
    width: 100% !important;
  }

  /* Target the actual input and button elements */
  .register-form :deep(.el-form-item .el-input) {
    width: 100% !important;
  }

  .register-form :deep(.el-form-item .el-input__wrapper) {
    width: 100% !important;
  }

  .register-form :deep(.el-form-item .el-button) {
    width: 100% !important;
  }

  .input-wrapper {
    width: 100%;
    box-sizing: border-box;
  }

  .custom-input {
    width: 100%;
  }

  .custom-input :deep(.el-input) {
    width: 100% !important;
  }

  .custom-input :deep(.el-input__wrapper) {
    width: 100% !important;
  }

  .submit-btn {
    margin-top: 0;
    height: 48px;
    width: 100% !important;
  }

  .form-footer {
    margin-top: 12px;
    padding-top: 0;
    width: 100%;
  }

  .login-prompt {
    font-size: 13px;
  }
}
</style>
