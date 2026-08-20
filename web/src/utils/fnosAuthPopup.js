export const FNOS_AUTH_RESULT_TYPE = 'notepad-fnos-auth-result'

/**
 * 打开飞牛 NAS 系统登录页弹窗（不跳转当前页面）。
 * 用户完成登录后，网关按 redirect_uri 跳回应用回调页，由回调页通过
 * postMessage 把结果回传，再自动关闭弹窗。
 * 解析值：
 *   'authorized'  已在弹窗中完成登录授权
 *   'cancelled'   用户直接关闭弹窗，未完成授权
 *   'blocked'     弹窗被浏览器拦截（未弹出）
 */
export function openFnOSAuthPopup({ width = 500, height = 560 } = {}) {
  // 回调页必须带应用 base 前缀（如 /app/techfunway-notepad/fnos-auth-callback），
  // 否则登录后网关会跳回站点根路径，回调页无法加载。
  const callbackURL = new URL(
    `${import.meta.env.BASE_URL}fnos-auth-callback`,
    window.location.origin
  ).href
  const authURL = new URL(
    `/login?redirect_uri=${encodeURIComponent(callbackURL)}`,
    window.location.origin
  ).href

  const popup = window.open(
    authURL,
    'notepad-fnos-auth',
    `width=${width},height=${height},left=${Math.max(0, (window.screen.width - width) / 2)},top=${Math.max(0, (window.screen.height - height) / 2)}`
  )

  return new Promise((resolve) => {
    if (!popup) {
      resolve('blocked')
      return
    }

    let settled = false
    let checkClosed = null

    const finish = (status) => {
      if (settled) return
      settled = true
      window.removeEventListener('message', onMessage)
      if (checkClosed) window.clearInterval(checkClosed)
      if (!popup.closed) popup.close()
      resolve(status)
    }

    const onMessage = (event) => {
      if (event.data && event.data.type === FNOS_AUTH_RESULT_TYPE) {
        finish('authorized')
      }
    }

    window.addEventListener('message', onMessage)
    checkClosed = window.setInterval(() => {
      if (popup.closed) {
        finish('cancelled')
      }
    }, 400)
  })
}
