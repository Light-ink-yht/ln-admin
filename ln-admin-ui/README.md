# LN Admin UI

基于 Vue 3 + Vite + Ant Design Vue 4.0 构建的现代化后台管理系统模板。

## ✨ 特性

- 🎨 **动态主题系统** - 支持主题色切换、亮色/暗色模式、跟随系统主题
- 💾 **主题持久化** - 自动保存用户主题偏好
- 🎯 **TypeScript** - 完整的类型支持
- 🚀 **Vue 3** - 使用 Composition API
- 📦 **Vite** - 快速的热重载开发体验
- 🎁 **Ant Design Vue 4.0** - 丰富的组件库
- 🔧 **Axios 封装** - 完善的请求拦截和错误处理
- 🌐 **代理配置** - 开发环境 API 代理支持

## 📁 项目结构

```
ln-admin-ui/
├── src/
│   ├── api/               # API 接口定义
│   │   └── auth.ts       # 认证相关接口
│   ├── theme/            # 主题系统
│   │   ├── config.ts     # 主题配置
│   │   └── utils.ts      # 主题工具函数
│   ├── stores/           # Pinia 状态管理
│   │   └── modules/
│   │       └── theme.ts  # 主题 Store (组合式 API)
│   ├── utils/            # 工具函数
│   │   └── request.ts    # Axios 封装
│   ├── views/            # 页面组件
│   │   ├── login/        # 登录页面
│   │   │   └── Login.vue
│   │   └── ThemeDemo.vue # 主题演示页面
│   ├── router/           # 路由配置
│   ├── App.vue           # 根组件
│   └── env.config.ts     # 环境配置
└── docs/                 # 文档
    ├── THEME.md          # 主题系统文档
    └── LOGIN.md          # 登录系统文档
```

## 🚀 快速开始

### 环境要求

- Node.js >= 20.19.0 或 >= 22.12.0

### 安装依赖

```sh
npm install
```

### 开发

```sh
npm run dev
```

服务将在 `http://localhost:5173` 启动，并自动打开浏览器。

### 构建生产版本

```sh
npm run build
```

### 预览生产构建

```sh
npm run preview
```

## ⚙️ 环境配置

### 方式一：使用环境变量文件（推荐）

在项目根目录创建环境变量文件：

**`.env.development`** - 开发环境
```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_URL=http://localhost:8080
VITE_APP_TITLE=LN Admin UI - 开发环境
```

**`.env.production`** - 生产环境
```env
VITE_API_BASE_URL=http://your-production-api.com/api
VITE_APP_TITLE=LN Admin UI
```

**`.env.local`** - 本地配置（可选，不会被提交到 git）
```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_URL=http://localhost:8080
```

### 方式二：直接在 vite.config.ts 中修改

```typescript
// 在 vite.config.ts 中修改 BACKEND_URL 常量
const BACKEND_URL = 'http://localhost:8080'
```

### 默认配置

| 配置项 | 默认值 |
|--------|--------|
| 后端地址 | `http://localhost:8080` |
| API 路径 | `/api` |
| 开发端口 | `5173` |

## 🌐 API 代理配置

开发环境下，所有 `/api` 开头的请求会被代理到后端服务器。

### 配置示例

```typescript
// vite.config.ts
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 后端地址
      changeOrigin: true,
      rewrite: (path) => path,
    },
  },
}
```

### 使用环境配置

代码中通过 `env.config.ts` 获取配置：

```typescript
import { envConfig } from '@/env.config'

console.log(envConfig.apiBaseUrl)  // /api
console.log(envConfig.apiProxyUrl) // http://localhost:8080
```

## 🎨 主题系统

本项目实现了完整的动态主题系统，包含以下功能：

### 核心功能

- ✅ **动态主题色切换** - 8 种预设颜色 + 自定义颜色
- ✅ **亮色/暗色模式** - 一键切换主题模式
- ✅ **跟随系统主题** - 自动检测系统主题偏好
- ✅ **主题持久化** - 刷新页面后保留用户设置
- ✅ **Vue3 组合式 API** - 使用最新的 Composition API
- ✅ **TypeScript 支持** - 完整的类型定义

### 快速使用

```vue
<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme'
import { ThemeMode } from '@/theme/config'

const themeStore = useThemeStore()

// 切换主题色
themeStore.setColorPrimary('#1890ff')

// 切换主题模式
themeStore.setThemeMode(ThemeMode.DARK)

// 切换模式
themeStore.toggleThemeMode()

// 跟随系统主题
themeStore.followSystemTheme()
</script>
```

### 演示页面

启动项目后访问 `/theme-demo` 查看完整的主题功能演示。

### 详细文档

查看 [主题系统使用文档](./docs/THEME.md) 了解完整的 API 和使用示例。

## 📚 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - JavaScript 的超集
- **Vite** - 下一代前端构建工具
- **Vue Router** - 官方路由管理器
- **Pinia** - Vue 官方状态管理
- **Ant Design Vue** - 企业级 UI 组件库
- **Pinia Persistedstate** - 状态持久化插件
- **Axios** - HTTP 客户端

## 🔐 后端接口格式

本项目已适配常见后端响应格式：

```typescript
// 成功响应
{
  "code": 0,  // 或 200
  "msg": "操作成功",  // 或 message
  "data": { ... }
}

// 登录成功响应
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "jwt_token",
    "user": {
      "ID": 1,
      "userId": "767406100189184",
      "phone": "18797131041",
      "nickname": "用户名",
      "avatar": "/uploads/avatar/logo.png",
      ...
    }
  }
}
```

## 🛠️ 开发工具推荐

### IDE

推荐使用 [VS Code](https://code.visualstudio.com/) + [Vue Language Features (Volar)](https://marketplace.visualstudio.com/items?itemName=Vue.volar)。

### 浏览器扩展

**Chromium 浏览器 (Chrome, Edge, Brave 等):**
- [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)

**Firefox:**
- [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)

## 📖 文档

- [主题系统完整文档](./docs/THEME.md) - 详细的 API 文档和使用示例
- [登录系统文档](./docs/LOGIN.md) - 登录表单和接口说明
- [Ant Design Vue 文档](https://www.antdv.com/docs/vue/introduce-cn)
- [Vue 3 文档](https://cn.vuejs.org/)
- [Vite 文档](https://cn.vitejs.dev/)

## 📝 常见问题

### Q: 如何修改后端 API 地址？

A: 创建 `.env.development` 文件并设置：
```env
VITE_API_PROXY_URL=http://your-backend-server:8080
```

### Q: 开发环境如何跳过代理直接请求？

A: 修改 `vite.config.ts` 中的 `baseURL` 或直接使用完整 URL。

### Q: 如何查看 API 请求日志？

A: 开发环境下，代理会在控制台输出请求日志。

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
