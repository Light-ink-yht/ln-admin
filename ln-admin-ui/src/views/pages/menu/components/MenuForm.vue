<template>
    <a-modal
        v-model:open="modalOpen"
        :title="isEdit ? '编辑菜单' : '添加菜单'"
        width="800px"
        :confirm-loading="loading"
        @ok="handleSubmit"
        @cancel="handleCancel"
    >
        <a-form
            ref="formRef"
            :model="formData"
            :rules="rules"
            :label-col="{ span: 6 }"
            :wrapper-col="{ span: 18 }"
        >
            <a-form-item label="菜单标识" name="menu_key" v-if="!isEdit">
                <a-input
                    v-model:value="formData.menu_key"
                    placeholder="请输入菜单标识，如：menu-list"
                />
            </a-form-item>
            <a-form-item label="菜单标题" name="title">
                <a-input
                    v-model:value="formData.title"
                    placeholder="请输入菜单标题"
                />
            </a-form-item>
            <a-form-item label="菜单类型" name="menu_type">
                <a-radio-group v-model:value="formData.menu_type">
                    <a-radio value="1">侧边栏菜单</a-radio>
                    <a-radio value="2">用户菜单</a-radio>
                </a-radio-group>
            </a-form-item>
            <a-form-item label="路由路径" name="path">
                <a-input
                    v-model:value="formData.path"
                    placeholder="请输入路由路径，如：/menu/list"
                />
            </a-form-item>
            <a-form-item label="图标" name="icon">
                <a-input
                    v-model:value="formData.icon"
                    placeholder="请输入图标名称，如：MenuOutlined"
                />
            </a-form-item>
            <a-form-item label="父菜单" name="parent_id">
                <a-tree-select
                    v-model:value="formData.parent_id"
                    :tree-data="parentMenuOptions"
                    placeholder="请选择父菜单（留空表示根菜单）"
                    allow-clear
                    tree-default-expand-all
                    :field-names="{ label: 'title', value: 'menu_id', children: 'children' }"
                />
            </a-form-item>
            <a-form-item label="排序" name="sort">
                <a-input-number
                    v-model:value="formData.sort"
                    :min="0"
                    placeholder="请输入排序号，数字越小越靠前"
                    style="width: 100%"
                />
            </a-form-item>
            <a-form-item label="权限标识" name="permission">
                <a-input
                    v-model:value="formData.permission"
                    placeholder="请输入权限标识，如：/menu/list:GET"
                />
            </a-form-item>
            <a-form-item label="状态" name="status" v-if="isEdit">
                <a-radio-group v-model:value="formData.status">
                    <a-radio value="1">启用</a-radio>
                    <a-radio value="2">禁用</a-radio>
                </a-radio-group>
            </a-form-item>
            <a-form-item label="菜单描述" name="description">
                <a-textarea
                    v-model:value="formData.description"
                    placeholder="请输入菜单描述"
                    :rows="3"
                />
            </a-form-item>
        </a-form>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { menuApi, type Menu, type CreateMenuRequest, type UpdateMenuRequest } from '@/api/menu'

const props = defineProps<{
    open: boolean
    menu?: Menu | null
}>()

const emit = defineEmits<{
    'update:open': [value: boolean]
    success: []
}>()

const modalOpen = ref(false)
const loading = ref(false)
const formRef = ref<FormInstance>()
const isEdit = ref(false)
const parentMenuOptions = ref<Menu[]>([])

const formData = reactive<CreateMenuRequest & { status?: string }>({
    menu_key: '',
    title: '',
    menu_type: '1',
    path: '',
    icon: '',
    parent_id: '',
    sort: 0,
    permission: '',
    description: '',
    status: '1',
})

const rules = {
    menu_key: [{ required: true, message: '请输入菜单标识', trigger: 'blur' }],
    title: [{ required: true, message: '请输入菜单标题', trigger: 'blur' }],
    menu_type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
}

// 加载父菜单选项
const loadParentMenus = async () => {
    try {
        const response = await menuApi.getMenuList({ menu_type: formData.menu_type })
        if (response.code === 200 || response.code === 0) {
            // 过滤掉当前编辑的菜单（避免选择自己作为父菜单）
            const menus = response.data?.list || []
            if (isEdit.value && props.menu) {
                parentMenuOptions.value = menus.filter((m: Menu) => m.menu_id !== props.menu?.menu_id)
            } else {
                parentMenuOptions.value = menus
            }
        }
    } catch (error: any) {
        console.error('加载父菜单失败:', error)
        parentMenuOptions.value = []
    }
}

// 监听 open 变化
watch(
    () => props.open,
    (newVal) => {
        modalOpen.value = newVal
        if (newVal) {
            initFormData()
        }
    },
    { immediate: true }
)

// 监听 menu_type 变化，重新加载父菜单
watch(
    () => formData.menu_type,
    () => {
        if (modalOpen.value) {
            loadParentMenus()
        }
    }
)

// 监听 modalOpen 变化，同步到父组件
watch(modalOpen, (newVal) => {
    emit('update:open', newVal)
})

// 初始化表单数据
const initFormData = () => {
    if (props.menu) {
        isEdit.value = true
        Object.assign(formData, {
            menu_key: props.menu.menu_key,
            title: props.menu.title,
            menu_type: props.menu.menu_type,
            path: props.menu.path || '',
            icon: props.menu.icon || '',
            parent_id: props.menu.parent_id || '',
            sort: props.menu.sort || 0,
            permission: props.menu.permission || '',
            description: props.menu.description || '',
            status: props.menu.status || '1',
        })
    } else {
        isEdit.value = false
        Object.assign(formData, {
            menu_key: '',
            title: '',
            menu_type: '1',
            path: '',
            icon: '',
            parent_id: '',
            sort: 0,
            permission: '',
            description: '',
            status: '1',
        })
    }
    formRef.value?.resetFields()
    if (modalOpen.value) {
        loadParentMenus()
    }
}

// 提交表单
const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        loading.value = true

        if (isEdit.value && props.menu) {
            const updateData: UpdateMenuRequest = {
                title: formData.title,
                path: formData.path,
                icon: formData.icon,
                parent_id: formData.parent_id || undefined,
                sort: formData.sort || 0,
                permission: formData.permission,
                status: formData.status,
                description: formData.description,
            }
            await menuApi.updateMenu(props.menu.menu_id, updateData)
            message.success('更新成功')
        } else {
            const createData: CreateMenuRequest = {
                menu_key: formData.menu_key,
                title: formData.title,
                menu_type: formData.menu_type,
                path: formData.path,
                icon: formData.icon,
                parent_id: formData.parent_id || undefined,
                sort: formData.sort || 0,
                permission: formData.permission,
                description: formData.description,
            }
            await menuApi.createMenu(createData)
            message.success('创建成功')
        }

        emit('success')
        handleCancel()
    } catch (error: any) {
        console.error('提交失败:', error)
        message.error(error.message || '提交失败')
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = () => {
    modalOpen.value = false
    formRef.value?.resetFields()
}

onMounted(() => {
    if (props.open) {
        loadParentMenus()
    }
})
</script>

