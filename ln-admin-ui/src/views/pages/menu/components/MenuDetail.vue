<template>
    <a-modal
        v-model:open="modalOpen"
        title="菜单详情"
        width="800px"
        :footer="null"
    >
        <a-descriptions :column="2" bordered v-if="menu">
            <a-descriptions-item label="菜单ID" :span="2">
                <code>{{ menu.menu_id }}</code>
            </a-descriptions-item>
            <a-descriptions-item label="菜单标识">
                <code>{{ menu.menu_key }}</code>
            </a-descriptions-item>
            <a-descriptions-item label="菜单标题">
                {{ menu.title }}
            </a-descriptions-item>
            <a-descriptions-item label="菜单类型">
                <a-tag :color="menu.menu_type === '1' ? 'blue' : 'green'">
                    {{ menu.menu_type === '1' ? '侧边栏菜单' : '用户菜单' }}
                </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
                <UserStatus :status="menu.status" />
            </a-descriptions-item>
            <a-descriptions-item label="路由路径">
                {{ menu.path || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="图标">
                {{ menu.icon || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="父菜单ID">
                {{ menu.parent_id || '无（根菜单）' }}
            </a-descriptions-item>
            <a-descriptions-item label="排序">
                {{ menu.sort }}
            </a-descriptions-item>
            <a-descriptions-item label="权限标识">
                {{ menu.permission || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="菜单描述" :span="2">
                {{ menu.description || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="创建人">
                {{ menu.creator_id || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="修改人">
                {{ menu.modifier_id || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">
                {{ menu.created_at }}
            </a-descriptions-item>
            <a-descriptions-item label="更新时间">
                {{ menu.updated_at }}
            </a-descriptions-item>
        </a-descriptions>
        <a-empty v-else description="暂无数据" />
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { UserStatus } from '@/views/pages/components'
import { menuApi, type Menu } from '@/api/menu'

const props = defineProps<{
    open: boolean
    menuId?: string | null
}>()

const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const modalOpen = ref(false)
const menu = ref<Menu | null>(null)

// 监听 open 变化
watch(
    () => props.open,
    (newVal) => {
        modalOpen.value = newVal
        if (newVal && props.menuId) {
            loadMenuDetail()
        }
    },
    { immediate: true }
)

// 监听 modalOpen 变化，同步到父组件
watch(modalOpen, (newVal) => {
    emit('update:open', newVal)
})

// 加载菜单详情
const loadMenuDetail = async () => {
    if (!props.menuId) return

    try {
        const response = await menuApi.getMenuDetail(props.menuId)
        if (response.code === 200 || response.code === 0) {
            menu.value = response.data
        }
    } catch (error: any) {
        console.error('加载菜单详情失败:', error)
        menu.value = null
    }
}
</script>

