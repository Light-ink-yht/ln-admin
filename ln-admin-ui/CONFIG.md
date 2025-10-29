# 配置说明

## 🎯 统一配置后端地址

后端地址**只需在 `vite.config.ts` 文件第 10 行修改一处即可**！

```typescript
// vite.config.ts
// ============================================
// 📌 后端服务配置 - 只在这里修改后端地址！
// ============================================
const BACKEND_URL = 'http://localhost:8080' // 👈 修改这里
```

## 📝 修改步骤

1. 打开 `vite.config.ts` 文件
2. 找到第 10 行的 `BACKEND_URL` 常量
3. 修改为你需要的后端地址
4. 保存文件
5. 重启开发服务器：`npm run dev`

## 💡 配置示例

### 本地开发
```typescript
const BACKEND_URL = 'http://localhost:8080'
```

### 局域网 IP
```typescript
const BACKEND_URL = 'http://192.168.1.100:8080'
```

### TL 服务器
```typescript
const BACKEND_URL = 'http://192.168.1.101:8080'
```

### 远程服务器
```typescript
const BACKEND_URL = 'http://api.example.com'
```

## ✅ 优点

- ✅ 只在一个地方配置，简单明了
- ✅ 避免多处配置不一致
- ✅ 默认值已经设置好，开箱即用
- ✅ 不需要创建 `.env` 文件

## 🚀 开始使用

1. 确保后端服务运行在 `http://localhost:8080`（或你的配置地址）
2. 运行 `npm run dev`
3. API 请求会自动代理到后端

