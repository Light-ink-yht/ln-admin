# 文件上传组件

通用文件上传组件库，支持图片上传、文件上传等多种场景。

## 组件列表

### 1. ImageUpload - 多图片上传组件

支持多张图片上传，带预览和删除功能。

#### 基本用法

```vue
<template>
  <ImageUpload
    v-model="images"
    :multiple="true"
    :max-count="9"
    :max-size="10"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ImageUpload } from '@/components/upload'
import type { ImageItem } from '@/components/upload'

const images = ref<ImageItem[]>([])

const handleUploadSuccess = (file, index) => {
  console.log('上传成功:', file, index)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 图片列表（v-model） | ImageItem[] | [] |
| multiple | 是否支持多选 | boolean | false |
| drag | 是否支持拖拽上传 | boolean | false |
| maxSize | 单个文件大小限制（MB） | number | 10 |
| maxCount | 最多上传数量 | number | 9 |
| accept | 接受的文件类型 | string | 'image/*' |
| placeholder | 上传按钮文字 | string | '上传图片' |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 图片列表更新 | ImageItem[] |
| change | 图片列表变化 | ImageItem[] |
| upload-success | 上传成功 | (file: FileResponse, index: number) |
| upload-error | 上传失败 | (error: Error, index: number) |

#### ImageItem 接口

```typescript
interface ImageItem {
  url?: string          // 图片URL
  file_id?: string      // 文件ID（上传后）
  name?: string         // 文件名
  uploading?: boolean   // 是否正在上传
  progress?: number     // 上传进度（0-100）
  error?: string        // 错误信息
  file?: File          // 原始文件对象
}
```

---

### 2. ImageUploadSingle - 单图片上传组件

专门用于单张图片上传，带预览和编辑功能。

#### 基本用法

```vue
<template>
  <ImageUploadSingle
    v-model="imageUrl"
    :max-size="10"
    placeholder="上传头像"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ImageUploadSingle } from '@/components/upload'
import type { FileResponse } from '@/api/file'

const imageUrl = ref('')

const handleUploadSuccess = (file: FileResponse) => {
  console.log('上传成功:', file)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 图片URL（v-model） | string | '' |
| maxSize | 文件大小限制（MB） | number | 10 |
| accept | 接受的文件类型 | string | 'image/*' |
| placeholder | 上传按钮文字 | string | '上传图片' |
| previewAlt | 预览图片alt属性 | string | '图片预览' |
| drag | 是否支持拖拽上传 | boolean | false |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 图片URL更新 | string |
| change | 图片URL变化 | string |
| upload-success | 上传成功 | (file: FileResponse) |
| upload-error | 上传失败 | (error: Error) |

---

### 3. FileUpload - 通用文件上传组件

支持各种类型文件的上传，包括图片、文档、视频等。

#### 基本用法

```vue
<template>
  <FileUpload
    v-model="files"
    :multiple="true"
    :max-size="100"
    :max-count="10"
    button-text="上传文件"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { FileUpload } from '@/components/upload'
import type { UploadFile } from 'ant-design-vue'
import type { FileResponse } from '@/api/file'

const files = ref<UploadFile[]>([])

const handleUploadSuccess = (file: FileResponse, uploadFile: UploadFile) => {
  console.log('上传成功:', file, uploadFile)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 文件列表（v-model） | UploadFile[] | [] |
| multiple | 是否支持多选 | boolean | true |
| drag | 是否支持拖拽上传 | boolean | false |
| accept | 接受的文件类型 | string | '*' |
| maxSize | 单个文件大小限制（MB） | number | 100 |
| maxCount | 最多上传数量 | number | 10 |
| buttonText | 上传按钮文字 | string | '上传文件' |
| placeholder | 拖拽区域提示文字 | string | '点击或拖拽文件到此区域上传' |
| hint | 拖拽区域说明文字 | string | '' |
| showUploadList | 是否显示上传列表 | boolean | true |
| showFileList | 是否显示文件列表 | boolean | true |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 文件列表更新 | UploadFile[] |
| change | 文件列表变化 | UploadFile[] |
| upload-success | 上传成功 | (file: FileResponse, uploadFile: UploadFile) |
| upload-error | 上传失败 | (error: Error, uploadFile: UploadFile) |

---

### 4. VideoUpload - 视频上传组件

专门用于视频文件上传，支持视频预览和播放。

#### 基本用法

```vue
<template>
  <VideoUpload
    v-model="videos"
    :multiple="true"
    :max-count="5"
    :max-size="500"
    placeholder="上传视频"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { VideoUpload } from '@/components/upload'
import type { VideoItem } from '@/components/upload'

const videos = ref<VideoItem[]>([])

const handleUploadSuccess = (file, index) => {
  console.log('上传成功:', file, index)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 视频列表（v-model） | VideoItem[] | [] |
| multiple | 是否支持多选 | boolean | false |
| drag | 是否支持拖拽上传 | boolean | false |
| maxSize | 单个文件大小限制（MB） | number | 500 |
| maxCount | 最多上传数量 | number | 5 |
| accept | 接受的文件类型 | string | 'video/*' |
| placeholder | 上传按钮文字 | string | '上传视频' |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 视频列表更新 | VideoItem[] |
| change | 视频列表变化 | VideoItem[] |
| upload-success | 上传成功 | (file: FileResponse, index: number) |
| upload-error | 上传失败 | (error: Error, index: number) |

#### VideoItem 接口

```typescript
interface VideoItem {
  url?: string          // 视频URL
  file_id?: string      // 文件ID（上传后）
  name?: string         // 文件名
  size?: number         // 文件大小（字节）
  duration?: number     // 视频时长（秒）
  uploading?: boolean   // 是否正在上传
  progress?: number     // 上传进度（0-100）
  error?: string        // 错误信息
  file?: File          // 原始文件对象
}
```

---

### 5. AudioUpload - 音频上传组件

专门用于音频文件上传，支持音频播放预览。

#### 基本用法

```vue
<template>
  <AudioUpload
    v-model="audios"
    :multiple="true"
    :max-count="10"
    :max-size="100"
    placeholder="上传音频"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { AudioUpload } from '@/components/upload'
import type { AudioItem } from '@/components/upload'

const audios = ref<AudioItem[]>([])

const handleUploadSuccess = (file, index) => {
  console.log('上传成功:', file, index)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 音频列表（v-model） | AudioItem[] | [] |
| multiple | 是否支持多选 | boolean | false |
| drag | 是否支持拖拽上传 | boolean | false |
| maxSize | 单个文件大小限制（MB） | number | 100 |
| maxCount | 最多上传数量 | number | 10 |
| accept | 接受的文件类型 | string | 'audio/*' |
| placeholder | 上传按钮文字 | string | '上传音频' |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 音频列表更新 | AudioItem[] |
| change | 音频列表变化 | AudioItem[] |
| upload-success | 上传成功 | (file: FileResponse, index: number) |
| upload-error | 上传失败 | (error: Error, index: number) |

#### AudioItem 接口

```typescript
interface AudioItem {
  url?: string          // 音频URL
  file_id?: string      // 文件ID（上传后）
  name?: string         // 文件名
  size?: number         // 文件大小（字节）
  duration?: number     // 音频时长（秒）
  uploading?: boolean   // 是否正在上传
  progress?: number     // 上传进度（0-100）
  error?: string        // 错误信息
  file?: File          // 原始文件对象
}
```

---

### 6. ArchiveUpload - 压缩包上传组件

专门用于压缩包文件上传，支持 ZIP、RAR、7Z 等格式。

#### 基本用法

```vue
<template>
  <ArchiveUpload
    v-model="archives"
    :multiple="true"
    :max-count="10"
    :max-size="500"
    placeholder="上传压缩包"
    @upload-success="handleUploadSuccess"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ArchiveUpload } from '@/components/upload'
import type { ArchiveItem } from '@/components/upload'

const archives = ref<ArchiveItem[]>([])

const handleUploadSuccess = (file, index) => {
  console.log('上传成功:', file, index)
}
</script>
```

#### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 压缩包列表（v-model） | ArchiveItem[] | [] |
| multiple | 是否支持多选 | boolean | false |
| drag | 是否支持拖拽上传 | boolean | false |
| maxSize | 单个文件大小限制（MB） | number | 500 |
| maxCount | 最多上传数量 | number | 10 |
| accept | 接受的文件类型 | string | '.zip,.rar,.7z,.tar,.gz' |
| placeholder | 上传按钮文字 | string | '上传压缩包' |
| autoUpload | 是否自动上传 | boolean | true |

#### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 压缩包列表更新 | ArchiveItem[] |
| change | 压缩包列表变化 | ArchiveItem[] |
| upload-success | 上传成功 | (file: FileResponse, index: number) |
| upload-error | 上传失败 | (error: Error, index: number) |

#### ArchiveItem 接口

```typescript
interface ArchiveItem {
  url?: string          // 压缩包URL
  file_id?: string      // 文件ID（上传后）
  name?: string         // 文件名
  size?: number         // 文件大小（字节）
  type?: string         // 压缩包类型（zip、rar等）
  uploading?: boolean   // 是否正在上传
  progress?: number     // 上传进度（0-100）
  error?: string        // 错误信息
  file?: File          // 原始文件对象
}
```

---

## 使用示例

### 示例1：用户头像上传

```vue
<template>
  <ImageUploadSingle
    v-model="avatarUrl"
    :max-size="5"
    placeholder="上传头像"
    preview-alt="用户头像"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ImageUploadSingle } from '@/components/upload'

const avatarUrl = ref('')
</script>
```

### 示例2：商品图片上传（多图）

```vue
<template>
  <ImageUpload
    v-model="productImages"
    :multiple="true"
    :max-count="9"
    :max-size="5"
    placeholder="上传商品图片"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ImageUpload } from '@/components/upload'
import type { ImageItem } from '@/components/upload'

const productImages = ref<ImageItem[]>([])
</script>
```

### 示例3：文档上传

```vue
<template>
  <FileUpload
    v-model="documents"
    :accept="'.pdf,.doc,.docx'"
    :max-size="50"
    :max-count="5"
    button-text="上传文档"
    placeholder="点击或拖拽文档到此区域上传"
    hint="支持 PDF、Word 格式，单个文件不超过 50MB"
  />
</template>

<script setup>
import { ref } from 'vue'
import { FileUpload } from '@/components/upload'
import type { UploadFile } from 'ant-design-vue'

const documents = ref<UploadFile[]>([])
</script>
```

### 示例4：拖拽上传

```vue
<template>
  <ImageUpload
    v-model="images"
    :drag="true"
    :multiple="true"
    placeholder="拖拽图片到此处上传"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ImageUpload } from '@/components/upload'
import type { ImageItem } from '@/components/upload'

const images = ref<ImageItem[]>([])
</script>
```

### 示例5：视频上传

```vue
<template>
  <VideoUpload
    v-model="videos"
    :multiple="true"
    :max-count="3"
    :max-size="500"
    :drag="true"
    placeholder="上传视频"
  />
</template>

<script setup>
import { ref } from 'vue'
import { VideoUpload } from '@/components/upload'
import type { VideoItem } from '@/components/upload'

const videos = ref<VideoItem[]>([])
</script>
```

### 示例6：音频上传

```vue
<template>
  <AudioUpload
    v-model="audios"
    :multiple="true"
    :max-count="10"
    :max-size="100"
    placeholder="上传音频文件"
  />
</template>

<script setup>
import { ref } from 'vue'
import { AudioUpload } from '@/components/upload'
import type { AudioItem } from '@/components/upload'

const audios = ref<AudioItem[]>([])
</script>
```

### 示例7：压缩包上传

```vue
<template>
  <ArchiveUpload
    v-model="archives"
    :multiple="true"
    :max-count="5"
    :max-size="500"
    :drag="true"
    placeholder="上传压缩包"
  />
</template>

<script setup>
import { ref } from 'vue'
import { ArchiveUpload } from '@/components/upload'
import type { ArchiveItem } from '@/components/upload'

const archives = ref<ArchiveItem[]>([])
</script>
```

---

## 注意事项

1. **文件大小限制**：所有组件都支持 `maxSize` 参数限制文件大小，单位为 MB
2. **文件类型限制**：通过 `accept` 参数控制，可以使用 MIME 类型或文件扩展名
3. **自动上传**：默认 `autoUpload` 为 `true`，文件选择后自动上传；设置为 `false` 时，需要手动触发上传
4. **内存管理**：使用对象 URL 预览时，组件会自动释放内存，无需手动处理
5. **错误处理**：上传失败时，组件会显示错误信息，并触发 `upload-error` 事件

---

## 样式定制

所有组件都支持通过 CSS 变量进行主题定制：

```css
:root {
  --color-primary: #1890ff; /* 主色调 */
}
```

组件会自动使用主题色，无需额外配置。

