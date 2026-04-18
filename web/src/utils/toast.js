const TOAST_DURATION = 2000

function showToast(message, type = 'info') {
  let container = document.getElementById('mobile-toast-container')
  if (!container) {
    container = document.createElement('div')
    container.id = 'mobile-toast-container'
    container.style.cssText = `
      position: fixed;
      top: 60px;
      left: 16px;
      right: 16px;
      z-index: 10000;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      pointer-events: none;
    `
    document.body.appendChild(container)
  }

  const icons = {
    success: '✓',
    error: '✕',
    warning: '!'
  }

  const colors = {
    success: '#10b981',
    error: '#ef4444',
    warning: '#f59e0b',
    info: '#667eea'
  }

  const toast = document.createElement('div')
  toast.className = 'mobile-toast'
  toast.style.cssText = `
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: ${colors[type] || colors.info};
    color: white;
    border-radius: 10px;
    font-size: 14px;
    font-weight: 500;
    width: 100%;
    box-sizing: border-box;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    opacity: 0;
    transform: translateY(-20px);
    transition: all 0.25s ease;
    text-align: left;
  `

  const iconSpan = document.createElement('span')
  iconSpan.style.cssText = `
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.25);
    font-size: 12px;
    font-weight: 700;
    flex-shrink: 0;
  `
  iconSpan.textContent = icons[type] || icons.info
  toast.appendChild(iconSpan)

  const textSpan = document.createElement('span')
  textSpan.textContent = message
  textSpan.style.flex = '1'
  toast.appendChild(textSpan)

  container.appendChild(toast)

  requestAnimationFrame(() => {
    toast.style.opacity = '1'
    toast.style.transform = 'translateY(0)'
  })

  setTimeout(() => {
    toast.style.opacity = '0'
    toast.style.transform = 'translateY(-20px)'
    setTimeout(() => {
      if (toast.parentNode) {
        toast.parentNode.removeChild(toast)
      }
      if (container.children.length === 0 && container.parentNode) {
        container.parentNode.removeChild(container)
      }
    }, 250)
  }, TOAST_DURATION)
}

export const toast = {
  success: (msg) => showToast(msg, 'success'),
  error: (msg) => showToast(msg, 'error'),
  warning: (msg) => showToast(msg, 'warning'),
  info: (msg) => showToast(msg, 'info')
}

export function isMobile() {
  return window.innerWidth <= 768
}
