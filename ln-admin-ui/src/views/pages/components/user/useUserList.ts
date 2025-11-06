import { ref, computed, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import type { Ref } from 'vue'
import type { UserResponse, PageResponse } from '@/api/user'
import type { ActionButton, ExtraAction } from '../DataList.vue'
import { mapUserList } from './userUtils'

/**
 * 用户列表组合式函数
 * @param fetchUserListFn 获取用户列表的函数
 * @param options 配置选项
 */
export function useUserList(
    fetchUserListFn: (params: any) => Promise<PageResponse<UserResponse[]>>,
    options?: {
        onRefresh?: () => void
        onView?: (userId: string) => void
        onEdit?: (user: UserResponse) => void
        onDelete?: (userId: string) => Promise<void>
        deleteConfirmText?: string
        showViewAction?: boolean
        showEditAction?: boolean
        showDeleteAction?: boolean
    }
) {
    const {
        onRefresh,
        onView,
        onEdit,
        onDelete,
        deleteConfirmText = '确定要删除该用户吗？删除后不可恢复！',
        showViewAction = true,
        showEditAction = true,
        showDeleteAction = true,
    } = options || {}

    // 权限相关
    const permissionDenied = ref(false)

    // 用户CRUD相关状态
    const userFormOpen = ref(false)
    const userDetailOpen = ref(false)
    const currentUser = ref<UserResponse | null>(null)
    const currentUserId = ref<string | null>(null)

    // 数据列表引用（需要在父组件中传递）
    const dataListRef = ref<any>(null)

    /**
     * 处理获取用户列表
     */
    const fetchUserList = async (params: any): Promise<PageResponse<UserResponse[]>> => {
        try {
            const response = await fetchUserListFn(params)
            
            // 检查权限错误
            if (response.code === 403) {
                permissionDenied.value = true
                throw new Error(response.msg || response.message || '没有权限访问该资源')
            }
            
            // 权限验证通过，隐藏权限提示
            permissionDenied.value = false
            
            // 确保数据格式正确，处理可能的字段名差异
            if (response.data && Array.isArray(response.data)) {
                response.data = mapUserList(response.data)
            }
            
            return response
        } catch (error: any) {
            console.error('获取用户列表失败:', error)
            // 如果是权限错误，已经在上面设置了 permissionDenied
            if (error.message?.includes('权限') || error.message?.includes('没有权限')) {
                permissionDenied.value = true
            }
            throw error
        }
    }

    /**
     * 处理查看用户
     */
    const handleViewUser = async (record: UserResponse) => {
        if (onView) {
            onView(record.userId)
        } else {
            // 先设置用户ID，再打开弹窗，确保数据加载正确
            currentUserId.value = record.userId
            // 使用 nextTick 确保 userId 更新后再打开弹窗
            await nextTick()
            userDetailOpen.value = true
        }
    }

    /**
     * 处理编辑用户
     */
    const handleEditUser = async (record: UserResponse) => {
        if (onEdit) {
            onEdit(record)
        } else {
            // 先设置用户数据，再打开弹窗，确保数据加载正确
            currentUser.value = record
            // 使用 nextTick 确保 user 更新后再打开弹窗
            await nextTick()
            userFormOpen.value = true
        }
    }

    /**
     * 处理删除用户
     */
    const handleDeleteUser = async (record: UserResponse) => {
        try {
            if (onDelete) {
                await onDelete(record.userId)
            } else {
                // 如果没有提供删除函数，这里可以调用默认的API
                throw new Error('请提供删除用户的处理函数')
            }
            message.success('删除成功')
            refreshList()
        } catch (error: any) {
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
            throw error
        }
    }

    /**
     * 处理添加用户
     */
    const handleAddUser = () => {
        currentUser.value = null
        userFormOpen.value = true
    }

    /**
     * 刷新列表
     */
    const refreshList = () => {
        if (dataListRef.value?.refresh) {
            dataListRef.value.refresh()
        }
        if (onRefresh) {
            onRefresh()
        }
    }

    /**
     * 处理表单成功回调
     */
    const handleFormSuccess = () => {
        refreshList()
        userFormOpen.value = false
    }

    // 操作按钮配置
    const actions = computed<ActionButton[]>(() => {
        const actionButtons: ActionButton[] = []

        if (showViewAction) {
            actionButtons.push({
                key: 'view',
                label: '查看',
                type: 'link',
                onClick: handleViewUser,
            })
        }

        if (showEditAction) {
            actionButtons.push({
                key: 'edit',
                label: '编辑',
                type: 'link',
                onClick: handleEditUser,
            })
        }

        if (showDeleteAction) {
            actionButtons.push({
                key: 'delete',
                label: '删除',
                type: 'link',
                danger: true,
                confirm: deleteConfirmText,
                onClick: handleDeleteUser,
                disabled: (record: UserResponse) => {
                    // 超级管理员账号不能删除
                    return record.phone === '18797131041'
                },
            })
        }

        return actionButtons
    })

    return {
        // 状态
        permissionDenied,
        userFormOpen,
        userDetailOpen,
        currentUser,
        currentUserId,
        dataListRef,
        
        // 方法
        fetchUserList,
        handleViewUser,
        handleEditUser,
        handleDeleteUser,
        handleAddUser,
        refreshList,
        handleFormSuccess,
        
        // 配置
        actions,
    }
}

