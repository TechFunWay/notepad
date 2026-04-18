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

export function showConfirm(message, options = {}) {
  return new Promise((resolve) => {
    const { title = '提示', confirmText = '确定', cancelText = '取消', type = 'warning' } = options

    const overlay = document.createElement('div')
    overlay.className = 'mobile-confirm-overlay'
    overlay.style.cssText = `
      position: fixed;
      top: 0; left: 0; right: 0; bottom: 0;
      background: rgba(0, 0, 0, 0.5);
      z-index: 20000;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
      opacity: 0;
      transition: opacity 0.2s ease;
    `

    const colors = {
      warning: '#f59e0b',
      error: '#ef4444',
      info: '#667eea'
    }

    const dialog = document.createElement('div')
    dialog.className = 'mobile-confirm-dialog'
    dialog.style.cssText = `
      background: white;
      border-radius: 16px;
      width: 100%;
      max-width: 320px;
      padding: 28px 24px 20px;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
      transform: scale(0.9);
      transition: transform 0.25s ease;
      text-align: center;
    `

    const iconContainer = document.createElement('div')
    iconContainer.style.cssText = `
      width: 48px; height: 48px;
      border-radius: 50%;
      background: ${colors[type] || colors.info}20;
      display: flex; align-items: center; justify-content: center;
      margin: 0 auto 16px;
      font-size: 24px;
    `
    const iconMap = { warning: '⚠', error: '✕', info: 'ℹ' }
    iconContainer.textContent = iconMap[type] || iconMap.warning
    dialog.appendChild(iconContainer)

    if (title) {
      const titleEl = document.createElement('div')
      titleEl.style.cssText = `font-size: 17px; font-weight: 600; color: #1f2937; margin-bottom: 8px;`
      titleEl.textContent = title
      dialog.appendChild(titleEl)
    }

    const msgEl = document.createElement('div')
    msgEl.style.cssText = `font-size: 14px; color: #6b7280; margin-bottom: 24px; line-height: 1.5;`
    msgEl.textContent = message
    dialog.appendChild(msgEl)

    const btnGroup = document.createElement('div')
    btnGroup.style.cssText = `display: flex; gap: 12px;`

    const cancelBtn = document.createElement('button')
    cancelBtn.textContent = cancelText
    cancelBtn.style.cssText = `
      flex: 1; padding: 12px; border: 1px solid #e5e7eb; border-radius: 10px;
      background: white; color: #6b7280; font-size: 15px; font-weight: 500;
      cursor: pointer; transition: all 0.15s;
    `
    cancelBtn.addEventListener('click', () => {
      overlay.style.opacity = '0'
      dialog.style.transform = 'scale(0.9)'
      setTimeout(() => { document.body.removeChild(overlay); resolve(false) }, 200)
    })
    btnGroup.appendChild(cancelBtn)

    const confirmBtn = document.createElement('button')
    confirmBtn.textContent = confirmText
    confirmBtn.style.cssText = `
      flex: 1; padding: 12px; border: none; border-radius: 10px;
      background: ${colors[type] || colors.info}; color: white; font-size: 15px; font-weight: 500;
      cursor: pointer; transition: all 0.15s;
    `
    confirmBtn.addEventListener('click', () => {
      overlay.style.opacity = '0'
      dialog.style.transform = 'scale(0.9)'
      setTimeout(() => { document.body.removeChild(overlay); resolve(true) }, 200)
    })
    btnGroup.appendChild(confirmBtn)

    dialog.appendChild(btnGroup)
    overlay.appendChild(dialog)
    document.body.appendChild(overlay)

    cancelBtn.addEventListener('touchstart', () => { cancelBtn.style.background = '#f3f4f6' })
    cancelBtn.addEventListener('touchend', () => { cancelBtn.style.background = 'white' })
    confirmBtn.addEventListener('touchstart', () => { confirmBtn.style.opacity = '0.8' })
    confirmBtn.addEventListener('touchend', () => { confirmBtn.style.opacity = '1' })

    requestAnimationFrame(() => {
      overlay.style.opacity = '1'
      dialog.style.transform = 'scale(1)'
    })
  })
}
