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
          role="button"
          tabindex="0"
          :aria-label="`打开笔记：${note.title || '无标题'}`"
          @click="selectNote(note)"
          @keyup.enter="selectNote(note)"
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

        <!-- View Mode -->
        <div v-if="!editMode" class="note-view">
          <div class="view-header">
            <div class="view-title-row">
              <el-button
                v-if="!isMobile"
                class="back-to-list-btn"
                :icon="Tickets"
                text
                @click="backToList"
              >
                <span>笔记列表</span>
              </el-button>
              <h1 class="view-title">{{ currentNote.title || '无标题' }}</h1>
              <el-button
                v-if="!isMobile"
                class="view-edit-btn"
                type="primary"
                :icon="EditPen"
                @click="enterEditMode"
              >
                <span>编辑</span>
              </el-button>
            </div>
            <div class="view-meta">
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                创建: {{ formatFullDate(currentNote.created_at) }}
              </span>
              <span class="meta-item">
                <el-icon><Edit /></el-icon>
                修改: {{ formatFullDate(currentNote.updated_at) }}
              </span>
            </div>
            <div v-if="currentTagInput.length" class="view-tags">
              <button
                v-for="t in currentTagInput"
                :key="t"
                type="button"
                class="view-tag"
                @click="openTagNotes(t)"
              >
                # {{ t }}
              </button>
            </div>
          </div>
          <div class="view-content" v-html="currentNote.content"></div>
          <div class="view-footer">
            <el-button
              class="reader-edit-action"
              type="primary"
              :icon="EditPen"
              @click="enterEditMode"
            >
              {{ isMobile ? '编辑笔记' : '编辑' }}
            </el-button>
            <el-popconfirm title="确定删除这个笔记吗？" @confirm="removeNote">
              <template #reference>
                <el-button
                  class="reader-delete-action"
                  :icon="Delete"
                  aria-label="删除笔记"
                >
                  <span v-if="!isMobile">删除</span>
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>

        <!-- Edit Mode -->
        <div v-if="editMode" class="editor-edit-area">
          <div class="editor-header">
            <el-button
              v-if="!isMobile"
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
            <span class="save-status" :class="{ dirty: isDirty }">
              {{ saving ? '保存中…' : isDirty ? '待保存' : '已保存' }}
            </span>
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
                <div class="tag-entry-shell">
                  <span v-for="tag in currentTagInput" :key="tag" class="editor-tag-chip">
                    {{ tag }}
                    <button
                      type="button"
                      :aria-label="`移除标签：${tag}`"
                      @click="removeCurrentTag(tag)"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </span>
                  <input
                    v-model="tagDraft"
                    class="tag-entry-input"
                    type="text"
                    list="note-tag-suggestions"
                    placeholder="添加标签，按回车确认"
                    aria-label="添加标签"
                    @keydown="handleTagKeydown"
                    @blur="addCurrentTag"
                  />
                  <datalist id="note-tag-suggestions">
                    <option v-for="tag in availableTagSuggestions" :key="tag" :value="tag" />
                  </datalist>
                </div>
              </div>
            </div>
          </div>

          <div class="editor-content">
            <TiptapEditor v-model="currentNote.content" @update:model-value="markDirty" />
          </div>

          <!-- 底部操作栏(桌面 + 移动通用,内容多时无需滚回顶部) -->
          <div class="editor-bottom-bar">
            <span class="bar-word-count">{{ wordCount }} 字</span>
            <div class="bar-actions">
              <el-button class="bar-btn" :icon="View" @click="exitEditMode">预览</el-button>
              <el-button
                type="primary"
                class="bar-btn bar-save-btn"
                :icon="Check"
                :loading="saving"
                @click="saveNote"
              >
                <span>保存</span>
              </el-button>
              <el-popconfirm title="确定删除这个笔记吗？" @confirm="removeNote">
                <template #reference>
                  <el-button class="bar-btn bar-delete-btn" :icon="Delete"></el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
      </template>

      <div v-else class="editor-empty">
        <div class="dashboard-container">
          <div class="dashboard-header">
            <div>
              <span class="dashboard-date">{{ dashboardDate }}</span>
              <h1>{{ greeting }}</h1>
              <p>给今天留一点空间，写下正在发生的事。</p>
            </div>
            <el-button class="dashboard-create" type="primary" :icon="Plus" @click="createNewNote">
              新建笔记
            </el-button>
          </div>

          <div class="dashboard-stats">
            <button
              type="button"
              class="stat-card stat-card-action"
              aria-label="查看全部笔记"
              @click="loadAllNotes"
            >
              <div class="stat-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ notes.length }}</div>
                <div class="stat-label">总笔记数</div>
              </div>
              <el-icon class="stat-arrow"><ArrowRight /></el-icon>
            </button>
            <button
              type="button"
              class="stat-card stat-card-action"
              aria-label="管理标签"
              @click="tagManagerOpen = true"
            >
              <div class="stat-icon">
                <el-icon><Tickets /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ allTags.length }}</div>
                <div class="stat-label">标签数</div>
              </div>
              <el-icon class="stat-arrow"><ArrowRight /></el-icon>
            </button>
          </div>

          <div class="recent-notes-section">
            <div class="section-header">
              <div>
                <span class="section-kicker">RECENT</span>
                <h2>最近笔记</h2>
              </div>
              <el-button text class="more-link" @click="loadAllNotes">
                <span>查看全部</span>
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
            <div class="recent-notes-list">
              <div
                v-for="note in recentNotes"
                :key="note.id"
                class="recent-note-item"
                role="button"
                tabindex="0"
                @click="selectNote(note)"
                @keyup.enter="selectNote(note)"
              >
                <div class="recent-note-head">
                  <span class="recent-note-index">{{ String(note.id).slice(-2).padStart(2, '0') }}</span>
                  <div class="recent-note-title">{{ note.title || '无标题' }}</div>
                </div>
                <div class="recent-note-preview">{{ stripHtml(note.content) }}</div>
                <div class="recent-note-footer">
                  <div class="recent-note-time">{{ formatDate(note.updated_at) }}</div>
                  <el-button
                    class="recent-note-edit-btn"
                    size="small"
                    :icon="EditPen"
                    @click.stop="quickEditNote(note)"
                  >
                    <span>编辑</span>
                  </el-button>
                </div>
              </div>
              <div v-if="recentNotes.length === 0" class="no-recent-notes">
                <p>{{ isMobile ? '还没有笔记，点击右下角开始写作' : '还没有笔记，点击上方按钮创建第一篇' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 移动端浮动新建按钮 -->
    <button v-if="isMobile && !currentNote" class="mobile-fab" type="button" aria-label="新建笔记" @click="createNewNote">
      <el-icon :size="24"><Plus /></el-icon>
    </button>

    <el-drawer
      v-model="tagManagerOpen"
      :direction="isMobile ? 'btt' : 'rtl'"
      :size="isMobile ? '72%' : '420px'"
      :with-header="false"
      append-to-body
      class="tag-manager-drawer"
    >
      <section class="tag-manager" aria-labelledby="tag-manager-title">
        <header class="tag-manager-header">
          <div>
            <span class="section-kicker">ORGANIZE</span>
            <h2 id="tag-manager-title">标签管理</h2>
            <p>点击标签筛选笔记，也可以重命名或移除标签。</p>
          </div>
          <button
            type="button"
            class="tag-manager-close"
            aria-label="关闭标签管理"
            @click="tagManagerOpen = false"
          >
            <el-icon><Close /></el-icon>
          </button>
        </header>

        <div class="tag-manager-summary">
          <span>{{ tagStats.length }} 个标签</span>
          <span>{{ tagStats.reduce((total, item) => total + item.count, 0) }} 次关联</span>
        </div>

        <div v-if="tagStats.length" class="tag-manager-list">
          <article v-for="item in tagStats" :key="item.name" class="tag-manager-item">
            <button type="button" class="tag-manager-open" @click="openTagNotes(item.name)">
              <span class="tag-manager-name"># {{ item.name }}</span>
              <small>{{ item.count }} 篇笔记</small>
            </button>
            <div class="tag-manager-actions">
              <el-button
                text
                :icon="EditPen"
                :aria-label="`重命名标签：${item.name}`"
                @click="renameManagedTag(item)"
              />
              <el-popconfirm
                :title="`移除标签「${item.name}」？笔记不会被删除。`"
                width="250"
                @confirm="removeManagedTag(item)"
              >
                <template #reference>
                  <el-button
                    text
                    :icon="Delete"
                    :aria-label="`删除标签：${item.name}`"
                  />
                </template>
              </el-popconfirm>
            </div>
          </article>
        </div>

        <div v-else class="tag-manager-empty">
          <el-icon><Tickets /></el-icon>
          <h3>还没有标签</h3>
          <p>编辑笔记时输入标签并按回车，即可创建。</p>
        </div>
      </section>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import { Plus, Check, Delete, Document, EditPen, ArrowLeft, FolderOpened, Notebook, Clock, Search, Tickets, Close, Calendar, Edit, View } from '@element-plus/icons-vue'
import {
  getNotes,
  createNote,
  updateNote,
  deleteNote,
  getNote,
  getTags,
  renameTag as renameTagApi,
  deleteTag as deleteTagApi
} from '@/api/note'
import { ElMessage, ElMessageBox } from 'element-plus'
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
const editMode = ref(false)
const isMobile = ref(false)
const sidebarOpen = ref(false)
const allTags = ref([])
const tagStats = ref([])
const activeTag = ref('')
const tagDraft = ref('')
const tagManagerOpen = ref(false)

const availableTagSuggestions = computed(() =>
  allTags.value.filter(tag => !currentTagInput.value.includes(tag))
)

const wordCount = computed(() => {
  const html = currentNote.value?.content || ''
  // 去掉 HTML 标签,再去掉空白,剩余字符数即为"字数"
  return html.replace(/<[^>]*>/g, '').replace(/\s/g, '').length
})

const now = new Date()
const greeting = computed(() => {
  const hour = now.getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})
const dashboardDate = computed(() =>
  new Intl.DateTimeFormat('zh-CN', {
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  }).format(now)
)
let searchTimer = null
let autoSaveTimer = null
const AUTO_SAVE_DELAY = 3000

// 最近笔记（最多显示5个）
const recentNotes = ref([])

// 加载所有笔记
function loadAllNotes() {
  router.push('/notes-list')
}

function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveNote(true)
  }
}

function handleBeforeUnload(e) {
  if (isDirty.value) {
    e.preventDefault()
    e.returnValue = ''
  }
}

onBeforeRouteLeave((to, from, next) => {
  if (isDirty.value) {
    clearTimeout(autoSaveTimer)
    ElMessageBox.confirm('笔记内容尚未保存，确定要离开吗？', '未保存的修改', {
      confirmButtonText: '离开',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      isDirty.value = false
      next()
    }).catch(() => {
      next(false)
    })
  } else {
    next()
  }
})

onMounted(async () => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('beforeunload', handleBeforeUnload)
  await loadNotes()
  await loadTags()
  await openNoteFromQuery()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('beforeunload', handleBeforeUnload)
  clearTimeout(autoSaveTimer)
})

watch(() => route.query.note_id, async (newId) => {
  if (!newId) return
  if (String(currentNote.value?.id) === String(newId)) return
  await nextTick()
  await loadNotes()
  await loadTags()
  await openNoteFromQuery()
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
    const { data } = await getTags()
    allTags.value = data.tags || []
    tagStats.value = data.items || allTags.value.map(name => ({
      name,
      count: notes.value.filter(note => splitTags(note.tags).includes(name)).length
    }))
  } catch (e) {}
}

function openTagNotes(tag) {
  tagManagerOpen.value = false
  router.push({ path: '/notes-list', query: { tag } })
}

async function renameManagedTag(item) {
  try {
    const { value } = await ElMessageBox.prompt(
      `为标签「${item.name}」输入新名称`,
      '重命名标签',
      {
        confirmButtonText: '保存',
        cancelButtonText: '取消',
        inputValue: item.name,
        inputPlaceholder: '标签名称',
        inputValidator: value => {
          const name = value?.trim()
          if (!name) return '标签名称不能为空'
          if (name.length > 80) return '标签名称不能超过 80 个字符'
          return true
        }
      }
    )

    const nextName = value.trim()
    if (nextName === item.name) return
    await renameTagApi(item.name, nextName)
    ElMessage.success('标签已重命名')
    await Promise.all([loadNotes(), loadTags()])
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error('标签重命名失败')
    }
  }
}

async function removeManagedTag(item) {
  try {
    await deleteTagApi(item.name)
    ElMessage.success('标签已移除，笔记内容不受影响')
    await Promise.all([loadNotes(), loadTags()])
  } catch {
    ElMessage.error('标签删除失败')
  }
}

async function openNoteFromQuery() {
  const noteId = route.query.note_id
  if (!noteId) return
  try {
    const { data } = await getNote(noteId)
    currentNote.value = { ...data }
    currentTagInput.value = data.tags ? data.tags.split(',').map(t => t.trim()).filter(Boolean) : []
    tagDraft.value = ''
    isDirty.value = false
    editMode.value = route.query.edit === '1'
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
    recentNotes.value = notes.value.slice(0, 5)
    currentNote.value = { ...data }
    currentTagInput.value = []
    tagDraft.value = ''
    isDirty.value = false
    editMode.value = true
    sidebarOpen.value = false
    router.replace({ path: '/', query: { note_id: data.id, edit: '1' } })
  } catch (e) {
    ElMessage.error('创建笔记失败')
  }
}

async function selectNote(note) {
  // If there are unsaved changes, save them silently first
  if (isDirty.value && currentNote.value) {
    clearTimeout(autoSaveTimer)
    const prevNote = currentNote.value
    try {
      const tags = currentTagInput.value.join(', ')
      await updateNote(prevNote.id, {
        title: prevNote.title,
        content: prevNote.content,
        tags
      })
    } catch (e) {/* ignore save errors on switch */}
  }
  currentNote.value = { ...note }
  currentTagInput.value = note.tags ? note.tags.split(',').map(t => t.trim()).filter(Boolean) : []
  tagDraft.value = ''
  isDirty.value = false
  editMode.value = false
  router.replace({ path: '/', query: { note_id: note.id } })
  if (isMobile.value) {
    sidebarOpen.value = false
  }
}

function backToList() {
  currentNote.value = null
  tagDraft.value = ''
  sidebarOpen.value = true
  router.replace({ path: '/' })
}

function markDirty() {
  if (!isDirty.value) {
    isDirty.value = true
  }
  clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => saveNote(false), AUTO_SAVE_DELAY)
}

function addCurrentTag() {
  const tags = tagDraft.value
    .split(/[,，]/)
    .map(tag => tag.trim())
    .filter(Boolean)

  if (tags.length === 0) return

  let changed = false
  for (const tag of tags) {
    if (!currentTagInput.value.includes(tag)) {
      currentTagInput.value.push(tag)
      changed = true
    }
  }
  tagDraft.value = ''
  if (changed) markDirty()
}

function handleTagKeydown(event) {
  if (event.key !== 'Enter' && event.key !== ',' && event.key !== '，') return
  event.preventDefault()
  addCurrentTag()
}

function removeCurrentTag(tag) {
  currentTagInput.value = currentTagInput.value.filter(item => item !== tag)
  markDirty()
}

function enterEditMode() {
  editMode.value = true
  if (currentNote.value) {
    router.replace({ path: '/', query: { note_id: currentNote.value.id, edit: '1' } })
  }
}

async function quickEditNote(note) {
  await selectNote(note)
  editMode.value = true
  router.replace({ path: '/', query: { note_id: note.id, edit: '1' } })
  if (isMobile.value) {
    sidebarOpen.value = false
  }
}

async function exitEditMode() {
  if (isDirty.value) {
    clearTimeout(autoSaveTimer)
    await saveNote(false)
  }
  editMode.value = false
  if (currentNote.value) {
    router.replace({ path: '/', query: { note_id: currentNote.value.id } })
  }
}

async function saveNote(showToast = true) {
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
    clearTimeout(autoSaveTimer)
    if (showToast) ElMessage.success('保存成功')
    await loadNotes()
    await loadTags()
  } catch (e) {
    if (showToast) ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function removeNote() {
  if (!currentNote.value) return
  try {
    await deleteNote(currentNote.value.id)
    currentNote.value = null
    router.replace({ path: '/' })
    ElMessage.success('删除成功')
    await loadNotes()
    await loadTags()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

import { stripHtml, splitTags, formatDate } from '@/utils/note'

function formatFullDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.notes-container {
  display: flex;
  height: calc(100vh - 64px - 48px);
  position: relative;
  background: var(--bg-secondary);
}

.notes-sidebar {
  width: 320px;
  background: var(--bg-primary);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color);
}

.sidebar-header {
  padding: 12px;
}

.create-btn {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  font-weight: 600;
  background: var(--gradient-primary);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s;
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.sidebar-search {
  padding: 0 12px 12px;
}

.search-box {
  position: relative;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  font-size: 18px;
  z-index: 1;
}

.search-input :deep(.el-input__wrapper) {
  padding-left: 44px !important;
  border-radius: 12px;
  box-shadow: 0 0 0 1px var(--border-color);
  background: var(--bg-secondary);
  transition: all 0.3s;
}

.search-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--primary-color);
  background: var(--bg-primary);
}

.search-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--primary-color);
  background: var(--bg-primary);
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
  color: var(--text-primary);
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
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border: none;
  transition: all 0.2s;
}

.tag-item:hover {
  background: var(--border-color);
  transform: translateY(-1px);
}

.tag-item.active {
  background: var(--gradient-primary);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.clear-tag {
  background: #fef2f2;
  color: var(--danger-color);
}

.clear-tag:hover {
  background: #fee2e2;
}

html[data-theme="dark"] .clear-tag {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

html[data-theme="dark"] .clear-tag:hover {
  background: rgba(239, 68, 68, 0.25);
}

.note-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 4px 12px;
}

.note-item {
  padding: 16px;
  margin-bottom: 8px;
  cursor: pointer;
  border-radius: 12px;
  transition: all 0.2s;
  background: var(--bg-tertiary);
}

.note-item:hover {
  background: var(--bg-tertiary);
  transform: translateY(-1px);
}

.note-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border: 2px solid var(--primary-color);
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
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.note-item-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--primary-color);
  opacity: 0;
  flex-shrink: 0;
}

.note-item.active .note-item-dot {
  opacity: 1;
}

.note-item-preview {
  font-size: 13px;
  color: var(--text-secondary);
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
  color: var(--text-muted);
}

.note-item-tags {
  display: flex;
  gap: 4px;
}

.mini-tag {
  padding: 2px 8px;
  background: var(--border-color);
  color: var(--text-secondary);
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
  color: var(--text-muted);
  margin-bottom: 16px;
}

.empty-list h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 6px;
}

.empty-list p {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
}

.notes-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  min-width: 0;
  overflow: hidden;
}

.editor-mobile-top {
  padding: 12px 16px 0;
  flex-shrink: 0;
  border-bottom: 1px solid var(--border-light);
}

.back-btn {
  color: var(--primary-color);
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
  color: var(--primary-color);
  font-weight: 500;
  font-size: 14px;
}

.back-to-list-btn:hover {
  color: var(--primary-hover);
}

.title-input :deep(.el-input__wrapper) {
  box-shadow: none !important;
  font-size: 28px;
  font-weight: 700;
  padding: 0;
}

.title-input :deep(.el-input__inner) {
  color: var(--text-primary);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 32px 16px;
  border-bottom: 1px solid var(--border-light);
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
  color: var(--text-secondary);
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
  color: var(--text-secondary);
}

.tags-select {
  min-width: 180px;
}

.editor-edit-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  /* 同上：避免 overflow-y 隐式触发 overflow-x: auto 产生水平滚动条 */
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
}

.editor-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  color: var(--text-secondary);
  text-align: center;
  padding: 8px 16px;
  overflow-y: auto;
}

.editor-empty h1 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.editor-empty p {
  font-size: 15px;
  color: var(--text-secondary);
  margin: 0 0 16px;
}

.empty-btn {
  border-radius: 12px;
  font-weight: 600;
  background: var(--gradient-primary);
  border: none;
}

/* 仪表板样式 */
.dashboard-container {
  width: 100%;
  padding: 8px 16px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 24px;
}

.dashboard-header h1 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 4px 0 2px;
  letter-spacing: -0.5px;
}

.dashboard-header p {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 8px;
}

.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 8px;
  margin-bottom: 16px;
}

.stat-card {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 20px 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  border: 1px solid var(--border-light);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 48px;
  height: 48px;
  background: var(--gradient-primary);
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
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: 2px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.recent-notes-section {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 20px 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid var(--border-light);
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
  color: var(--text-primary);
  margin: 0;
}

.more-link {
  color: var(--primary-color) !important;
  font-weight: 500;
  font-size: 12px;
}

.more-link:hover {
  color: var(--primary-hover) !important;
}

.recent-notes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recent-note-item {
  padding: 16px;
  border-radius: 12px;
  background: var(--bg-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  border: 2px solid transparent;
}

.recent-note-item:hover {
  background: var(--bg-tertiary);
  transform: translateY(-1px);
  border-color: var(--border-color);
}

.recent-note-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recent-note-preview {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.recent-note-time {
  font-size: 11px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

.recent-note-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: 4px;
}

.recent-note-edit-btn {
  flex-shrink: 0;
  border-radius: 8px;
  font-weight: 500;
}

.no-recent-notes {
  text-align: center;
  padding: 32px 16px;
  color: var(--text-muted);
  font-size: 13px;
}

.empty-btn {
  margin-top: 16px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  padding: 10px 20px;
  background: var(--gradient-primary);
  border: none;
}

/* View mode styles */
.note-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  /* 同上：预览态同样避免出现多余的水平滚动条 */
  overflow-x: hidden;
  padding: 32px 40px;
}

.view-header {
  position: sticky;
  top: 0;
  z-index: 5;
  background: var(--bg-primary);
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border-light);
}

.view-title-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.view-title {
  flex: 1;
  min-width: 0;
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.3;
  word-break: break-word;
}

.view-edit-btn {
  flex-shrink: 0;
  border-radius: 10px;
  font-weight: 600;
  background: var(--gradient-primary);
  border: none;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.view-meta {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
  color: var(--text-secondary);
  font-size: 14px;
}

.view-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.view-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.view-content {
  flex: 1;
  line-height: 1.8;
  font-size: 16px;
  color: var(--text-primary);
  padding: 8px 0 24px;
  /* 同编辑器：强制超长串/URL 折行，兼容老内核，避免横向撑宽 */
  word-break: break-word;
  overflow-wrap: anywhere;
}

.view-content :deep(h1) {
  font-size: 2em;
  margin: 0.67em 0;
}

.view-content :deep(h2) {
  font-size: 1.5em;
  margin: 0.75em 0;
}

.view-content :deep(h3) {
  font-size: 1.17em;
  margin: 0.83em 0;
}

.view-content :deep(p) {
  margin: 1em 0;
}

.view-content :deep(ul),
.view-content :deep(ol) {
  padding-left: 2em;
  margin: 1em 0;
}

.view-content :deep(li) {
  margin: 0.5em 0;
}

.view-content :deep(blockquote) {
  border-left: 4px solid var(--primary-color);
  padding-left: 16px;
  margin: 1em 0;
  color: var(--text-secondary);
  font-style: italic;
}

.view-content :deep(code) {
  background: var(--bg-tertiary);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
}

.view-content :deep(pre) {
  background: var(--bg-tertiary);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1em 0;
  /* 代码块保留自身横向滚动，不继承正文的强制折行 */
  white-space: pre;
  word-break: normal;
  overflow-wrap: normal;
}

.view-content :deep(pre code) {
  background: none;
  padding: 0;
}

.view-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 1em 0;
}

.view-footer {
  display: flex;
  gap: 12px;
  padding: 24px 0 8px;
  border-top: 1px solid var(--border-light);
}

/* 移动端特定优化 */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 4px 4px;
  }
  
  .dashboard-header h1 {
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
    height: calc(100vh - 56px - 16px);
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

  .tags-select {
    flex: 1;
  }

  .note-view {
    padding: 16px;
  }

  .view-title {
    font-size: 24px;
  }

  .view-title-row {
    flex-wrap: wrap;
    row-gap: 8px;
  }

  .view-edit-btn {
    font-size: 13px;
    padding: 6px 14px;
    height: 32px;
  }

  .view-meta {
    flex-direction: column;
    gap: 8px;
  }

  .view-footer {
    flex-direction: column;
    gap: 8px;
  }

  .view-footer .el-button {
    width: 100%;
  }
}

/* 移动端浮动新建按钮 */
.mobile-fab {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99;
  transition: all 0.3s ease;
}

.mobile-fab:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.mobile-fab:active {
  transform: scale(0.95);
}

/* Also show the FAB on the notes-list page when sidebar is hidden */
@media (max-width: 768px) {
  .mobile-fab {
    bottom: 20px;
    right: 20px;
    width: 52px;
    height: 52px;
  }
}

/* 底部操作栏(桌面 + 移动通用) */
.editor-bottom-bar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 16px;
  padding-bottom: calc(8px + env(safe-area-inset-bottom, 0px));
  background: var(--bg-primary);
  border-top: 1px solid var(--border-color);
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.1);
}

.bar-word-count {
  flex-shrink: 0;
  font-size: 13px;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
  user-select: none;
  white-space: nowrap;
}

.bar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  justify-content: flex-end;
  min-width: 0;
}

.bar-btn {
  flex: 1;
  height: 44px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 15px;
}

.bar-save-btn {
  flex: 2;
  background: var(--gradient-primary);
  border: none;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.bar-delete-btn {
  flex: 0 0 44px;
  padding: 0;
}

@media (max-width: 768px) {
  .editor-bottom-bar {
    padding: 8px 12px;
  }
  .bar-word-count {
    font-size: 12px;
  }
}
</style>
