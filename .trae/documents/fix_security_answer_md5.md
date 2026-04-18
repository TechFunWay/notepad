# 修复安全问题MD5加密不一致问题

## 问题分析
前端在设置和验证安全问题答案时，MD5加密不一致：
- **忘记密码页面**：安全答案使用 `md5()` 加密后传输 ✓
- **注册页面**：安全答案直接明文传输 ✗
- **个人中心（修改安全问题）**：安全答案直接明文传输 ✗

这导致：
1. 注册时存储的安全答案是明文
2. 忘记密码时用MD5加密后的答案去验证，与明文不匹配
3. 个人中心修改的安全答案也是明文
4. 重置密码后无法用新密码登录（因为验证步骤就失败了）

## 修复步骤

### 1. 修改注册页面
- 文件：`web/src/views/RegisterView.vue`
- 在 `handleRegister` 函数中，将 `security_answer` 使用 MD5 加密：
  ```javascript
  security_answer: md5(form.security_answer)
  ```

### 2. 修改个人中心页面
- 文件：`web/src/views/ProfileView.vue`
- 在 `handleSecurityChange` 函数中，将 `security_answer` 使用 MD5 加密：
  ```javascript
  security_answer: md5(securityForm.value.answer)
  ```

### 3. 验证
- 确保所有三个地方（注册、个人中心修改、忘记密码验证）都统一使用 MD5 加密
- 编译构建并测试
