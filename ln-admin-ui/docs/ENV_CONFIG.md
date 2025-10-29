# 环境配置说明

## 📋 概述

本文档介绍如何配置 LN Admin UI 项目的环境变量和代理设置。

## ⚙️ 环境变量配置

### 方式一：使用环境变量文件（推荐）

在项目根目录创建以下环境变量文件：

#### 1. 开发环境配置 (`.env.development`)

```env
# API 基础地址
VITE_API_BASE_URL=/api

# 应用标题
VITE_APP_TITLE=LN Admin UI - 开发环境

# 后端服务地址（用于代理）
VITE_API_PROXY_URL=http://localhost:8080
```

#### 2. 生产环境配置 (`.env.production`)

```env
# API 基础地址（生产环境直接使用完整路径）
VITE_API_BASE_URL=http://your-production-api.com/api

# 应用标题
VITE_APP_TITLE=LN Admin UI

# 后端服务地址
VITE_API_PROXY_URL=http://your-production-api.com
```

#### 3. 本地配置 (`.env.local` - 可选)

本地开发时使用，不会被提交到 git：

```env
# API 基础地址
VITE_API_BASE_URL=/api

# 应用标题
VITE_APP_TITLE=LN Admin UI - 本地环境

# 后端服务地址（用于代理）
VITE_API_PROXY_URL=http://localhost:8080
```

### 方式二：在代码中直接配置

如果不想使用环境变量文件，可以直接在 `src/env.config.ts` 中修改默认值。

## 🌐 代理配置

### 开发环境代理

在开发环境下，`vite.config.ts` 配置了代理，所有 `/api` 开头的请求会被代理到后端服务器。

### 配置示例

```typescript
// vite.config.ts
server: {
  port: 5173,
  host: true,
  open: true,
  cors: true,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 后端地址
      changeOrigin: true,                // 改变源地址
      rewrite: (path) => path,           // 不重写路径
    },
  },
}
```

### 代理配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `target` | 后端服务器地址 | `http://localhost:8080` |
| `changeOrigin` | 改变请求头中的 origin | `true` |
| `rewrite` | 重写请求路径 | 不重写 |
| `port` | 开发服务器端口 | `5173` |
| `host` | 监听所有地址 | `true` |

### 修改后端地址

#### 方法一：使用环境变量

创建 `.env.development` 文件：

```env
VITE_API_PROXY_URL=http://192.168.1.100:8080
```

#### 方法二：直接修改 vite.config.ts

```typescript
proxy: {
  '/api': {
    target: 'http://192.168.1.100:8080',  // 修改这里
    changeOrigin: true,
    rewrite: (path) => path,
  },
}
```

## 🔧 环境变量说明

### 变量列表

| 变量名 | 说明 | 必需 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `VITE_API_BASE_URL` | API 基础地址 | 否 | `/api` | `/api` |
| `VITE_API_PROXY_URL` | 后端服务地址 | 否 | `http://localhost:8080` | `http://localhost:8080` |
| `VITE_APP_TITLE` | 应用标题 | 否 | `LN Admin UI` | `LN Admin UI` |

### 环境变量优先级

1. `.env.local` - 本地配置（最高优先级）
2. `.env.development` / `.env.production` - 环境配置
3. `env.config.ts` 中的默认值（最低优先级）

### 使用环境配置

在代码中使用 `env.config.ts` 获取配置：

```typescript
import { envConfig } from '@/env.config'

// 获取 API 基础地址
const apiUrl = envConfig.apiBaseUrl

// 获取代理地址
const proxyUrl = envConfig.apiProxyUrl

// 获取应用标题
const title = envConfig.appTitle

// 判断环境
const isDev = envConfig.dev
const isProd = envConfig.prod
```

## 📝 配置示例

### 示例 1: 本地开发

后端运行在 `http://localhost:8080`

创建 `.env.development`：
```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_URL=http://localhost:8080
```

### 示例 2: 局域网开发

后端运行在局域网 IP `http://192.168.1.100:8080`

创建 `.env.development`：
```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_URL=http://192.168.1.100:8080
```

### 示例 3: 远程服务器开发

后端运行在远程服务器 `http://dev.example.com`

创建 `.env.development`：
```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_URL=http://dev.example.com
```

### 示例 4: 生产环境

生产环境直接请求完整 API 地址

创建 `.env.production`：
```env
VITE_API_BASE_URL=https://api.example.com/api
VITE_API_PROXY_URL=https://api.example.com
```

## 🐛 常见问题

### Q1: 请求被 CORS 阻止？

**A:** 确保：
1. 代理配置中 `changeOrigin: true`
2. 后端已配置 CORS 允许前端域名

### Q2: 代理不生效？

**A:** 检查：
1. 是否重启了开发服务器
2. 请求路径是否以 `/api` 开头
3. `vite.config.ts` 中代理配置是否正确

### Q3: 如何查看代理日志？

**A:** 开发环境下，代理会在控制台输出请求日志：
```
Sending Request to the Target: GET /api/auth/captcha
Received Response from the Target: 200 /api/auth/captcha
```

### Q4: 生产环境是否需要配置代理？

**A:** 不需要。生产环境下，直接使用完整的 API 地址，不需要代理。

### Q5: 如何修改开发服务器端口？

**A:** 在 `vite.config.ts` 中修改：
```typescript
server: {
  port: 3000,  // 修改为 3000
}
```

## 📚 参考资料

- [Vite 环境变量](https://cn.vitejs.dev/guide/env-and-mode.html)
- [Vite 服务器配置](https://cn.vitejs.dev/config/server-options.html)
- [Http Proxy Middleware](https://github.com/chimurai/http-proxy-middleware)

## 📄 许可证

MIT License

