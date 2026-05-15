import { ref, watch } from 'vue'

const isDark = ref(localStorage.getItem('theme') === 'dark')

function applyTheme(dark) {
  if (dark) {
    document.documentElement.setAttribute('data-theme', 'dark')
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
    document.documentElement.classList.remove('dark')
  }
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
