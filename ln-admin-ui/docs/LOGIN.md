# 登录系统使用文档

## 📋 概述

本文档介绍 LN Admin UI 登录系统的使用方法，包括三个核心功能表单：登录、注册、忘记密码。

## 🚀 功能特性

- ✅ **三种表单切换** - 登录、注册、忘记密码无缝切换
- ✅ **双验证码系统** - 图片验证码 + 短信验证码
- ✅ **表单验证** - 完整的前端表单验证
- ✅ **主题适配** - 自动适配当前主题色
- ✅ **响应式设计** - 完美支持移动端和桌面端
- ✅ **暗色模式支持** - 完美适配亮色和暗色主题

## 📝 表单字段说明

### 1. 登录表单

| 字段 | 类型 | 说明 | 验证规则 |
|------|------|------|---------|
| phone | string | 手机号 | 必填，11位手机号格式 |
| password | string | 密码 | 必填，6-20位 |
| code | string | 图片验证码 | 必填 |
| captchaId | string | 验证码ID | 自动获取 |

### 2. 注册表单

| 字段 | 类型 | 说明 | 验证规则 |
|------|------|------|---------|
| phone | string | 手机号 | 必填，11位手机号格式 |
| code | string | 短信验证码 | 必填 |
| password | string | 密码 | 必填，6-20位 |
| confirmPassword | string | 确认密码 | 必填，需与密码一致 |

### 3. 忘记密码表单

| 字段 | 类型 | 说明 | 验证规则 |
|------|------|------|---------|
| phone | string | 手机号 | 必填，11位手机号格式 |
| code | string | 短信验证码 | 必填 |
| password | string | 新密码 | 必填，6-20位 |
| confirmPassword | string | 确认新密码 | 必填，需与新密码一致 |

## 🔧 API 接口

### 1. 获取图片验证码

```typescript
GET /api/auth/captcha

Response:
{
  code: 200,
  message: "success",
  data: {
    captchaId: "captcha_123456",
    image: "data:image/png;base64,..."
  }
}
```

### 2. 发送短信验证码

```typescript
POST /api/auth/send-sms

Request:
{
  phone: "13800138000",
  type: "register" | "forgot"
}

Response:
{
  code: 200,
  message: "验证码已发送"
}
```

### 3. 用户登录

```typescript
POST /api/auth/login

Request:
{
  phone: "13800138000",
  password: "password123",
  code: "ABC123",
  captchaId: "captcha_123456"
}

Response:
{
  code: 200,
  message: "登录成功",
  data: {
    token: "jwt_token_here",
    user: {
      id: "user_123",
      phone: "13800138000",
      nickname: "用户名"
    }
  }
}
```

### 4. 用户注册

```typescript
POST /api/auth/register

Request:
{
  phone: "13800138000",
  code: "123456",
  password: "password123",
  confirmPassword: "password123"
}

Response:
{
  code: 200,
  message: "注册成功",
  data: {
    id: "user_123",
    phone: "13800138000",
    nickname: "用户名"
  }
}
```

### 5. 忘记密码

```typescript
POST /api/auth/forgot-password

Request:
{
  phone: "13800138000",
  code: "123456",
  password: "newpassword123",
  confirmPassword: "newpassword123"
}

Response:
{
  code: 200,
  message: "密码重置成功"
}
```

## 💻 使用示例

### 基础使用

访问 `/login` 路由即可使用登录页面：

```vue
<template>
  <router-view />
</template>
```

### 自定义验证码刷新

```typescript
import { authApi } from '@/api/auth'

const refreshCaptcha = async () => {
  const res = await authApi.getCaptcha()
  if (res.code === 200) {
    captchaImage.value = res.data.image
    captchaId.value = res.data.captchaId
  }
}
```

### 发送短信验证码

```typescript
import { authApi } from '@/api/auth'

const sendSmsCode = async () => {
  const res = await authApi.sendSms({
    phone: '13800138000',
    type: 'register'
  })
  if (res.code === 200) {
    message.success('验证码已发送')
  }
}
```

## 🎨 主题适配

登录页面会自动使用当前主题色：

```less
// 按钮使用主题色
.login-button {
  background: var(--color-primary);
}
```

### 暗色模式

页面已内置暗色模式样式，切换主题时会自动适配：

```less
body.dark-mode {
  .login-container {
    background: linear-gradient(135deg, #1f1f1f 0%, #000 100%);
  }
}
```

## 📱 响应式设计

### 桌面端（> 768px）

- 左右分栏布局
- 左侧品牌展示
- 右侧表单区域

### 移动端（≤ 768px）

- 单栏布局
- 隐藏左侧品牌区
- 优化触摸体验

## 🔐 安全特性

1. **图片验证码** - 防止暴力破解
2. **短信验证码** - 60秒防刷倒计时
3. **密码强度校验** - 6-20位密码要求
4. **Token 存储** - localStorage 安全存储
5. **自动跳转** - 登录成功自动跳转

## 📦 依赖说明

```json
{
  "dependencies": {
    "axios": "^1.x.x",  // HTTP 请求
    "ant-design-vue": "^4.2.6",  // UI 组件库
    "vue-router": "^4.x.x",  // 路由管理
    "pinia": "^3.x.x"  // 状态管理
  }
}
```

## ⚙️ 配置

### 环境变量

在 `.env` 文件中配置 API 基础地址：

```env
VITE_API_BASE_URL=/api
```

### 修改 API 接口

编辑 `src/api/auth.ts` 文件：

```typescript
export const authApi = {
  login(params: LoginParams) {
    return request.post('/auth/login', params)
  },
  // 修改其他接口...
}
```

## 🐛 常见问题

### Q1: 验证码不显示？

确保后端返回的图片格式正确：

```typescript
{
  image: "data:image/png;base64,..."
}
```

### Q2: 短信验证码没有发送？

检查手机号格式是否正确：

```typescript
/^1[3-9]\d{9}$/  // 支持13-19开头的手机号
```

### Q3: Token 未保存？

确保登录接口返回的数据格式正确：

```typescript
{
  data: {
    token: "jwt_token"
  }
}
```

### Q4: 页面样式错乱？

确保已正确引入主题样式：

```typescript
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
themeStore.initTheme()
```

## 📄 许可证

MIT License

