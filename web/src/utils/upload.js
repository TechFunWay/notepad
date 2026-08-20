const appBase = import.meta.env.BASE_URL

export function toRuntimeUploadUrl(url) {
  if (!url?.startsWith('/uploads/') || appBase === '/') return url
  return `${appBase}${url.slice(1)}`
}

export function toRuntimeUploadHtml(html) {
  if (!html || appBase === '/') return html
  return html.replace(/(\bsrc=["'])\/uploads\//gi, `$1${appBase}uploads/`)
}

export function toStoredUploadHtml(html) {
  if (!html || appBase === '/') return html
  const escapedBase = appBase.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return html.replace(
    new RegExp(`(\\bsrc=["'])${escapedBase}uploads/`, 'gi'),
    '$1/uploads/'
  )
}
