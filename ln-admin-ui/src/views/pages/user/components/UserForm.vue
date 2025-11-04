<template>
    <a-modal
        v-model:open="visible"
        :title="isEdit ? '编辑用户' : '添加用户'"
        :width="800"
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
            <a-form-item label="手机号" name="phone" v-if="!isEdit">
                <a-input v-model:value="formData.phone" placeholder="请输入手机号" />
            </a-form-item>

            <a-form-item label="密码" name="password" v-if="!isEdit">
                <a-input-password v-model:value="formData.password" placeholder="请输入密码" />
            </a-form-item>

            <a-form-item label="邮箱" name="email">
                <a-input v-model:value="formData.email" placeholder="请输入邮箱" />
            </a-form-item>

            <a-form-item label="昵称" name="nickname">
                <a-input v-model:value="formData.nickname" placeholder="请输入昵称" />
            </a-form-item>

            <a-form-item label="姓名" name="fullName">
                <a-input v-model:value="formData.fullName" placeholder="请输入姓名" />
            </a-form-item>

            <a-form-item label="性别" name="gender">
                <a-select v-model:value="formData.gender" placeholder="请选择性别">
                    <a-select-option value="1">男</a-select-option>
                    <a-select-option value="2">女</a-select-option>
                    <a-select-option value="3">未知</a-select-option>
                </a-select>
            </a-form-item>

            <a-form-item label="状态" name="status">
                <a-select v-model:value="formData.status" placeholder="请选择状态">
                    <a-select-option value="1">启用</a-select-option>
                    <a-select-option value="2">禁用</a-select-option>
                </a-select>
            </a-form-item>

            <a-form-item label="头像" name="avatar">
                <a-input v-model:value="formData.avatar" placeholder="请输入头像URL" />
            </a-form-item>

            <a-form-item label="生日" name="birthday">
                <a-date-picker
                    v-model:value="formData.birthday"
                    placeholder="请选择生日"
                    style="width: 100%"
                />
            </a-form-item>

            <a-form-item label="备注" name="remarks">
                <a-textarea
                    v-model:value="formData.remarks"
                    placeholder="请输入备注"
                    :rows="3"
                />
            </a-form-item>

            <a-form-item label="角色" name="roleIds">
                <a-select
                    v-model:value="formData.roleIds"
                    mode="multiple"
                    placeholder="请选择角色"
                    :loading="rolesLoading"
                    :options="roleOptions"
                />
            </a-form-item>

            <!-- 权限授权 -->
            <a-form-item label="权限授权" v-if="isEdit">
                <a-button type="primary" @click="handleGrantPermission">
                    <template #icon>
                        <SafetyOutlined />
                    </template>
                    管理权限
                </a-button>
                <span style="margin-left: 8px; color: #8c8c8c; font-size: 12px">
                    为用户直接分配权限（独立于角色权限）
                </span>
            </a-form-item>
        </a-form>

        <!-- 用户权限授权弹窗 -->
        <UserPermissionGrant
            v-model:open="grantPermissionOpen"
            :user-id="isEdit && props.user ? props.user.userId : null"
            :current-permissions="currentPermissions"
            @success="handleGrantSuccess"
        />
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, computed, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { SafetyOutlined } from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'
import dayjs, { type Dayjs } from 'dayjs'
import { userApi, type CreateUserRequest, type UpdateUserRequest, type UserResponse, type RoleInfo, type PermissionInfo } from '@/api/user'
import { roleApi, type Role } from '@/api/role'
import UserPermissionGrant from './UserPermissionGrant.vue'

interface Props {
    open: boolean
    user?: UserResponse | null
}

const props = withDefaults(defineProps<Props>(), {
    open: false,
    user: null,
})

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
    (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)
const rolesLoading = ref(false)
const roleOptions = ref<Array<{ label: string; value: string }>>([])
const grantPermissionOpen = ref(false)
const currentPermissions = ref<PermissionInfo[]>([])

const isEdit = computed(() => !!props.user)

const formData = reactive<CreateUserRequest & { birthday?: Dayjs | null }>({
    phone: '',
    email: '',
    password: '',
    nickname: '',
    fullName: '',
    gender: '3',
    status: '1',
    avatar: '',
    birthday: null,
    remarks: '',
    roleIds: [],
})

const rules = computed(() => ({
    phone: [{ required: !isEdit.value, message: '请输入手机号', trigger: 'blur' }],
    password: [{ required: !isEdit.value, message: '请输入密码', trigger: 'blur' }],
    nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
    fullName: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
    email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
}))

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 加载角色列表
const loadRoles = async () => {
    rolesLoading.value = true
    try {
        const response = await roleApi.getAllRoles()
        if (response.code === 200 || response.code === 0) {
            roleOptions.value = (response.data || []).map((role: Role) => ({
                label: role.roleName,
                value: role.roleId,
            }))
        }
    } catch (error: any) {
        console.error('加载角色列表失败:', error)
        // 如果是404，说明接口还未实现，不显示错误提示，只记录日志
        if (error.response?.status === 404) {
            console.warn('角色列表接口未实现，使用空列表')
            roleOptions.value = []
        } else {
            message.warning('加载角色列表失败，角色选择功能可能不可用')
        }
    } finally {
        rolesLoading.value = false
    }
}

// 初始化表单数据
const initFormData = async () => {
    if (props.user) {
        // 编辑模式
        formData.email = props.user.email || ''
        formData.nickname = props.user.nickname || ''
        formData.fullName = props.user.fullName || ''
        formData.gender = props.user.gender || '3'
        formData.status = props.user.status || '1'
        formData.avatar = props.user.avatar || ''
        formData.remarks = props.user.remarks || ''
        if (props.user.birthday) {
            formData.birthday = dayjs(props.user.birthday)
        }
        // 加载用户的角色和权限
        try {
            const detailResponse = await userApi.getUserDetail(props.user.userId)
            if (detailResponse.code === 200 || detailResponse.code === 0) {
                formData.roleIds = (detailResponse.data?.roles || []).map((r: RoleInfo) => r.roleId)
                currentPermissions.value = detailResponse.data?.permissions || []
            }
        } catch (error) {
            console.error('加载用户信息失败:', error)
            formData.roleIds = []
            currentPermissions.value = []
        }
    } else {
        // 添加模式
        formData.phone = ''
        formData.email = ''
        formData.password = ''
        formData.nickname = ''
        formData.fullName = ''
        formData.gender = '3'
        formData.status = '1'
        formData.avatar = ''
        formData.birthday = null
        formData.remarks = ''
        formData.roleIds = []
    }
}

// 提交表单
const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        loading.value = true

        const submitData: any = {
            email: formData.email,
            nickname: formData.nickname,
            fullName: formData.fullName,
            gender: formData.gender,
            status: formData.status,
            avatar: formData.avatar,
            remarks: formData.remarks,
            roleIds: formData.roleIds,
        }

        if (formData.birthday) {
            submitData.birthday = formData.birthday.format('YYYY-MM-DD')
        }

        if (isEdit.value && props.user) {
            // 更新用户
            const response = await userApi.updateUser(props.user.userId, submitData as UpdateUserRequest)
            if (response.code === 200 || response.code === 0) {
                message.success('更新用户成功')
                emit('success')
                handleCancel()
            }
        } else {
            // 创建用户
            submitData.phone = formData.phone
            submitData.password = formData.password
            const response = await userApi.createUser(submitData as CreateUserRequest)
            if (response.code === 200 || response.code === 0) {
                message.success('创建用户成功')
                emit('success')
                handleCancel()
            }
        }
    } catch (error: any) {
        console.error('提交失败:', error)
        if (error.errorFields) {
            return
        }
        message.error(error.message || '操作失败')
    } finally {
        loading.value = false
    }
}

// 处理授权权限
const handleGrantPermission = () => {
    grantPermissionOpen.value = true
}

// 授权成功回调
const handleGrantSuccess = () => {
    grantPermissionOpen.value = false
    // 重新加载用户信息和权限
    if (props.user) {
        userApi.getUserDetail(props.user.userId).then((response) => {
            if (response.code === 200 || response.code === 0) {
                currentPermissions.value = response.data?.permissions || []
            }
        }).catch((error) => {
            console.error('加载用户权限失败:', error)
        })
    }
    message.success('权限授权已更新')
}

// 取消
const handleCancel = () => {
    formRef.value?.resetFields()
    emit('update:open', false)
}

// 监听打开状态和用户数据，加载数据
watch(
    () => [props.open, props.user],
    async ([open, user]) => {
        if (open) {
            loadRoles()
            // 使用 nextTick 确保在 DOM 更新后再初始化表单数据
            await nextTick()
            initFormData()
        } else {
            // 关闭时重置表单
            formRef.value?.resetFields()
            Object.assign(formData, {
                phone: '',
                email: '',
                password: '',
                nickname: '',
                fullName: '',
                gender: '3',
                status: '1',
                avatar: '',
                birthday: null,
                remarks: '',
                roleIds: [],
            })
        }
    },
    { immediate: true }
)
</script>

<style scoped lang="less">
:deep(.ant-form-item-label > label) {
    font-weight: 500;
}
</style>

