<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>找回密码</h2>
      <p class="subtitle">通过安全问题重置密码</p>

      <el-steps :active="step" finish-status="success" style="margin-bottom: 30px">
        <el-step title="输入用户名" />
        <el-step title="验证安全问题" />
        <el-step title="重置密码" />
      </el-steps>

      <el-form v-if="step === 0" :model="form">
        <el-form-item>
          <el-input v-model="form.username" placeholder="请输入用户名" prefix-icon="User" size="large" @keyup.enter="checkUsername" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="checkUsername" size="large" style="width: 100%">下一步</el-button>
        </el-form-item>
      </el-form>

      <el-form v-else-if="step === 1" :model="form">
        <div class="security-question-display">
          <el-icon><QuestionFilled /></el-icon>
          <span>{{ securityQuestion }}</span>
        </div>
        <el-form-item>
          <el-input v-model="form.security_answer" placeholder="请输入安全答案" size="large" @keyup.enter="verifyAnswer" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="verifyAnswer" size="large" style="width: 100%">下一步</el-button>
        </el-form-item>
      </el-form>

      <el-form v-else :model="form">
        <el-form-item>
          <el-input v-model="form.new_password" type="password" placeholder="新密码（至少6位）" prefix-icon="Lock" size="large" show-password @keyup.enter="resetPassword" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="resetPassword" size="large" style="width: 100%">重置密码</el-button>
        </el-form-item>
      </el-form>

      <div class="auth-links">
        <router-link to="/login">返回登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getSecurityQuestion, forgotPassword } from '../api/auth'
import { ElMessage } from 'element-plus'
import { md5 } from '@/utils/crypto'

const router = useRouter()
const loading = ref(false)
const step = ref(0)
const securityQuestion = ref('')
const form = ref({
  username: '',
  security_answer: '',
  new_password: ''
})

async function checkUsername() {
  if (!form.value.username) {
    ElMessage.warning('请输入用户名')
    return
  }
  loading.value = true
  try {
    const { data } = await getSecurityQuestion(form.value.username)
    securityQuestion.value = data.security_question
    step.value = 1
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '用户不存在')
  } finally {
    loading.value = false
  }
}

async function verifyAnswer() {
  if (!form.value.security_answer) {
    ElMessage.warning('请输入安全答案')
    return
  }
  step.value = 2
}

async function resetPassword() {
  if (!form.value.new_password || form.value.new_password.length < 6) {
    ElMessage.warning('密码至少6位')
    return
  }
  loading.value = true
  try {
    await forgotPassword({
      username: form.value.username,
      security_answer: form.value.security_answer,
      new_password: md5(form.value.new_password)
    })
    ElMessage.success('密码重置成功，请登录')
    router.push('/login')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '重置失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
  padding: 20px;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: white;
  padding: 40px 30px;
  border-radius: 16px;
  box-shadow: var(--card-shadow);
}

.auth-card h2 {
  text-align: center;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.auth-card .subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 30px;
  font-size: 14px;
}

.security-question-display {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
  color: #303133;
}

.auth-links {
  text-align: center;
  margin-top: 24px;
}

.auth-links a {
  color: var(--primary-color);
  text-decoration: none;
  font-size: 14px;
}

.auth-links a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .auth-card {
    padding: 30px 20px;
    border-radius: 12px;
  }
}
</style>