import { ref, watch } from 'vue'

const savedTheme = localStorage.getItem('theme')
const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches
const isDark = ref(savedTheme ? savedTheme === 'dark' : prefersDark)

function applyTheme(dark) {
  if (dark) {
    document.documentElement.setAttribute('data-theme', 'dark')
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
    document.documentElement.classList.remove('dark')
  }
  document.querySelector('meta[name="theme-color"]')?.setAttribute(
    'content',
    dark ? '#09191e' : '#0891b2'
  )
}

function initTheme() {
  applyTheme(isDark.value)
}

function toggleTheme() {
  isDark.value = !isDark.value
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  applyTheme(isDark.value)
}

watch(isDark, () => {
  applyTheme(isDark.value)
})

export function useTheme() {
  return { isDark, toggleTheme, initTheme }
}
