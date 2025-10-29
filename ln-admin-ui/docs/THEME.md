# 主题系统使用文档

## 📋 目录

- [简介](#简介)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [使用示例](#使用示例)
- [常见问题](#常见问题)

## 📖 简介

本项目基于 Ant Design Vue 4.0 实现了完整的动态主题系统，支持：

- ✅ 动态切换主题色（8种预设颜色 + 自定义颜色）
- ✅ 亮色/暗色模式切换
- ✅ 跟随系统主题
- ✅ 主题持久化（刷新后保留）
- ✅ Vue3 组合式 API 实现
- ✅ TypeScript 完整支持

## 🚀 快速开始

### 1. 在应用中使用主题

在任何 Vue 组件中使用主题：

```vue
<template>
  <div>
    <a-button type="primary">主要按钮</a-button>
  </div>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
</script>
```

### 2. 切换主题色

```typescript
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()

// 使用预设颜色
themeStore.setColorPrimary('#1890ff')

// 使用自定义颜色
themeStore.setColorPrimary('#ff5722')
```

### 3. 切换主题模式

```typescript
import { useThemeStore } from '@/stores/modules/theme'
import { ThemeMode } from '@/theme/config'

const themeStore = useThemeStore()

// 设置为暗色模式
themeStore.setThemeMode(ThemeMode.DARK)

// 设置为亮色模式
themeStore.setThemeMode(ThemeMode.LIGHT)

// 切换模式
themeStore.toggleThemeMode()

// 跟随系统主题
themeStore.followSystemTheme()
```

## ⚙️ 配置说明

### 预设主题色

系统预设了 8 种主题色，定义在 `src/theme/config.ts`：

```typescript
export const PRESET_COLORS = [
  '#1890ff', // 蓝色（默认）
  '#f5222d', // 红色
  '#52c41a', // 绿色
  '#faad14', // 橙色
  '#13c2c2', // 青色
  '#722ed1', // 紫色
  '#eb2f96', // 粉色
  '#fa8c16', // 橙黄
]
```

### 主题模式

```typescript
export enum ThemeMode {
  LIGHT = 'light', // 亮色模式
  DARK = 'dark',   // 暗色模式
}
```

### 修改预设配置

编辑 `src/theme/config.ts` 文件：

```typescript
// 添加新的预设颜色
export const PRESET_COLORS = [
  '#1890ff',
  '#f5222d',
  // 添加你的颜色
  '#your-color',
]

// 添加颜色名称映射
export const COLOR_NAMES: Record<string, string> = {
  '#1890ff': '蓝色',
  '#f5222d': '红色',
  '#your-color': '你的颜色',
}

// 修改默认主题
export const DEFAULT_THEME = {
  colorPrimary: '#1890ff', // 默认主题色
  mode: ThemeMode.LIGHT,   // 默认模式
}
```

## 📚 API 文档

### ThemeStore

#### 状态 (State)

| 属性 | 类型 | 说明 |
|------|------|------|
| `colorPrimary` | `string` | 当前主题色 |
| `mode` | `ThemeMode \| string` | 当前主题模式 |
| `antdConfigProvider` | `any` | Ant Design ConfigProvider 实例 |

#### 计算属性 (Getters)

| 属性 | 类型 | 说明 |
|------|------|------|
| `themeConfig` | `ThemeConfig` | 完整的主题配置对象 |
| `isDark` | `boolean` | 是否为暗色模式 |

#### 方法 (Actions)

| 方法 | 参数 | 说明 |
|------|------|------|
| `setColorPrimary` | `color: string` | 设置主题色 |
| `setThemeMode` | `mode: ThemeMode` | 设置主题模式 |
| `toggleThemeMode` | - | 切换主题模式（亮色↔暗色） |
| `followSystemTheme` | - | 跟随系统主题设置 |
| `applyTheme` | - | 应用主题配置（通常不需要手动调用） |
| `initTheme` | - | 初始化主题（在应用启动时调用） |

### 类型定义

```typescript
// 主题配置
interface ThemeConfig {
  colorPrimary: string
  mode: ThemeMode
}

// 主题模式枚举
enum ThemeMode {
  LIGHT = 'light',
  DARK = 'dark'
}
```

## 💡 使用示例

### 示例 1: 主题色选择器

```vue
<template>
  <a-card title="选择主题色">
    <a-space wrap>
      <div
        v-for="color in PRESET_COLORS"
        :key="color"
        class="color-item"
        :style="{ backgroundColor: color }"
        @click="handleColorChange(color)"
      >
        <CheckOutlined v-if="themeStore.colorPrimary === color" />
        <span>{{ COLOR_NAMES[color] }}</span>
      </div>
    </a-space>
  </a-card>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme'
import { PRESET_COLORS, COLOR_NAMES } from '@/theme/config'
import { CheckOutlined } from '@ant-design/icons-vue'

const themeStore = useThemeStore()

const handleColorChange = (color: string) => {
  themeStore.setColorPrimary(color)
}
</script>
```

### 示例 2: 主题模式切换器

```vue
<template>
  <a-space>
    <a-switch
      :checked="themeStore.isDark"
      @change="handleModeToggle"
      checked-children="暗色"
      un-checked-children="亮色"
    />
    <a-button @click="followSystem">跟随系统</a-button>
  </a-space>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme'
import { ThemeMode } from '@/theme/config'

const themeStore = useThemeStore()

const handleModeToggle = (checked: boolean) => {
  themeStore.setThemeMode(checked ? ThemeMode.DARK : ThemeMode.LIGHT)
}

const followSystem = () => {
  themeStore.followSystemTheme()
}
</script>
```

### 示例 3: 自定义颜色选择器

```vue
<template>
  <a-space>
    <a-color-picker v-model:value="customColor" show-text />
    <a-button type="primary" @click="applyCustomColor">
      应用自定义颜色
    </a-button>
  </a-space>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
const customColor = ref('#1890ff')

const applyCustomColor = () => {
  themeStore.setColorPrimary(customColor.value)
}
</script>
```

### 示例 4: 根据主题模式调整内容

```vue
<template>
  <div :class="{ 'dark-content': themeStore.isDark }">
    <h1>内容标题</h1>
    <p>这是根据主题自适应颜色的内容</p>
  </div>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
</script>

<style scoped>
.dark-content {
  background-color: #141414;
  color: rgba(255, 255, 255, 0.85);
}

.dark-content h1 {
  color: #fff;
}
</style>
```

## 🎨 主题定制

### CSS 变量

主题系统会自动设置以下 CSS 变量到 `:root`：

```css
:root {
  --color-primary: #1890ff;
  --color-primary-hover: #40a9ff;
  --color-primary-active: #096dd9;
}
```

### 使用 CSS 变量

在样式文件中使用：

```css
.my-button {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}

.my-button:hover {
  background-color: var(--color-primary-hover);
}

.my-button:active {
  background-color: var(--color-primary-active);
}
```

### 暗色模式样式

当切换到暗色模式时，`body` 元素会自动添加 `dark-mode` class：

```css
/* 亮色模式默认样式 */
.my-card {
  background-color: #fff;
  color: #333;
}

/* 暗色模式样式 */
body.dark-mode .my-card {
  background-color: #1f1f1f;
  color: rgba(255, 255, 255, 0.85);
}
```

### 自定义 Ant Design 组件样式

在 `src/theme/utils.ts` 中修改 `getAntdThemeConfig` 函数：

```typescript
export function getAntdThemeConfig(config: ThemeConfig): any {
  const isDark = config.mode === ThemeMode.DARK

  return {
    algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: config.colorPrimary,
      borderRadius: 6, // 修改全局圆角
      // 添加更多 token 配置
    },
    components: {
      Button: {
        borderRadius: 6, // 按钮圆角
      },
      Card: {
        borderRadius: 8, // 卡片圆角
      },
      // 添加更多组件配置
    },
  }
}
```

## 🔧 持久化配置

主题配置默认使用 `localStorage` 持久化，key 为 `theme-store`。

修改持久化配置（在 `src/stores/modules/theme.ts` 中）：

```typescript
{
  persist: {
    key: 'your-custom-key', // 自定义 key
    storage: sessionStorage, // 使用 sessionStorage 代替 localStorage
  },
}
```

## ❓ 常见问题

### Q1: 如何在页面初始化时应用主题？

在应用的入口文件（如 `src/App.vue`）中调用：

```typescript
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()

onMounted(() => {
  themeStore.initTheme()
})
```

### Q2: 主题切换不生效？

确保在 `ConfigProvider` 中正确绑定主题：

```vue
<template>
  <ConfigProvider :theme="antdTheme">
    <!-- 你的内容 -->
  </ConfigProvider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ConfigProvider } from 'ant-design-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { getAntdThemeConfig } from '@/theme/utils'

const themeStore = useThemeStore()

const antdTheme = computed(() => {
  return getAntdThemeConfig(themeStore.themeConfig)
})
</script>
```

### Q3: 如何添加新的预设颜色？

编辑 `src/theme/config.ts`：

```typescript
export const PRESET_COLORS = [
  // ... 现有颜色
  '#your-new-color',
]

export const COLOR_NAMES: Record<string, string> = {
  // ... 现有映射
  '#your-new-color': '新颜色名称',
}
```

### Q4: 支持哪些颜色格式？

目前支持十六进制颜色格式（`#RRGGBB`），例如：
- `#1890ff`
- `#f5222d`

### Q5: 如何禁用主题持久化？

在 `src/stores/modules/theme.ts` 中移除或注释 `persist` 配置：

```typescript
{
  // persist: {
  //   key: 'theme-store',
  //   storage: localStorage,
  // },
}
```

## 📝 更新日志

### v1.0.0
- ✨ 初始版本
- ✅ 支持动态主题色切换
- ✅ 支持亮色/暗色模式
- ✅ 支持主题持久化
- ✅ Vue3 组合式 API 实现
- ✅ TypeScript 支持

## 📄 许可证

MIT License

