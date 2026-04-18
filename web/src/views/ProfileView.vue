<template>
  <div class="profile-container">
    <div class="profile-layout">
      <div class="profile-sidebar">
        <div class="sidebar-header">
          <div class="avatar">
            <el-icon :size="32"><UserFilled /></el-icon>
          </div>
          <h2>个人设置</h2>
        </div>
        <div class="menu-list">
          <div 
            class="menu-item" 
            :class="{ active: activeTab === 'password' }"
            @click="activeTab = 'password'"
          >
            <el-icon><Key /></el-icon>
            <span>密码修改</span>
          </div>
          <div 
            class="menu-item" 
            :class="{ active: activeTab === 'security' }"
            @click="activeTab = 'security'"
          >
            <el-icon><QuestionFilled /></el-icon>
            <span>安全问题</span>
          </div>
        </div>
      </div>

      <div class="profile-content">
        <div v-if="activeTab === 'password'" class="tab-panel">
          <div class="panel-header">
            <el-icon :size="40" class="panel-icon"><Key /></el-icon>
            <div>
              <h3>密码修改</h3>
              <p>修改您的登录密码</p>
            </div>
          </div>
          <el-form :model="passwordForm" class="form-wrapper" @submit.prevent>
            <div class="form-group">
              <label class="group-label">当前密码</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><Lock /></el-icon>
                <el-input 
                  v-model="passwordForm.current_password" 
                  type="password" 
                  placeholder="请输入当前密码"
                  show-password
                  @keyup.enter="focusNewPassword"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="group-label">新密码</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><Lock /></el-icon>
                <el-input 
                  v-model="passwordForm.new_password" 
                  type="password" 
                  placeholder="请输入新密码（至少6位）"
                  show-password
                  @keyup.enter="focusConfirmPassword"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="group-label">确认密码</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><Lock /></el-icon>
                <el-input 
                  v-model="passwordForm.confirm_password" 
                  type="password" 
                  placeholder="请再次输入新密码"
                  show-password
                  @keyup.enter="handlePasswordChange"
                />
              </div>
            </div>
            <el-button 
              type="primary" 
              class="submit-btn"
              :loading="passwordLoading" 
              @click="handlePasswordChange"
            >
              确认修改
            </el-button>
          </el-form>
        </div>

        <div v-else class="tab-panel">
          <div class="panel-header">
            <el-icon :size="40" class="panel-icon"><QuestionFilled /></el-icon>
            <div>
              <h3>安全问题</h3>
              <p>设置用于找回密码的安全问题和答案</p>
            </div>
          </div>
          <el-form :model="securityForm" class="form-wrapper" @submit.prevent>
            <div class="form-group">
              <label class="group-label">安全问题</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><QuestionFilled /></el-icon>
                <el-input 
                  v-model="securityForm.question" 
                  placeholder="例如：您第一所学校的名称？"
                  @keyup.enter="focusAnswer"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="group-label">安全答案</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><Key /></el-icon>
                <el-input 
                  v-model="securityForm.answer" 
                  type="password" 
                  placeholder="请输入安全答案"
                  show-password
                  @keyup.enter="focusConfirmAnswer"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="group-label">确认答案</label>
              <div class="input-box">
                <el-icon class="input-box-icon"><Key /></el-icon>
                <el-input 
                  v-model="securityForm.confirm_answer" 
                  type="password" 
                  placeholder="请再次输入安全答案"
                  show-password
                  @keyup.enter="handleSecurityChange"
                />
              </div>
            </div>
            <el-button 
              type="primary" 
              class="submit-btn"
              :loading="securityLoading" 
              @click="handleSecurityChange"
            >
              确认修改
            </el-button>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { changePassword, updateSecurityQuestion } from '../api/auth'
import { message } from '@/utils/message'
import { Lock, Key, QuestionFilled, UserFilled } from '@element-plus/icons-vue'
import { md5 } from '@/utils/crypto'

const activeTab = ref('password')

const passwordLoading = ref(false)
const passwordForm = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const securityLoading = ref(false)
const securityForm = ref({
  question: '',
  answer: '',
  confirm_answer: ''
})

async function focusNewPassword() {
  await nextTick()
  const inputs = document.querySelectorAll('.input-box')
  if (inputs[1]) {
    const input = inputs[1].querySelector('input')
    input?.focus()
  }
}

async function focusConfirmPassword() {
  await nextTick()
  const inputs = document.querySelectorAll('.input-box')
  if (inputs[2]) {
    const input = inputs[2].querySelector('input')
    input?.focus()
  }
}

async function focusAnswer() {
  await nextTick()
  const inputs = document.querySelectorAll('.input-box')
  if (inputs[1]) {
    const input = inputs[1].querySelector('input')
    input?.focus()
  }
}

async function focusConfirmAnswer() {
  await nextTick()
  const inputs = document.querySelectorAll('.input-box')
  if (inputs[2]) {
    const input = inputs[2].querySelector('input')
    input?.focus()
  }
}

async function handlePasswordChange() {
  if (!passwordForm.value.current_password || !passwordForm.value.new_password || !passwordForm.value.confirm_password) {
    message.warning('请填写完整密码信息')
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    message.warning('新密码与确认密码不一致')
    return
  }
  if (passwordForm.value.new_password.length < 6) {
    message.warning('新密码至少6位')
    return
  }
  passwordLoading.value = true
  try {
    await changePassword({
      current_password: md5(passwordForm.value.current_password),
      new_password: md5(passwordForm.value.new_password)
    })
    message.success('密码修改成功')
    passwordForm.value = { current_password: '', new_password: '', confirm_password: '' }
  } catch (e) {
    message.error(e.response?.data?.error || '密码修改失败')
  } finally {
    passwordLoading.value = false
  }
}

async function handleSecurityChange() {
  if (!securityForm.value.question || !securityForm.value.answer || !securityForm.value.confirm_answer) {
    message.warning('请填写完整安全问题信息')
    return
  }
  if (securityForm.value.answer !== securityForm.value.confirm_answer) {
    message.warning('两次输入的安全答案不一致')
    return
  }
  securityLoading.value = true
  try {
    await updateSecurityQuestion({
      security_question: securityForm.value.question,
      security_answer: md5(securityForm.value.answer)
    })
    message.success('安全问题修改成功')
    securityForm.value = { question: '', answer: '', confirm_answer: '' }
  } catch (e) {
    message.error(e.response?.data?.error || '安全问题修改失败')
  } finally {
    securityLoading.value = false
  }
}
</script>

<style scoped>
.profile-container {
  min-height: calc(100vh - 64px);
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
}

.profile-layout {
  display: flex;
  width: 100%;
  min-height: calc(100vh - 64px);
  background: white;
}

.profile-sidebar {
  width: 200px;
  flex-shrink: 0;
  background: #fafbfc;
  border-right: 1px solid #f0f0f0;
  padding: 24px 0;
}

.sidebar-header {
  text-align: center;
  padding: 0 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin: 0 auto 12px;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25);
}

.sidebar-header h2 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.menu-list {
  padding: 16px 12px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.menu-item:hover {
  background: #f3f4f6;
  color: #374151;
}

.menu-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.menu-item .el-icon {
  font-size: 18px;
}

.profile-content {
  flex: 1;
  padding: 32px 36px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.tab-panel {
  max-width: 440px;
  width: 100%;
  margin: 0 auto;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
}

.panel-icon {
  color: #667eea;
  flex-shrink: 0;
}

.panel-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px;
}

.panel-header p {
  font-size: 13px;
  color: #9ca3af;
  margin: 0;
}

.form-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.group-label {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.input-box {
  position: relative;
  display: flex;
  align-items: center;
}

.input-box-icon {
  position: absolute;
  left: 14px;
  font-size: 18px;
  color: #9ca3af;
  z-index: 1;
}

.input-box :deep(.el-input__wrapper) {
  padding-left: 42px !important;
  border-radius: 10px;
  box-shadow: 0 0 0 1px #e5e7eb;
}

.input-box :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #667eea;
}

.submit-btn {
  height: 44px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .profile-container {
    padding: 0;
  }

  .profile-layout {
    flex-direction: column;
    border-radius: 0;
    min-height: calc(100vh - 56px);
    max-width: 100%;
  }

  .profile-sidebar {
    width: 100%;
    flex-shrink: 0;
    border-right: none;
    border-bottom: 1px solid #f0f0f0;
    padding: 12px 0 10px;
  }

  .sidebar-header {
    padding: 0 16px 10px;
  }

  .avatar {
    width: 40px;
    height: 40px;
    margin-bottom: 8px;
  }

  .sidebar-header h2 {
    font-size: 15px;
  }

  .menu-list {
    display: flex;
    padding: 0 10px;
    gap: 6px;
  }

  .menu-item {
    flex: 1;
    justify-content: center;
    padding: 8px 6px;
    font-size: 13px;
    gap: 4px;
  }

  .profile-content {
    padding: 16px 16px 24px;
    display: block;
  }

  .tab-panel {
    max-width: none;
  }

  .panel-header {
    margin-bottom: 16px;
  }

  .panel-icon {
    display: none;
  }

  .panel-header h3 {
    font-size: 18px;
  }

  .form-wrapper {
    gap: 14px;
  }
}
</style>
