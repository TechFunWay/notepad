<template>
  <div class="notes-container">
    <div
      v-if="isMobile && sidebarOpen"
      class="sidebar-overlay"
      @click="sidebarOpen = false"
    ></div>

    <div class="notes-sidebar" :class="{ 'sidebar-open': isMobile && sidebarOpen }">
      <div class="sidebar-header">
        <el-button type="primary" class="create-btn" :icon="Plus" @click="createNewNote">
          <span>新建笔记</span>
        </el-button>
      </div>
      <div class="sidebar-search">
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
      </div>
      <div v-if="allTags.length > 0" class="sidebar-tags">
        <div class="tags-title">
          <el-icon><Tickets /></el-icon>
          <span>标签</span>
        </div>
        <div class="tags-list">
          <el-tag
            v-for="tag in allTags"
            :key="tag"
            class="tag-item"
            :class="{ active: activeTag === tag }"
            @click="toggleTag(tag)"
          >{{ tag }}</el-tag>
          <el-tag v-if="activeTag" class="tag-item clear-tag" @click="activeTag = ''; loadNotes()">
            <el-icon><Close /></el-icon>
            <span>清除筛选</span>
          </el-tag>
        </div>
      </div>
      <div class="note-list">
        <div
          v-for="note in notes"
          :key="note.id"
          class="note-item"
          :class="{ active: currentNote?.id === note.id }"
          @click="selectNote(note)"
        >
          <div class="note-item-header">
            <div class="note-item-title">{{ note.title || '无标题' }}</div>
            <div class="note-item-dot"></div>
          </div>
          <div class="note-item-preview">{{ stripHtml(note.content) }}</div>
          <div class="note-item-footer">
            <span class="note-item-time">
              <el-icon><Clock /></el-icon>
              {{ formatDate(note.updated_at) }}
            </span>
            <div v-if="note.tags" class="note-item-tags">
              <span v-for="t in splitTags(note.tags).slice(0, 2)" :key="t" class="mini-tag">{{ t }}</span>
            </div>
          </div>
        </div>
        <div v-if="notes.length === 0 && !loading" class="empty-list">
          <div class="empty-icon">
            <el-icon :size="56"><Document /></el-icon>
          </div>
          <h3>{{ searchQuery ? '没有找到笔记' : '还没有笔记' }}</h3>
          <p>{{ searchQuery ? '试试其他关键词' : '点击上方按钮创建新笔记' }}</p>
        </div>
      </div>
    </div>

    <div class="notes-editor">
      <template v-if="currentNote">
        <div v-if="isMobile" class="editor-mobile-top">
          <el-button class="back-btn" :icon="ArrowLeft" text @click="backToList">
            <span>返回</span>
          </el-button>
        </div>

        <div class="editor-header">
          <el-button
            v-if="!isMobile"
            type="primary"
            class="back-to-list-btn"
            :icon="Tickets"
            text
            @click="backToList"
          >
            <span>笔记列表</span>
          </el-button>
          <el-input
            v-model="currentNote.title"
            placeholder="笔记标题"
            class="title-input"
            @input="markDirty"
          />
        </div>

        <div class="editor-toolbar">
          <div class="toolbar-left">
            <div class="meta-item">
              <el-icon><Calendar /></el-icon>
              <span>创建: {{ formatFullDate(currentNote.created_at) }}</span>
            </div>
            <div class="meta-item">
              <el-icon><Edit /></el-icon>
              <span>修改: {{ formatFullDate(currentNote.updated_at) }}</span>
            </div>
          </div>
          <div class="toolbar-right">
            <div class="tags-selector">
              <el-icon><Tickets /></el-icon>
              <el-select
                v-model="currentTagInput"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="添加标签..."
                size="small"
                class="tags-select"
                @change="markDirty"
                @remove-tag="markDirty"
              >
                <el-option v-for="tag in allTags" :key="tag" :label="tag" :value="tag" />
              </el-select>
            </div>
            <div class="action-buttons">
              <el-button 
                type="primary" 
                class="save-btn"
                :icon="Check" 
                :loading="saving" 
                @click="saveNote"
              >
                <span>保存</span>
              </el-button>
              <el-popconfirm title="确定删除这个笔记吗？" @confirm="removeNote">
                <template #reference>
                  <el-button class="delete-btn" :icon="Delete"></el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>

        <div class="editor-content">
          <TiptapEditor v-model="currentNote.content" @update:model-value="markDirty" />
        </div>
      </template>

      <div v-else class="editor-empty">
        <div class="dashboard-container">
          <div class="dashboard-header">
            <div class="empty-animation">
              <div class="empty-icon-large">
                <el-icon :size="isMobile ? 80 : 96"><EditPen /></el-icon>
              </div>
            </div>
            <h2>欢迎回来</h2>
            <p>开始记录您的灵感与想法</p>
          </div>
          
          <div class="dashboard-stats">
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ notes.length }}</div>
                <div class="stat-label">总笔记数</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><Tickets /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ allTags.length }}</div>
                <div class="stat-label">标签数</div>
              </div>
            </div>
          </div>
          
          <div class="recent-notes-section">
            <div class="section-header">
              <h3>最近笔记</h3>
              <el-button text class="more-link" @click="loadAllNotes">
                <span>查看全部</span>
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
            <div class="recent-notes-list">
              <div v-for="note in recentNotes" :key="note.id" class="recent-note-item" @click="selectNote(note)">
                <div class="recent-note-title">{{ note.title || '无标题' }}</div>
                <div class="recent-note-preview">{{ stripHtml(note.content) }}</div>
                <div class="recent-note-time">{{ formatDate(note.updated_at) }}</div>
              </div>
              <div v-if="recentNotes.length === 0" class="no-recent-notes">
                <p>还没有笔记，点击上方按钮创建第一个笔记</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Plus, Check, Delete, Document, EditPen, ArrowLeft, FolderOpened, Notebook, Clock, Search, Tickets, Close, Calendar, Edit } from '@element-plus/icons-vue'
import { getNotes, createNote, updateNote, deleteNote, getNote } from '@/api/note'
import api from '@/api/request'
import { ElMessage } from 'element-plus'
import TiptapEditor from '@/components/TiptapEditor.vue'

const router = useRouter()
const route = useRoute()

const notes = ref([])
const currentNote = ref(null)
const currentTagInput = ref([])
const searchQuery = ref('')
const loading = ref(false)
const saving = ref(false)
const isDirty = ref(false)
const isMobile = ref(false)
const sidebarOpen = ref(false)
const allTags = ref([])
const activeTag = ref('')
let searchTimer = null

// 最近笔记（最多显示5个）
const recentNotes = ref([])

// 加载所有笔记
function loadAllNotes() {
  router.push('/notes-list')
}

onMounted(async () => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  await loadNotes()
  await loadTags()
  await openNoteFromQuery()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

async function loadNotes() {
  loading.value = true
  try {
    const params = {}
    if (searchQuery.value) params.q = searchQuery.value
    if (activeTag.value) params.tag = activeTag.value
    const { data } = await getNotes(params)
    notes.value = data.notes
    
    // 更新最近笔记（按更新时间排序，最多显示5个）
    recentNotes.value = [...notes.value]
      .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
      .slice(0, 5)
  } catch (e) {
    ElMessage.error('加载笔记失败')
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

async function openNoteFromQuery() {
  const noteId = route.query.note_id
  if (!noteId) return
  try {
    const { data } = await getNote(noteId)
    currentNote.value = { ...data }
    currentTagInput.value = data.tags ? data.tags.split(',').map(t => t.trim()).filter(Boolean) : []
    isDirty.value = false
  } catch (e) {
    ElMessage.error('加载笔记失败')
  }
}

function toggleTag(tag) {
  activeTag.value = activeTag.value === tag ? '' : tag
  loadNotes()
}

function debouncedSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => loadNotes(), 300)
}

async function createNewNote() {
  try {
    const today = new Date()
    const yy = String(today.getFullYear()).slice(2)
    const mm = String(today.getMonth() + 1).padStart(2, '0')
    const dd = String(today.getDate()).padStart(2, '0')
    const title = `${yy}${mm}${dd}的笔记`
    const { data } = await createNote({ title, content: '<p></p>', tags: '' })
    notes.value.unshift(data)
    currentNote.value = { ...data }
    currentTagInput.value = []
    isDirty.value = false
    sidebarOpen.value = false
  } catch (e) {
    ElMessage.error('创建笔记失败')
  }
}

function selectNote(note) {
  currentNote.value = { ...note }
  currentTagInput.value = note.tags ? note.tags.split(',').map(t => t.trim()).filter(Boolean) : []
  isDirty.value = false
  if (isMobile.value) {
    sidebarOpen.value = false
  }
}

function backToList() {
  currentNote.value = null
  sidebarOpen.value = true
}

function markDirty() {
  isDirty.value = true
}

async function saveNote() {
  if (!currentNote.value) return
  saving.value = true
  try {
    const tags = currentTagInput.value.join(', ')
    await updateNote(currentNote.value.id, {
      title: currentNote.value.title,
      content: currentNote.value.content,
      tags
    })
    isDirty.value = false
    ElMessage.success('保存成功')
    await loadNotes()
    await loadTags()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function removeNote() {
  if (!currentNote.value) return
  try {
    await deleteNote(currentNote.value.id)
    currentNote.value = null
    ElMessage.success('删除成功')
    await loadNotes()
    await loadTags()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

function stripHtml(html) {
  if (!html) return ''
  const div = document.createElement('div')
  div.innerHTML = html
  const text = div.textContent || div.innerText || ''
  return text.length > 80 ? text.substring(0, 80) + '...' : text
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

function formatFullDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.notes-container {
  display: flex;
  height: calc(100vh - 64px);
  position: relative;
  background: #f5f7fa;
}

.notes-sidebar {
  width: 320px;
  background: white;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  border-right: 1px solid #e8eaed;
}

.sidebar-header {
  padding: 20px;
}

.create-btn {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s;
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.sidebar-search {
  padding: 0 20px 16px;
}

.search-box {
  position: relative;
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

.sidebar-tags {
  padding: 0 20px 16px;
}

.tags-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 10px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  background: #f3f4f6;
  color: #4b5563;
  border: none;
  transition: all 0.2s;
}

.tag-item:hover {
  background: #e5e7eb;
  transform: translateY(-1px);
}

.tag-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.clear-tag {
  background: #fef2f2;
  color: #dc2626;
}

.clear-tag:hover {
  background: #fee2e2;
}

.note-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 12px 12px;
}

.note-item {
  padding: 16px;
  margin-bottom: 8px;
  cursor: pointer;
  border-radius: 12px;
  transition: all 0.2s;
  background: #fafafa;
}

.note-item:hover {
  background: #f3f4f6;
  transform: translateY(-1px);
}

.note-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border: 2px solid #667eea;
}

.note-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.note-item-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.note-item-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #667eea;
  opacity: 0;
  flex-shrink: 0;
}

.note-item.active .note-item-dot {
  opacity: 1;
}

.note-item-preview {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 10px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.note-item-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.note-item-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.note-item-tags {
  display: flex;
  gap: 4px;
}

.mini-tag {
  padding: 2px 8px;
  background: #e5e7eb;
  color: #4b5563;
  font-size: 11px;
  border-radius: 6px;
}

.empty-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  color: #d1d5db;
  margin-bottom: 16px;
}

.empty-list h3 {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 6px;
}

.empty-list p {
  font-size: 13px;
  color: #9ca3af;
  margin: 0;
}

.notes-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  min-width: 0;
  overflow: hidden;
}

.editor-mobile-top {
  padding: 12px 16px 0;
  flex-shrink: 0;
  border-bottom: 1px solid #f3f4f6;
}

.back-btn {
  color: #667eea;
  font-weight: 500;
}

.editor-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 32px 8px;
  flex-shrink: 0;
}

.back-to-list-btn {
  flex-shrink: 0;
  color: #667eea;
  font-weight: 500;
  font-size: 14px;
}

.back-to-list-btn:hover {
  color: #764ba2;
}

.title-input :deep(.el-input__wrapper) {
  box-shadow: none !important;
  font-size: 28px;
  font-weight: 700;
  padding: 0;
}

.title-input :deep(.el-input__inner) {
  color: #1f2937;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 32px 16px;
  border-bottom: 1px solid #f3f4f6;
  flex-shrink: 0;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-left {
  display: flex;
  gap: 20px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tags-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
}

.tags-select {
  min-width: 180px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.save-btn {
  border-radius: 10px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.delete-btn {
  border-radius: 10px;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
}

.editor-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  color: #6b7280;
  text-align: center;
  padding: 16px;
}

.empty-animation {
  margin-bottom: 16px;
}

.empty-icon-large {
  color: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  opacity: 0.6;
}

.editor-empty h2 {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px;
}

.editor-empty p {
  font-size: 15px;
  color: #6b7280;
  margin: 0 0 16px;
}

.empty-btn {
  border-radius: 12px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

/* 仪表板样式 */
.dashboard-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 4px 0;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 8px;
}

.dashboard-header h2 {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin: 4px 0 2px;
  letter-spacing: -0.5px;
}

.dashboard-header p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 8px;
}

.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 8px;
  margin-bottom: 16px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 20px 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  border: 1px solid #f3f4f6;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-icon .el-icon {
  font-size: 24px;
}

.stat-content {
  flex: 1;
  text-align: left;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 2px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.recent-notes-section {
  background: white;
  border-radius: 16px;
  padding: 20px 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid #f3f4f6;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.more-link {
  color: #667eea !important;
  font-weight: 500;
  font-size: 12px;
}

.more-link:hover {
  color: #764ba2 !important;
}

.recent-notes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recent-note-item {
  padding: 16px;
  border-radius: 12px;
  background: #f9fafb;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  border: 2px solid transparent;
}

.recent-note-item:hover {
  background: #f3f4f6;
  transform: translateY(-1px);
  border-color: #e5e7eb;
}

.recent-note-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recent-note-preview {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 8px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.recent-note-time {
  font-size: 11px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  gap: 4px;
}

.no-recent-notes {
  text-align: center;
  padding: 32px 16px;
  color: #9ca3af;
  font-size: 13px;
}

.empty-btn {
  margin-top: 16px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

/* 移动端特定优化 */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 4px 4px;
  }
  
  .dashboard-header h2 {
    font-size: 16px;
    margin: 2px 0 1px;
  }
  
  .dashboard-header p {
    font-size: 12px;
    margin: 0 0 8px;
  }
  
  .dashboard-stats {
    grid-template-columns: repeat(2, 1fr);
    gap: 6px;
  }
  
  .stat-card {
    padding: 10px 8px;
  }
  
  .stat-icon {
    width: 32px;
    height: 32px;
  }
  
  .stat-icon .el-icon {
    font-size: 16px;
  }
  
  .stat-number {
    font-size: 16px;
  }
  
  .stat-label {
    font-size: 10px;
  }
  
  .recent-notes-section {
    padding: 10px 8px;
  }
  
  .section-header h3 {
    font-size: 14px;
  }
  
  .recent-note-item {
    padding: 8px;
  }
  
  .recent-note-title {
    font-size: 12px;
  }
  
  .recent-note-preview {
    font-size: 11px;
  }
  
  .empty-btn {
    margin-top: 8px;
    padding: 6px 12px;
    font-size: 12px;
  }
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
  z-index: 100;
}

@media (max-width: 768px) {
  .notes-container {
    height: calc(100vh - 56px);
  }

  .notes-sidebar {
    display: none;
  }

  .sidebar-overlay {
    display: none;
  }

  .notes-editor {
    width: 100%;
  }

  .editor-header {
    padding: 16px 16px 8px;
  }

  .title-input :deep(.el-input__wrapper) {
    font-size: 24px;
  }

  .editor-toolbar {
    padding: 12px 16px;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .toolbar-left {
    flex-direction: column;
    gap: 6px;
  }

  .toolbar-right {
    width: 100%;
    justify-content: space-between;
    padding-right: 4px;
  }

  .action-buttons {
    flex-shrink: 0;
  }

  .tags-select {
    flex: 1;
  }

  .editor-content {
    padding-bottom: 20px;
  }
}
</style>
