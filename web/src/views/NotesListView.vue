<template>
  <div class="notes-list-container">
    <div class="notes-list-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon><Document /></el-icon>
          <span>我的笔记</span>
        </h1>
        <span class="notes-count">共 {{ total }} 条笔记</span>
      </div>
      <div class="header-right">
        <el-button type="primary" class="create-btn" :icon="Plus" @click="createNewNote">
          <span>新建笔记</span>
        </el-button>
      </div>
    </div>

    <div class="notes-list-toolbar">
      <div class="search-box">
        <el-icon class="search-icon"><Search /></el-icon>
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索笔记..." 
          clearable 
          @input="debouncedSearch"
          class="search-input"
        />
      </div>
      
      <div class="filter-group">
        <div v-if="allTags.length > 0" class="tag-filter">
          <el-select
            v-model="activeTag"
            placeholder="按标签筛选"
            clearable
            @change="loadNotes"
            class="tag-select"
          >
            <el-option v-for="tag in allTags" :key="tag" :label="tag" :value="tag" />
          </el-select>
        </div>
        
        <div class="sort-filter">
          <el-select
            v-model="sortBy"
            placeholder="排序方式"
            @change="loadNotes"
            class="sort-select"
          >
            <el-option label="最近修改" value="updated_at" />
            <el-option label="最近创建" value="created_at" />
            <el-option label="标题A-Z" value="title" />
          </el-select>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="6" animated />
    </div>

    <div v-else-if="notes.length === 0" class="empty-container">
      <div class="empty-content">
        <el-icon :size="80"><Document /></el-icon>
        <h3>{{ searchQuery || activeTag ? '没有找到笔记' : '还没有笔记' }}</h3>
        <p>{{ searchQuery || activeTag ? '试试其他关键词或筛选条件' : '点击上方按钮创建新笔记' }}</p>
        <el-button v-if="!searchQuery && !activeTag" type="primary" :icon="Plus" @click="createNewNote">
          <span>创建第一个笔记</span>
        </el-button>
      </div>
    </div>

    <div v-else class="notes-grid">
      <div 
        v-for="note in notes" 
        :key="note.id" 
        class="note-card"
        @click="openNote(note)"
      >
        <div class="note-card-header">
          <h3 class="note-title">{{ note.title || '无标题' }}</h3>
          <div class="note-actions">
            <el-button 
              text 
              :icon="Edit" 
              @click.stop="openNote(note)"
            />
            <el-popconfirm title="确定删除这个笔记吗？" @confirm.stop="deleteNoteItem(note)">
              <template #reference>
                <el-button text :icon="Delete" />
              </template>
            </el-popconfirm>
          </div>
        </div>
        
        <div class="note-preview">{{ stripHtml(note.content) }}</div>
        
        <div class="note-footer">
          <div class="note-meta">
            <span class="note-time">
              <el-icon><Clock /></el-icon>
              {{ formatDate(note.updated_at) }}
            </span>
          </div>
          <div v-if="note.tags" class="note-tags">
            <el-tag v-for="t in splitTags(note.tags).slice(0, 3)" :key="t" size="small">{{ t }}</el-tag>
          </div>
        </div>
      </div>
    </div>

    <div v-if="notes.length > 0 && total > pageSize" class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadNotes"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Document, Search, Edit, Delete, Clock } from '@element-plus/icons-vue'
import { getNotes, createNote, deleteNote } from '@/api/note'
import api from '@/api/request'
import { message } from '@/utils/message'

const router = useRouter()

const notes = ref([])
const loading = ref(false)
const searchQuery = ref('')
const activeTag = ref('')
const sortBy = ref('updated_at')
const currentPage = ref(1)
const pageSize = ref(24)
const total = ref(0)
const allTags = ref([])

let searchTimer = null

onMounted(() => {
  loadNotes()
  loadTags()
})

async function loadNotes() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchQuery.value) params.q = searchQuery.value
    if (activeTag.value) params.tag = activeTag.value
    params.sort_by = sortBy.value

    const { data } = await getNotes(params)
    notes.value = data.notes
    total.value = data.total || notes.value.length
  } catch (e) {
    message.error('加载笔记失败')
  } finally {
    loading.value = false
  }
}

async function loadTags() {
  try {
    const { data } = await api.get('/notes/tags')
    allTags.value = data.tags || []
  } catch (e) {}
}

function debouncedSearch() {
  clearTimeout(searchTimer)
  currentPage.value = 1
  searchTimer = setTimeout(() => loadNotes(), 300)
}

function handleSizeChange() {
  currentPage.value = 1
  loadNotes()
}

async function createNewNote() {
  try {
    const today = new Date()
    const yy = String(today.getFullYear()).slice(2)
    const mm = String(today.getMonth() + 1).padStart(2, '0')
    const dd = String(today.getDate()).padStart(2, '0')
    const title = `${yy}${mm}${dd}的笔记`
    const { data } = await createNote({ title, content: '<p></p>', tags: '' })
    message.success('创建成功')
    router.push({ path: '/', query: { note_id: data.id } })
  } catch (e) {
    message.error('创建笔记失败')
  }
}

function openNote(note) {
  router.push('/')
}

async function deleteNoteItem(note) {
  try {
    await deleteNote(note.id)
    message.success('删除成功')
    await loadNotes()
    await loadTags()
  } catch (e) {
    message.error('删除失败')
  }
}

function stripHtml(html) {
  if (!html) return ''
  const div = document.createElement('div')
  div.innerHTML = html
  const text = div.textContent || div.innerText || ''
  return text.length > 100 ? text.substring(0, 100) + '...' : text
}

function splitTags(tagsStr) {
  if (!tagsStr) return []
  return tagsStr.split(',').map(t => t.trim()).filter(Boolean)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'
  return d.toLocaleDateString()
}
</script>

<style scoped>
.notes-list-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.notes-list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.page-title .el-icon {
  color: #667eea;
}

.notes-count {
  font-size: 14px;
  color: #6b7280;
  padding: 4px 12px;
  background: #f3f4f6;
  border-radius: 12px;
}

.create-btn {
  border-radius: 12px;
  font-weight: 600;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.notes-list-toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 280px;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 18px;
  z-index: 1;
}

.search-input :deep(.el-input__wrapper) {
  padding-left: 44px !important;
  border-radius: 12px;
  box-shadow: 0 0 0 1px #e5e7eb;
  background: #f9fafb;
  transition: all 0.3s;
}

.search-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #667eea;
  background: white;
}

.search-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #667eea;
  background: white;
}

.filter-group {
  display: flex;
  gap: 12px;
}

.tag-select,
.sort-select {
  width: 160px;
}

.loading-container {
  padding: 40px;
  background: white;
  border-radius: 16px;
}

.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.empty-content {
  text-align: center;
  padding: 40px;
}

.empty-content .el-icon {
  color: #d1d5db;
  margin-bottom: 24px;
}

.empty-content h3 {
  font-size: 20px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px;
}

.empty-content p {
  font-size: 14px;
  color: #9ca3af;
  margin: 0 0 24px;
}

.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.note-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.3s;
}

.note-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #e5e7eb;
}

.note-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 12px;
  gap: 12px;
}

.note-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.note-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.note-card:hover .note-actions {
  opacity: 1;
}

.note-preview {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.6;
  margin-bottom: 16px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.note-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.note-meta {
  display: flex;
  align-items: center;
}

.note-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.note-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 24px 0;
}

@media (max-width: 768px) {
  .notes-list-container {
    padding: 12px;
  }

  .notes-list-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-left {
    justify-content: space-between;
  }

  .page-title {
    font-size: 20px;
  }

  .create-btn {
    width: 100%;
    justify-content: center;
  }

  .notes-list-toolbar {
    flex-direction: column;
    gap: 12px;
  }

  .search-box {
    width: 100%;
    min-width: auto;
  }

  .filter-group {
    width: 100%;
    gap: 8px;
  }

  .tag-select,
  .sort-select {
    flex: 1;
    width: auto;
  }

  .notes-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .note-card {
    padding: 16px;
  }

  .note-title {
    font-size: 14px;
  }

  .note-preview {
    font-size: 13px;
    -webkit-line-clamp: 2;
  }

  .pagination-container {
    margin-top: 16px;
    padding: 16px 0;
  }

  :deep(.el-pagination) {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
