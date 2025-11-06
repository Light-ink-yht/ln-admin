<template>
    <a-modal
        v-model:open="visible"
        :title="isEdit ? '编辑用户' : '添加用户'"
        :width="900"
        :confirm-loading="loading"
        @ok="handleSubmit"
        @cancel="handleCancel"
        :styles="{ body: { padding: '24px' } }"
    >
        <a-form
            ref="formRef"
            :model="formData"
            :rules="rules"
            layout="vertical"
        >
            <!-- 头像上传区域 -->
            <div class="form-header-section">
                <div class="avatar-upload-card">
                    <ImageUploadSingle
                        v-model="formData.avatar"
                        :max-size="5"
                        placeholder="上传头像"
                        preview-alt="用户头像"
                    />
                    <div class="avatar-info">
                        <p class="avatar-tip">支持 JPG、PNG、GIF 格式，大小不超过 5MB</p>
                    </div>
                </div>
            </div>

            <a-divider style="margin: 24px 0;" />

            <!-- 基本信息 -->
            <div class="form-section">
                <div class="section-title">
                    <UserOutlined />
                    <span>基本信息</span>
                </div>
                <a-row :gutter="16">
                    <a-col :span="12">
                        <a-form-item label="手机号" name="phone" v-if="!isEdit">
                            <a-input v-model:value="formData.phone" placeholder="请输入手机号" size="large" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="密码" name="password" v-if="!isEdit">
                            <a-input-password v-model:value="formData.password" placeholder="请输入密码" size="large" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="邮箱" name="email">
                            <a-input v-model:value="formData.email" placeholder="请输入邮箱" size="large" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="昵称" name="nickname">
                            <a-input v-model:value="formData.nickname" placeholder="请输入昵称" size="large" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="姓名" name="fullName">
                            <a-input v-model:value="formData.fullName" placeholder="请输入姓名" size="large" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="性别" name="gender">
                            <a-select v-model:value="formData.gender" placeholder="请选择性别" size="large">
                                <a-select-option value="1">男</a-select-option>
                                <a-select-option value="2">女</a-select-option>
                                <a-select-option value="3">未知</a-select-option>
                            </a-select>
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="状态" name="status">
                            <a-select v-model:value="formData.status" placeholder="请选择状态" size="large">
                                <a-select-option value="1">启用</a-select-option>
                                <a-select-option value="2">禁用</a-select-option>
                            </a-select>
                        </a-form-item>
                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="生日" name="birthday">
                            <a-date-picker
                                v-model:value="formData.birthday"
                                placeholder="请选择生日"
                                style="width: 100%"
                                size="large"
                            />
                        </a-form-item>
                    </a-col>
                    <a-col :span="24">
                        <a-form-item label="备注" name="remarks">
                            <a-textarea
                                v-model:value="formData.remarks"
                                placeholder="请输入备注"
                                :rows="3"
                            />
                        </a-form-item>
                    </a-col>
                </a-row>
            </div>

            <a-divider v-if="!isSuperAdmin" style="margin: 24px 0;" />

            <!-- 角色和权限（超级管理员不显示） -->
            <div v-if="!isSuperAdmin" class="form-section">
                <div class="section-title">
                    <SafetyOutlined />
                    <span>角色权限</span>
                </div>
                <a-row :gutter="16">
                    <a-col :span="24">
                        <a-form-item label="角色" name="roleIds">
                            <a-select
                                v-model:value="formData.roleIds"
                                mode="multiple"
                                placeholder="请选择角色"
                                :loading="rolesLoading"
                                :options="roleOptions"
                                size="large"
                            />
                        </a-form-item>
                    </a-col>
                    <a-col :span="24">
                        <a-form-item label="权限授权">
                            <a-button 
                                type="primary" 
                                @click="handleGrantPermission"
                                :disabled="!isEdit && !createdUserId"
                            >
                                <template #icon>
                                    <SafetyOutlined />
                                </template>
                                {{ isEdit ? '管理权限' : '分配权限' }}
                            </a-button>
                            <span class="form-item-tip">
                                {{ isEdit ? '为用户直接分配权限（独立于角色权限）' : '创建用户成功后，可以为用户分配权限' }}
                            </span>
                        </a-form-item>
                    </a-col>
                </a-row>
            </div>
        </a-form>

        <!-- 用户权限授权弹窗 -->
        <UserPermissionGrant
            v-model:open="grantPermissionOpen"
            :user-id="(isEdit && props.user ? props.user.userId : createdUserId) || null"
            :current-permissions="currentPermissions"
            @success="handleGrantSuccess"
        />
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, computed, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { SafetyOutlined, UserOutlined } from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'
import dayjs, { type Dayjs } from 'dayjs'
import { userApi, type CreateUserRequest, type UpdateUserRequest, type UserResponse, type RoleInfo, type PermissionInfo } from '@/api/user'
import { roleApi, type Role } from '@/api/role'
import { ImageUploadSingle } from '@/components/upload'
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
const createdUserId = ref<string | null>(null) // 用于存储新创建的用户ID

const isEdit = computed(() => !!props.user)

// 用户详情（用于判断角色）
const userDetail = ref<any>(null)

// 判断是否是超级管理员
const isSuperAdmin = computed(() => {
    if (!props.user) {
        return false
    }
    // 通过手机号判断
    if (props.user.phone === '18797131041') {
        return true
    }
    // 通过角色判断（从用户详情中获取角色信息）
    if (userDetail.value && userDetail.value.roles && Array.isArray(userDetail.value.roles)) {
        return userDetail.value.roles.some((role: any) => role.roleKey === 'super_admin')
    }
    return false
})

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

// 加载角色列表（过滤掉超级管理员角色）
const loadRoles = async () => {
    rolesLoading.value = true
    try {
        const response = await roleApi.getAllRoles()
        if (response.code === 200 || response.code === 0) {
            // 过滤掉超级管理员角色
            roleOptions.value = (response.data || [])
                .filter((role: Role) => role.roleKey !== 'super_admin')
                .map((role: Role) => ({
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
        // 重置创建用户ID，确保不会意外打开权限授权弹窗
        createdUserId.value = null
        grantPermissionOpen.value = false
        
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
                userDetail.value = detailResponse.data // 保存用户详情用于判断超级管理员
                formData.roleIds = (detailResponse.data?.roles || []).map((r: RoleInfo) => r.roleId)
                currentPermissions.value = detailResponse.data?.permissions || []
            }
        } catch (error) {
            console.error('加载用户信息失败:', error)
            userDetail.value = null
            formData.roleIds = []
            currentPermissions.value = []
        }
    } else {
        // 添加模式
        // 重置创建用户ID和权限授权弹窗状态
        createdUserId.value = null
        grantPermissionOpen.value = false
        
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
        }
        
        // 超级管理员不提交角色信息
        if (!isSuperAdmin.value) {
            submitData.roleIds = formData.roleIds
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
                // 保存新创建的用户ID（注意后端返回的字段名是 userId）
                const newUserId = response.data?.userId || null
                if (newUserId) {
                    createdUserId.value = newUserId
                    // 创建成功后，自动打开权限授权弹窗
                    await nextTick()
                    grantPermissionOpen.value = true
                    // 加载用户的权限信息（新用户应该没有权限）
                    currentPermissions.value = []
                } else {
                    // 如果无法获取用户ID，直接关闭并刷新列表
                    emit('success')
                    handleCancel()
                }
            }
        }
    } catch (error: any) {
        console.error('提交失败:', error)
        if (error.errorFields) {
            return
        }
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        loading.value = false
    }
}

// 处理授权权限
const handleGrantPermission = async () => {
    // 编辑模式下，需要先加载用户权限，然后打开授权弹窗
    if (isEdit.value && props.user) {
        // 确保权限已加载
        if (currentPermissions.value.length === 0) {
            try {
                const detailResponse = await userApi.getUserDetail(props.user.userId)
                if (detailResponse.code === 200 || detailResponse.code === 0) {
                    currentPermissions.value = detailResponse.data?.permissions || []
                }
            } catch (error) {
                console.error('加载用户权限失败:', error)
            }
        }
        grantPermissionOpen.value = true
        return
    }
    // 添加模式下，如果用户已创建，打开授权弹窗
    if (createdUserId.value) {
        grantPermissionOpen.value = true
        return
    }
    // 否则提示用户先创建用户
    message.warning('请先创建用户，然后再分配权限')
}

// 授权成功回调
const handleGrantSuccess = () => {
    grantPermissionOpen.value = false
    const userId = (isEdit.value && props.user) ? props.user.userId : createdUserId.value
    // 重新加载用户信息和权限
    if (userId) {
        userApi.getUserDetail(userId).then((response) => {
            if (response.code === 200 || response.code === 0) {
                currentPermissions.value = response.data?.permissions || []
            }
        }).catch((error) => {
            console.error('加载用户权限失败:', error)
        })
    }
    message.success('权限授权已更新')
    // 如果是新创建的用户，授权成功后关闭表单并刷新列表
    if (!isEdit.value && createdUserId.value) {
        emit('success')
        handleCancel()
    }
}


// 取消
const handleCancel = () => {
    formRef.value?.resetFields()
    createdUserId.value = null // 重置创建的用户ID
    currentPermissions.value = [] // 重置权限列表
    userDetail.value = null // 重置用户详情
    emit('update:open', false)
}

// 监听打开状态和用户数据，加载数据
watch(
    () => [props.open, props.user],
    async ([open, user]) => {
        if (open) {
            // 确保权限授权弹窗关闭
            grantPermissionOpen.value = false
            // 重置创建用户ID
            createdUserId.value = null
            
            loadRoles()
            // 使用 nextTick 确保在 DOM 更新后再初始化表单数据
            await nextTick()
            initFormData()
        } else {
            // 关闭时重置表单和状态
            grantPermissionOpen.value = false
            createdUserId.value = null
            currentPermissions.value = []
            userDetail.value = null // 重置用户详情
            
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

<style scoped>
.form-header-section {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;
}

.avatar-upload-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
}


.avatar-info {
    text-align: center;
    margin-top: 12px;
}

.avatar-tip {
    margin: 0;
    color: rgba(0, 0, 0, 0.45);
    font-size: 12px;
}

.form-section {
    margin-bottom: 24px;
}

.section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    color: rgba(0, 0, 0, 0.85);
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 2px solid #f0f0f0;
}

.section-title :deep(.anticon) {
    font-size: 18px;
    color: #1890ff;
}

.form-item-tip {
    margin-left: 12px;
    color: rgba(0, 0, 0, 0.45);
    font-size: 12px;
}

:deep(.ant-form-item-label > label) {
    font-weight: 500;
    color: rgba(0, 0, 0, 0.85);
}

:deep(.ant-input),
:deep(.ant-select),
:deep(.ant-picker) {
    border-radius: 6px;
}

:deep(.ant-input:focus),
:deep(.ant-select-focused .ant-select-selector),
:deep(.ant-picker-focused) {
    border-color: #1890ff;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}
</style>

