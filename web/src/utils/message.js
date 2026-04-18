import { ElMessage } from 'element-plus'
import { toast, isMobile } from './toast'

export const message = {
  success: (msg) => isMobile() ? toast.success(msg) : ElMessage.success(msg),
  error: (msg) => isMobile() ? toast.error(msg) : ElMessage.error(msg),
  warning: (msg) => isMobile() ? toast.warning(msg) : ElMessage.warning(msg),
  info: (msg) => isMobile() ? toast.info(msg) : ElMessage.info(msg)
}
