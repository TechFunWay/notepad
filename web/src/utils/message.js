import { ElMessage, ElMessageBox } from 'element-plus'
import { toast, isMobile, showConfirm } from './toast'

export const message = {
  success: (msg) => isMobile() ? toast.success(msg) : ElMessage.success(msg),
  error: (msg) => isMobile() ? toast.error(msg) : ElMessage.error(msg),
  warning: (msg) => isMobile() ? toast.warning(msg) : ElMessage.warning(msg),
  info: (msg) => isMobile() ? toast.info(msg) : ElMessage.info(msg),
  confirm: (msg, options = {}) => {
    if (isMobile()) {
      return showConfirm(msg, options)
    }
    return ElMessageBox.confirm(msg, options.title || '提示', {
      confirmButtonText: options.confirmButtonText || '确定',
      cancelButtonText: options.cancelButtonText || '取消',
      type: options.type || 'warning'
    }).then(() => true).catch(() => false)
  }
}
