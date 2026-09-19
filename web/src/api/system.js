import api from '@/api/request'

// 版本信息（含 started_at，用于前端区分「本次运行」与「重启后」）
export async function getVersion() {
  return api.get('/version')
}

// 「已支持」上报一次匿名支持计数（金额仅接收端统计用，不做支付核验）
export async function donateSupport(amount) {
  return api.post('/donate/support', { amount: amount || 0 })
}
