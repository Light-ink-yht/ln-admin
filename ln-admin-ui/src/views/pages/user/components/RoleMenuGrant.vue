<template>
    <a-modal
        v-model:open="visible"
        title="角色菜单授权"
        :width="800"
        :confirm-loading="loading"
        @ok="handleSubmit"
        @cancel="handleCancel"
    >
        <a-spin :spinning="menusLoading">
            <div class="grant-menu-content">
                <a-alert
                    message="提示"
                    description="为角色分配菜单后，只有拥有该角色的用户才能看到相应的菜单。超级管理员可以看到所有菜单。"
                    type="info"
                    show-icon
                    style="margin-bottom: 16px"
                />

                <a-input
                    v-model:value="searchKeyword"
                    placeholder="搜索菜单（标题、路径）"
                    allow-clear
                    style="margin-bottom: 16px"
                >
                    <template #prefix>
                        <SearchOutlined />
                    </template>
                </a-input>

                <a-tree
                    v-if="filteredMenuTree && filteredMenuTree.length > 0"
                    v-model:checkedKeys="selectedMenuKeys"
                    checkable
                    :tree-data="filteredMenuTree"
                    :field-names="{ title: 'title', key: 'menu_id', children: 'children' }"
                    :expanded-keys="expandedKeys"
                    @expand="onExpand"
                    :check-strictly="false"
                    style="max-height: 500px; overflow-y: auto;"
                />
                <a-empty v-else description="暂无菜单" :image="Empty.PRESENTED_IMAGE_SIMPLE" />

                <div v-if="selectedMenuKeys.length > 0" class="selected-info">
                    <a-divider>已选择 {{ selectedMenuKeys.length }} 个菜单</a-divider>
                </div>
            </div>
        </a-spin>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { message, Empty } from 'ant-design-vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import { menuApi, type Menu } from '@/api/menu'

interface Props {
    open: boolean
    roleId?: string | null
}

const props = withDefaults(defineProps<Props>(), {
    open: false,
    roleId: null,
})

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

const loading = ref(false)
const menusLoading = ref(false)
const menuTree = ref<Menu[]>([])
const selectedMenuKeys = ref<string[]>([])
const expandedKeys = ref<string[]>([])
const searchKeyword = ref('')

// 过滤菜单树
const filteredMenuTree = computed(() => {
    if (!searchKeyword.value) {
        return menuTree.value
    }
    const keyword = searchKeyword.value.toLowerCase()
    return filterMenuTree(menuTree.value, keyword)
})

// 过滤菜单树函数
const filterMenuTree = (menus: Menu[], keyword: string): Menu[] => {
    const result: Menu[] = []
    for (const menu of menus) {
        const matches = 
            menu.title?.toLowerCase().includes(keyword) ||
            menu.path?.toLowerCase().includes(keyword)
        
        const filteredChildren = menu.children ? filterMenuTree(menu.children, keyword) : []
        
        if (matches || filteredChildren.length > 0) {
            result.push({
                ...menu,
                children: filteredChildren.length > 0 ? filteredChildren : menu.children,
            })
            // 如果有匹配的子菜单，展开父菜单
            if (filteredChildren.length > 0 && !expandedKeys.value.includes(menu.menu_id)) {
                expandedKeys.value.push(menu.menu_id)
            }
        }
    }
    return result
}

// 展开/收起节点
const onExpand = (keys: string[]) => {
    expandedKeys.value = keys
}

// 加载菜单树
const loadMenuTree = async () => {
    menusLoading.value = true
    try {
        const response = await menuApi.getMenuTree('1') // 只加载侧边栏菜单
        if (response.code === 200 || response.code === 0) {
            menuTree.value = response.data || []
            // 默认展开所有节点
            expandedKeys.value = getAllMenuIds(menuTree.value)
        }
    } catch (error: any) {
        console.error('加载菜单树失败:', error)
        message.error('加载菜单树失败')
    } finally {
        menusLoading.value = false
    }
}

// 获取所有菜单ID（用于默认展开）
const getAllMenuIds = (menus: Menu[]): string[] => {
    const ids: string[] = []
    for (const menu of menus) {
        ids.push(menu.menu_id)
        if (menu.children && menu.children.length > 0) {
            ids.push(...getAllMenuIds(menu.children))
        }
    }
    return ids
}

// 加载角色的菜单列表
const loadRoleMenus = async () => {
    if (!props.roleId) {
        selectedMenuKeys.value = []
        return
    }
    try {
        const response = await menuApi.getRoleMenus(props.roleId)
        if (response.code === 200 || response.code === 0) {
            selectedMenuKeys.value = response.data?.menu_ids || []
        }
    } catch (error: any) {
        console.error('加载角色菜单失败:', error)
        message.error('加载角色菜单失败')
    }
}

// 提交授权
const handleSubmit = async () => {
    if (!props.roleId) {
        message.error('角色ID不能为空')
        return
    }

    loading.value = true
    try {
        await menuApi.grantRoleMenus(props.roleId, selectedMenuKeys.value)
        message.success('授权成功')
        emit('success')
        handleCancel()
    } catch (error: any) {
        console.error('授权失败:', error)
        message.error(error.message || '授权失败')
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = () => {
    selectedMenuKeys.value = []
    searchKeyword.value = ''
    emit('update:open', false)
}

// 监听打开状态
watch(
    () => props.open,
    (val) => {
        if (val) {
            loadMenuTree()
            loadRoleMenus()
        }
    }
)
</script>

<style scoped lang="less">
.grant-menu-content {
    .selected-info {
        margin-top: 16px;
    }
}

:deep(.ant-tree) {
    .ant-tree-node-content-wrapper {
        &:hover {
            background-color: rgba(24, 144, 255, 0.06);
        }
    }

    .ant-tree-checkbox-checked .ant-tree-checkbox-inner {
        background-color: var(--color-primary);
        border-color: var(--color-primary);
    }
}

/* 深色模式 */
body.dark-mode {
    :deep(.ant-tree) {
        background: #141414;
        color: rgba(255, 255, 255, 0.85);

        .ant-tree-node-content-wrapper {
            &:hover {
                background-color: rgba(24, 144, 255, 0.1);
            }
        }
    }
}
</style>

