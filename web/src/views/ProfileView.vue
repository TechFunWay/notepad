<template>
  <div class="profile-container">
    <div class="profile-wrapper">
      <div class="profile-header">
        <div class="header-icon">
          <el-icon :size="48"><Lock /></el-icon>
        </div>
        <h1>个人设置</h1>
        <p>管理您的账号安全</p>
      </div>

      <div class="profile-card">
        <div class="card-title">
          <el-icon><Key /></el-icon>
          <span>修改密码</span>
        </div>
        <el-form :model="form" class="profile-form">
          <el-form-item>
            <div class="form-item-wrapper">
              <div class="form-label">当前密码</div>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.current_password" 
                  type="password" 
                  placeholder="请输入当前密码"
                  show-password 
                  class="custom-input"
                  size="large"
                />
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <div class="form-item-wrapper">
              <div class="form-label">新密码</div>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.new_password" 
                  type="password" 
                  placeholder="请输入新密码（至少6位）"
                  show-password 
                  class="custom-input"
                  size="large"
                />
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <div class="form-item-wrapper">
              <div class="form-label">确认密码</div>
              <div class="input-wrapper">
                <el-icon class="input-icon"><Lock /></el-icon>
                <el-input 
                  v-model="form.confirm_password" 
                  type="password" 
                  placeholder="请再次输入新密码"
                  show-password 
                  class="custom-input"
                  size="large"
                />
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn"
              :loading="loading" 
              @click="handleChange"
            >
              <span>确认修改</span>
              <el-icon><Check /></el-icon>
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { changePassword } from '../api/auth'
import { ElMessage } from 'element-plus'
import { Lock, Key, Check } from '@element-plus/icons-vue'
import { md5 } from '@/utils/crypto'

const loading = ref(false)
const form = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

async function handleChange() {
  if (!form.value.current_password || !form.value.new_password) {
    ElMessage.warning('请填写所有字段')
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  if (form.value.new_password.length < 6) {
    ElMessage.warning('新密码至少6位')
    return
  }
  loading.value = true
  try {
    await changePassword({
      current_password: md5(form.value.current_password),
      new_password: md5(form.value.new_password)
    })
    ElMessage.success('密码修改成功')
    form.value = { current_password: '', new_password: '', confirm_password: '' }
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '修改失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.profile-container {
  min-height: calc(100vh - 64px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
}

.profile-wrapper {
  width: 100%;
  max-width: 520px;
  transform: translateY(-10px);
}

.profile-header {
  text-align: center;
  margin-bottom: 28px;
}

.header-icon {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin: 0 auto 16px;
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.2);
}

.profile-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 6px;
  letter-spacing: -0.5px;
}

.profile-header p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.profile-card {
  background: white;
  border-radius: 24px;
  padding: 36px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.profile-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 28px;
}

.card-title .el-icon {
  color: #667eea;
}

.profile-form {
  margin-bottom: 0;
}

.form-item-wrapper {
  width: 100%;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
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
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s;
  margin-top: 8px;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.submit-btn:active {
  transform: translateY(0);
}

@media (max-width: 768px) {
  .profile-container {
    min-height: calc(100vh - 56px);
    padding: 24px 16px;
    align-items: flex-start;
  }

  .header-icon {
    width: 64px;
    height: 64px;
    border-radius: 16px;
  }

  .header-icon .el-icon {
    font-size: 32px;
  }

  .profile-header h1 {
    font-size: 24px;
  }

  .profile-card {
    padding: 28px 20px;
    border-radius: 16px;
  }
}
</style>
