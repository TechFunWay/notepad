<template>
  <div class="tiptap-editor" v-if="editor">
    <div class="editor-toolbar">
      <div class="toolbar-group toolbar-marks">
        <button type="button" aria-label="粗体" @click="editor.chain().focus().toggleBold().run()" :class="{ active: editor.isActive('bold') }" title="粗体">
          <strong>B</strong>
        </button>
        <button type="button" aria-label="斜体" @click="editor.chain().focus().toggleItalic().run()" :class="{ active: editor.isActive('italic') }" title="斜体">
          <em>I</em>
        </button>
        <button type="button" aria-label="下划线" @click="editor.chain().focus().toggleUnderline().run()" :class="{ active: editor.isActive('underline') }" title="下划线">
          <u>U</u>
        </button>
        <button type="button" aria-label="删除线" @click="editor.chain().focus().toggleStrike().run()" :class="{ active: editor.isActive('strike') }" title="删除线">
          <s>S</s>
        </button>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-headings">
        <button type="button" aria-label="一级标题" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ active: editor.isActive('heading', { level: 1 }) }" title="一级标题">
          H1
        </button>
        <button type="button" aria-label="二级标题" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ active: editor.isActive('heading', { level: 2 }) }" title="二级标题">
          H2
        </button>
        <button type="button" aria-label="三级标题" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ active: editor.isActive('heading', { level: 3 }) }" title="三级标题">
          H3
        </button>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-lists">
        <button type="button" aria-label="无序列表" @click="editor.chain().focus().toggleBulletList().run()" :class="{ active: editor.isActive('bulletList') }" title="无序列表">
          •
        </button>
        <button type="button" aria-label="有序列表" @click="editor.chain().focus().toggleOrderedList().run()" :class="{ active: editor.isActive('orderedList') }" title="有序列表">
          1.
        </button>
        <button type="button" aria-label="引用" @click="editor.chain().focus().toggleBlockquote().run()" :class="{ active: editor.isActive('blockquote') }" title="引用">
          "
        </button>
        <button type="button" aria-label="分割线" @click="editor.chain().focus().setHorizontalRule().run()" title="分割线">
          —
        </button>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-alignment">
        <button type="button" aria-label="左对齐" @click="editor.chain().focus().setTextAlign('left').run()" :class="{ active: editor.isActive({ textAlign: 'left' }) }" title="左对齐">
          ←
        </button>
        <button type="button" aria-label="居中对齐" @click="editor.chain().focus().setTextAlign('center').run()" :class="{ active: editor.isActive({ textAlign: 'center' }) }" title="居中对齐">
          ─
        </button>
        <button type="button" aria-label="右对齐" @click="editor.chain().focus().setTextAlign('right').run()" :class="{ active: editor.isActive({ textAlign: 'right' }) }" title="右对齐">
          →
        </button>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-color">
        <input type="color" :value="textColor || '#000000'" @input="setColor($event.target.value)" class="color-input" title="文字颜色" aria-label="选择文字颜色" />
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-media">
        <button type="button" aria-label="插入图片" @click="triggerImageUpload" title="插入图片">
          <el-icon><Picture /></el-icon>
        </button>
        <input ref="imageInput" type="file" accept="image/*" @change="handleImageSelect" style="display: none" />
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-group toolbar-history">
        <button type="button" aria-label="撤销" @click="editor.chain().focus().undo().run()" title="撤销">
          <el-icon><RefreshLeft /></el-icon>
        </button>
        <button type="button" aria-label="重做" @click="editor.chain().focus().redo().run()" title="重做">
          <el-icon><RefreshRight /></el-icon>
        </button>
      </div>
    </div>
    <editor-content class="editor-content" :editor="editor" />
  </div>
</template>

<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import { computed, watch, ref } from 'vue'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Color from '@tiptap/extension-color'
import TextStyle from '@tiptap/extension-text-style'
import Underline from '@tiptap/extension-underline'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'
import { ElMessage } from 'element-plus'
import { Picture, RefreshLeft, RefreshRight } from '@element-plus/icons-vue'
import api from '@/api/request'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: '开始写作...'
    }),
    Color,
    TextStyle,
    Underline,
    TextAlign.configure({
      types: ['heading', 'paragraph']
    }),
    Image.configure({
      inline: false,
      allowBase64: true
    })
  ],
  content: props.modelValue,
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

const textColor = computed(() => {
  if (!editor.value) return '#000000'
  return editor.value.getAttributes('textStyle').color || '#000000'
})

watch(() => props.modelValue, (newValue) => {
  if (editor.value && newValue !== editor.value.getHTML()) {
    editor.value.commands.setContent(newValue)
  }
})

const imageInput = ref(null)

function triggerImageUpload() {
  imageInput.value?.click()
}

function handleImageSelect(event) {
  const file = event.target.files[0]
  if (file) {
    uploadImage(file)
    event.target.value = ''
  }
}

function handlePaste(event) {
  const items = event.clipboardData?.items
  if (!items) return

  for (const item of items) {
    if (item.type.indexOf('image') === 0) {
      event.preventDefault()
      const file = item.getAsFile()
      if (file) {
        uploadImage(file)
      }
      break
    }
  }
}

async function uploadImage(file) {
  if (!editor.value) return
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const { data } = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    editor.value.chain().focus().insertContent(`<img src="${data.url}" alt="图片" />`).run()
  } catch (error) {
    ElMessage.error('图片上传失败，请重试')
    console.error('Image upload error:', error)
  }
}

setTimeout(() => {
  const editorElement = document.querySelector('.ProseMirror')
  if (editorElement) {
    editorElement.addEventListener('paste', handlePaste)
  }
}, 100)

function setColor(color) {
  if (!editor.value) return
  editor.value.chain().focus().setColor(color).run()
}
</script>

<style scoped>
.tiptap-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-tertiary);
  background: color-mix(in srgb, var(--bg-tertiary) 88%, transparent);
  flex-wrap: wrap;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 4px;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 8px;
}

.editor-toolbar button {
  width: 44px;
  height: 44px;
  border: none;
  background: transparent;
  border-radius: 9px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 600;
  transition: all 0.15s;
}

.editor-toolbar button:hover {
  background: var(--bg-primary);
  color: var(--primary-color);
}

.editor-toolbar button.active {
  background: var(--primary-light);
  color: var(--primary-strong);
  box-shadow: inset 0 0 0 1px var(--border-color);
  box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--primary-color) 24%, transparent);
}

.color-input {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  padding: 0;
}

@media (max-width: 768px) {
  .editor-toolbar {
    gap: 8px;
    padding: 8px 10px;
  }

  .toolbar-group {
    gap: 8px;
  }

  .toolbar-divider {
    display: none;
  }
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  /* 仅设 overflow-y 时，overflow-x 的 visible 会被规范隐式计算为 auto，
     导致内容略宽即出现不该有的水平滚动条；显式置为 hidden 消除。 */
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
  padding: 32px clamp(20px, 5vw, 72px);
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 100%;
  max-width: 100%;
  /* 强制超长串/URL 折行，兼容不支持 break-spaces 的老内核(如搜狗极速模式)，
     避免内容横向撑宽后被裁切 */
  word-break: break-word;
  overflow-wrap: anywhere;
  width: min(100%, 780px);
  margin: 0 auto;
  color: var(--text-primary);
  font-size: 17px;
  line-height: 1.85;
}

.editor-content :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  color: var(--text-muted);
  float: left;
  height: 0;
  pointer-events: none;
  content: attr(data-placeholder);
}

.editor-content :deep(.ProseMirror p) {
  margin: 0.5em 0;
  line-height: 1.7;
}

.editor-content :deep(.ProseMirror h1) {
  font-size: 2em;
  margin: 0.67em 0;
  font-weight: 700;
}

.editor-content :deep(.ProseMirror h2) {
  font-size: 1.5em;
  margin: 0.83em 0;
  font-weight: 600;
}

.editor-content :deep(.ProseMirror h3) {
  font-size: 1.17em;
  margin: 1em 0;
  font-weight: 600;
}

.editor-content :deep(.ProseMirror ul),
.editor-content :deep(.ProseMirror ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}

.editor-content :deep(.ProseMirror li) {
  margin: 0.25em 0;
}

.editor-content :deep(.ProseMirror blockquote) {
  border-left: 3px solid var(--primary-color);
  padding-left: 1em;
  margin: 1em 0;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  padding: 12px 16px;
  border-radius: 8px;
}

.editor-content :deep(.ProseMirror hr) {
  border: none;
  border-top: 1px solid var(--border-color);
  margin: 2em 0;
}

.editor-content :deep(.ProseMirror a) {
  color: var(--primary-color);
  text-decoration: underline;
}

.editor-content :deep(.ProseMirror img) {
  max-width: 100%;
  max-height: 90vh;
  height: auto;
  display: block;
  border-radius: 8px;
  object-fit: contain;
}
</style>
