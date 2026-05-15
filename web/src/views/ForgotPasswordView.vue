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

      <el-form v-if="step === 0" :model="form" @submit.prevent>
        <el-form-item>
          <el-input ref="usernameInput" v-model="form.username" placeholder="请输入用户名" prefix-icon="User" size="large" @keyup.enter="checkUsername" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="checkUsername" @keyup.enter="checkUsername" size="large" style="width: 100%">下一步</el-button>
        </el-form-item>
      </el-form>

      <el-form v-else-if="step === 1" :model="form" @submit.prevent>
        <div class="security-question-display">
          <el-icon><QuestionFilled /></el-icon>
          <span>{{ securityQuestion }}</span>
        </div>
        <el-form-item>
          <el-input ref="answerInput" v-model="form.security_answer" placeholder="请输入安全答案" size="large" @keyup.enter="verifyAnswer" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="verifyAnswer" size="large" style="width: 100%">下一步</el-button>
        </el-form-item>
      </el-form>

      <el-form v-else :model="form" @submit.prevent>
        <el-form-item>
          <el-input ref="passwordInput" v-model="form.new_password" type="password" placeholder="新密码（至少6位）" prefix-icon="Lock" size="large" show-password @keyup.enter="handlePasswordEnter" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.confirm_password" type="password" placeholder="请再次输入新密码" prefix-icon="Lock" size="large" show-password @keyup.enter="resetPassword" />
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
import { ref, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getSecurityQuestion, forgotPassword } from '../api/auth'
import api from '../api/request'
import { message } from '@/utils/message'
import { md5 } from '@/utils/crypto'

const router = useRouter()
const loading = ref(false)
const step = ref(0)
const securityQuestion = ref('')
const usernameInput = ref(null)
const answerInput = ref(null)
const passwordInput = ref(null)
const form = ref({
  username: '',
  security_answer: '',
  new_password: '',
  confirm_password: ''
})

onMounted(async () => {
  await nextTick()
  usernameInput.value?.focus()
})

async function checkUsername() {
  if (!form.value.username) {
    message.warning('请输入用户名')
    return
  }
  loading.value = true
  try {
    const { data } = await getSecurityQuestion(form.value.username)
    if (!data.security_question) {
      message.error('该用户未设置安全问题')
      return
    }
    securityQuestion.value = data.security_question
    step.value = 1
    await nextTick()
    form.value.security_answer = ''
    answerInput.value?.focus()
  } catch (e) {
    message.error(e.response?.data?.error || '用户不存在')
  } finally {
    loading.value = false
  }
}

async function verifyAnswer() {
  if (!form.value.security_answer) {
    message.warning('请输入安全答案')
    return
  }
  loading.value = true
  try {
    await api.post('/auth/verify-answer', {
      username: form.value.username,
      security_answer: md5(form.value.security_answer)
    })
    step.value = 2
    await nextTick()
    form.value.new_password = ''
    form.value.confirm_password = ''
    passwordInput.value?.focus()
  } catch (e) {
    message.error(e.response?.data?.error || '安全答案错误')
  } finally {
    loading.value = false
  }
}

function handlePasswordEnter() {
  if (!form.value.new_password || form.value.new_password.length < 6) {
    message.warning('密码至少6位')
    return
  }
  document.querySelector('input[placeholder="请再次输入新密码"]')?.focus()
}

async function resetPassword() {
  if (!form.value.new_password || form.value.new_password.length < 6) {
    message.warning('密码至少6位')
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    message.warning('两次输入的密码不一致')
    return
  }
  loading.value = true
  try {
    await forgotPassword({
      username: form.value.username,
      security_answer: md5(form.value.security_answer),
      new_password: md5(form.value.new_password)
    })
    message.success('密码重置成功')
    window.location.href = '/login'
  } catch (e) {
    message.error(e.response?.data?.error || '密码重置失败')
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
  background: var(--bg-primary);
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
  background: var(--bg-secondary);
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
