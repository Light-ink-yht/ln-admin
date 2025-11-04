# 用户相关组件和使用说明

本目录包含用户列表相关的可复用组件和工具函数。

## 组件

### UserAvatar
用户头像组件，自动显示头像或首字母。

**Props:**
- `avatar?: string` - 头像URL
- `nickname?: string` - 昵称
- `fullName?: string` - 姓名
- `size?: number` - 头像大小，默认44
- `defaultText?: string` - 默认显示文本，默认'U'

**示例:**
```vue
<UserAvatar
    :avatar="user.avatar"
    :nickname="user.nickname"
    :full-name="user.fullName"
    :size="44"
/>
```

### UserStatus
用户状态显示组件，使用Badge显示启用/禁用状态。

**Props:**
- `status?: string` - 状态值（'1'启用，'2'禁用）
- `enabledText?: string` - 启用文本，默认'启用'
- `disabledText?: string` - 禁用文本，默认'禁用'
- `unknownText?: string` - 未知状态文本，默认'未知'

**示例:**
```vue
<UserStatus :status="user.status" />
```

### UserGender
用户性别显示组件。

**Props:**
- `gender?: string` - 性别值（'1'男，'2'女，'3'未知）

**示例:**
```vue
<UserGender :gender="user.gender" />
```

### UserText
文本显示组件，自动处理空值显示为'-'。

**Props:**
- `value?: string | null` - 要显示的值
- `emptyText?: string` - 空值时的显示文本，默认'-'

**示例:**
```vue
<UserText :value="user.phone" />
```

### PermissionDeniedAlert
权限不足提示组件。

**Props:**
- `visible?: boolean` - 是否显示
- `description?: string` - 提示描述

**Events:**
- `update:visible` - 更新显示状态
- `close` - 关闭事件

**示例:**
```vue
<PermissionDeniedAlert
    v-model:visible="permissionDenied"
    description="您没有访问该资源的权限"
/>
```

## 工具函数

### getUserSearchFields
获取用户搜索字段配置。

**参数:**
- `customFields?: Partial<SearchField>[]` - 自定义字段配置

**返回:** `SearchField[]`

**示例:**
```typescript
// 使用默认配置
const searchFields = getUserSearchFields()

// 自定义配置
const searchFields = getUserSearchFields([
    { key: 'phone', width: '200px' }, // 覆盖phone字段
    { key: 'custom', label: '自定义', type: 'input' }, // 添加新字段
])
```

### getUserColumns
获取用户表格列配置。

**参数:**
- `options?: { ... }` - 配置选项
  - `showAvatar?: boolean` - 显示头像列，默认true
  - `showGender?: boolean` - 显示性别列，默认true
  - `showStatus?: boolean` - 显示状态列，默认true
  - `showLoginCount?: boolean` - 显示登录次数列，默认true
  - `showLastLoginTime?: boolean` - 显示最后登录时间列，默认true
  - `showCreatedAt?: boolean` - 显示创建时间列，默认true
  - `customColumns?: TableColumnsType` - 自定义列

**返回:** `TableColumnsType`

**示例:**
```typescript
// 使用默认配置
const columns = getUserColumns()

// 自定义配置
const columns = getUserColumns({
    showLoginCount: false, // 隐藏登录次数列
    showLastLoginTime: false, // 隐藏最后登录时间列
    customColumns: [
        { title: '自定义列', dataIndex: 'custom', key: 'custom' }
    ]
})
```

### useUserList
用户列表组合式函数，封装了用户列表的通用逻辑。

**参数:**
- `fetchUserListFn: (params: any) => Promise<PageResponse<UserResponse[]>>` - 获取用户列表的函数
- `options?: { ... }` - 配置选项
  - `onRefresh?: () => void` - 刷新回调
  - `onView?: (userId: string) => void` - 查看用户回调
  - `onEdit?: (user: UserResponse) => void` - 编辑用户回调
  - `onDelete?: (userId: string) => Promise<void>` - 删除用户回调
  - `deleteConfirmText?: string` - 删除确认文本
  - `showViewAction?: boolean` - 显示查看按钮，默认true
  - `showEditAction?: boolean` - 显示编辑按钮，默认true
  - `showDeleteAction?: boolean` - 显示删除按钮，默认true

**返回:**
```typescript
{
    permissionDenied: Ref<boolean>
    userFormOpen: Ref<boolean>
    userDetailOpen: Ref<boolean>
    currentUser: Ref<UserResponse | null>
    currentUserId: Ref<string | null>
    dataListRef: Ref<any>
    fetchUserList: (params: any) => Promise<PageResponse<UserResponse[]>>
    handleViewUser: (record: UserResponse) => void
    handleEditUser: (record: UserResponse) => void
    handleDeleteUser: (record: UserResponse) => Promise<void>
    handleAddUser: () => void
    refreshList: () => void
    handleFormSuccess: () => void
    actions: ComputedRef<ActionButton[]>
}
```

**示例:**
```vue
<script setup>
import { useUserList, getUserSearchFields, getUserColumns } from '@/views/pages/components'
import { userApi } from '@/api/user'

const {
    permissionDenied,
    dataListRef,
    fetchUserList,
    actions,
    handleFormSuccess,
} = useUserList(
    async (params) => await userApi.getUserList(params),
    {
        onDelete: async (userId) => {
            await userApi.deleteUser(userId)
        },
    }
)

const searchFields = getUserSearchFields()
const columns = getUserColumns()
</script>
```

## 其他工具函数

### getGenderText
获取性别文本。

**参数:**
- `gender: string` - 性别值

**返回:** `string`

### mapUserData
映射单个用户数据，处理字段名转换。

**参数:**
- `user: any` - 原始用户数据

**返回:** `UserResponse`

### mapUserList
批量映射用户数据。

**参数:**
- `users: any[]` - 用户数据数组

**返回:** `UserResponse[]`

## 常量

### USER_STATUS_OPTIONS
用户状态选项数组。

### USER_GENDER_OPTIONS
用户性别选项数组。

