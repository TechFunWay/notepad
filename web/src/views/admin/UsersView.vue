<template>
  <div class="admin-container">
    <div class="admin-header">
      <h2>用户管理</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog">新建用户</el-button>
    </div>

    <!-- Desktop table -->
    <el-table v-if="!isMobile" :data="users" stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="role" label="角色" width="120">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
            {{ row.role === 'admin' ? '管理员' : '用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="showEditDialog(row)">编辑</el-button>
          <el-popconfirm title="确定删除该用户？" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- Mobile card list -->
    <div v-else class="mobile-card-list">
      <div v-for="user in users" :key="user.id" class="mobile-card">
        <div class="mobile-card-header">
          <span class="mobile-card-title">{{ user.username }}</span>
          <el-tag :type="user.role === 'admin' ? 'danger' : 'info'" size="small">
            {{ user.role === 'admin' ? '管理员' : '用户' }}
          </el-tag>
        </div>
        <div class="mobile-card-info">
          <span>ID: {{ user.id }}</span>
          <span>{{ new Date(user.created_at).toLocaleDateString() }}</span>
        </div>
        <div class="mobile-card-actions">
          <el-button size="small" @click="showEditDialog(user)">编辑</el-button>
          <el-popconfirm title="确定删除该用户？" @confirm="handleDelete(user.id)">
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
      <div v-if="users.length === 0" class="empty-list">
        <p>暂无用户</p>
      </div>
    </div>

    <!-- Create Dialog -->
    <el-dialog v-model="createDialogVisible" title="新建用户" :width="isMobile ? '95%' : '400px'">
      <el-form :model="createForm" :label-width="isMobile ? 'auto' : '80px'">
        <el-form-item label="用户名">
          <el-input v-model="createForm.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="createForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="createForm.role">
            <el-option label="用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- Edit Dialog -->
    <el-dialog v-model="editDialogVisible" title="编辑用户" :width="isMobile ? '95%' : '400px'">
      <el-form :model="editForm" :label-width="isMobile ? 'auto' : '80px'">
        <el-form-item label="用户名">
          <el-input :model-value="editForm.username" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="editForm.role" :disabled="editForm.isSelf">
            <el-option label="用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="editForm.password" type="password" show-password placeholder="留空不修改" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { getUsers, createUser, updateUser, deleteUser } from '../../api/user'
import { message } from '@/utils/message'
import { md5 } from '@/utils/crypto'

const users = ref([])
const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const isMobile = ref(false)
const createForm = ref({ username: '', password: '', role: 'user' })
const editForm = ref({ id: 0, username: '', role: 'user', password: '', isSelf: false })

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadUsers()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

async function loadUsers() {
  const { data } = await getUsers()
  users.value = data.users
}

function showCreateDialog() {
  createForm.value = { username: '', password: '', role: 'user' }
  createDialogVisible.value = true
}

async function handleCreate() {
  try {
    await createUser({
      username: createForm.value.username,
      password: md5(createForm.value.password),
      role: createForm.value.role
    })
    message.success('创建成功')
    createDialogVisible.value = false
    loadUsers()
  } catch (e) {
    message.error(e.response?.data?.error || '创建失败')
  }
}

function showEditDialog(user) {
  const currentUserId = parseInt(localStorage.getItem('userId') || '0')
  editForm.value = { 
    id: user.id, 
    username: user.username, 
    role: user.role, 
    password: '',
    isSelf: user.id === currentUserId
  }
  editDialogVisible.value = true
}

async function handleEdit() {
  try {
    const payload = { role: editForm.value.role }
    if (editForm.value.password) payload.password = md5(editForm.value.password)
    await updateUser(editForm.value.id, payload)
    message.success('更新成功')
    editDialogVisible.value = false
    loadUsers()
  } catch (e) {
    message.error(e.response?.data?.error || '更新失败')
  }
}

async function handleDelete(id) {
  try {
    await deleteUser(id)
    message.success('删除成功')
    loadUsers()
  } catch (e) {
    message.error(e.response?.data?.error || '删除失败')
  }
}
</script>

<style scoped>
.admin-container {
  padding: 24px;
  overflow-y: auto;
  height: 100%;
  -webkit-overflow-scrolling: touch;
  background: var(--bg-primary);
  border-radius: 16px;
  box-shadow: var(--card-shadow);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.admin-header h2 {
  margin: 0;
}

.mobile-card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mobile-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 14px 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.mobile-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.mobile-card-title {
  font-weight: 500;
  font-size: 15px;
}

.mobile-card-info {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.mobile-card-actions {
  display: flex;
  gap: 8px;
}

.empty-list {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .admin-container {
    padding: 16px;
  }

  .admin-header h2 {
    font-size: 18px;
  }
}
</style>
