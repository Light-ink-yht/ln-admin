<template>
    <div class="user-profile-container">
        <a-card title="个人资料" :bordered="false">
            <div class="profile-header">
                <div class="avatar-section">
                    <ImageUploadSingle
                        v-model="avatarUrl"
                        :max-size="5"
                        placeholder="上传头像"
                        preview-alt="用户头像"
                        @upload-success="handleAvatarUploadSuccess"
                    />
                </div>
                <div class="profile-info">
                    <h2 class="profile-name">{{ userInfo?.fullName || userInfo?.nickname || '未设置' }}</h2>
                    <p class="profile-role">{{ userInfo?.email || '-' }}</p>
                </div>
            </div>

            <a-descriptions :column="2" bordered style="margin-top: 24px;">
                <a-descriptions-item label="用户ID">
                    {{ userInfo?.userId || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="昵称">
                    {{ userInfo?.nickname || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="姓名">
                    {{ userInfo?.fullName || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="手机号">
                    {{ userInfo?.phone || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="邮箱">
                    {{ userInfo?.email || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="性别">
                    {{ userInfo?.gender === 'M' || userInfo?.gender === '1' ? '男' : userInfo?.gender === 'F' || userInfo?.gender === '2' ? '女' : '未知' }}
                </a-descriptions-item>
                <a-descriptions-item label="生日">
                    {{ userInfo?.birthday || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="状态">
                    <a-tag v-if="userInfo?.status" :color="userInfo.status === '1' ? 'success' : 'error'">
                        {{ userInfo.status === '1' ? '正常' : '禁用' }}
                    </a-tag>
                    <span v-else>-</span>
                </a-descriptions-item>
                <a-descriptions-item label="登录次数">
                    {{ userInfo?.loginCount || 0 }}
                </a-descriptions-item>
                <a-descriptions-item label="最后登录时间">
                    {{ userInfo?.lastLoginTime || '-' }}
                </a-descriptions-item>
            </a-descriptions>
            <div class="profile-actions" style="margin-top: 24px;">
                <a-button type="primary" @click="handleEditProfile">编辑资料</a-button>
                <a-button style="margin-left: 8px;" @click="handleChangePassword">修改密码</a-button>
            </div>
        </a-card>

        <!-- 编辑资料弹窗 -->
        <a-modal
            v-model:open="editProfileVisible"
            title="编辑资料"
            :confirm-loading="editProfileLoading"
            @ok="handleSaveProfile"
            @cancel="handleCancelEditProfile"
            width="600px"
        >
            <a-form
                ref="editProfileFormRef"
                :model="editProfileForm"
                :rules="editProfileRules"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 18 }"
            >
                <a-form-item label="昵称" name="nickname">
                    <a-input v-model:value="editProfileForm.nickname" placeholder="请输入昵称" />
                </a-form-item>
                <a-form-item label="姓名" name="fullName">
                    <a-input v-model:value="editProfileForm.fullName" placeholder="请输入姓名" />
                </a-form-item>
                <a-form-item label="邮箱" name="email">
                    <a-input v-model:value="editProfileForm.email" placeholder="请输入邮箱" />
                </a-form-item>
                <a-form-item label="性别" name="gender">
                    <a-radio-group v-model:value="editProfileForm.gender">
                        <a-radio value="M">男</a-radio>
                        <a-radio value="F">女</a-radio>
                    </a-radio-group>
                </a-form-item>
                <a-form-item label="生日" name="birthday">
                    <a-date-picker
                        v-model:value="editProfileForm.birthday"
                        placeholder="请选择生日"
                        style="width: 100%"
                        format="YYYY-MM-DD"
                    />
                </a-form-item>
                <a-form-item label="头像" name="avatar">
                    <ImageUploadSingle
                        v-model="editProfileForm.avatar"
                        :max-size="5"
                        placeholder="上传头像"
                        preview-alt="用户头像"
                    />
                </a-form-item>
            </a-form>
        </a-modal>

        <!-- 修改密码弹窗 -->
        <a-modal
            v-model:open="changePasswordVisible"
            title="修改密码"
            :confirm-loading="changePasswordLoading"
            @ok="handleSavePassword"
            @cancel="handleCancelChangePassword"
            width="500px"
        >
            <a-form
                ref="changePasswordFormRef"
                :model="changePasswordForm"
                :rules="changePasswordRules"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 18 }"
            >
                <a-form-item label="旧密码" name="oldPassword">
                    <a-input-password v-model:value="changePasswordForm.oldPassword" placeholder="请输入旧密码" />
                </a-form-item>
                <a-form-item label="新密码" name="newPassword">
                    <a-input-password v-model:value="changePasswordForm.newPassword" placeholder="请输入新密码" />
                </a-form-item>
                <a-form-item label="确认新密码" name="confirmPassword">
                    <a-input-password v-model:value="changePasswordForm.confirmPassword" placeholder="请再次输入新密码" />
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { UserOutlined } from '@ant-design/icons-vue'
import dayjs, { type Dayjs } from 'dayjs'
import { useUserStore } from '@/stores/modules/user'
import { userApi, type UpdateUserRequest } from '@/api/user'
import { ImageUploadSingle } from '@/components/upload'
import type { FileResponse } from '@/api/file'

const userStore = useUserStore()

// 从store获取用户信息
const userInfo = computed(() => userStore.userInfo)

// 头像URL（用于顶部显示）
const avatarUrl = ref(userInfo.value?.avatar || '')

// 监听用户信息变化，更新头像
watch(() => userInfo.value?.avatar, (newAvatar) => {
    avatarUrl.value = newAvatar || ''
})

// 编辑资料相关
const editProfileVisible = ref(false)
const editProfileLoading = ref(false)
const editProfileFormRef = ref<FormInstance>()
const editProfileForm = reactive({
    nickname: '',
    fullName: '',
    email: '',
    gender: 'M',
    birthday: null as Dayjs | null,
    avatar: '',
})

const editProfileRules = {
    nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
    fullName: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
    email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
    ],
}

// 修改密码相关
const changePasswordVisible = ref(false)
const changePasswordLoading = ref(false)
const changePasswordFormRef = ref<FormInstance>()
const changePasswordForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
})

const validateConfirmPassword = (_rule: any, value: string) => {
    if (!value) {
        return Promise.reject('请再次输入新密码')
    }
    if (value !== changePasswordForm.newPassword) {
        return Promise.reject('两次输入的密码不一致')
    }
    return Promise.resolve()
}

const changePasswordRules = {
    oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
    newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    ],
    confirmPassword: [
        { required: true, validator: validateConfirmPassword, trigger: 'blur' },
    ],
}

// 打开编辑资料弹窗
const handleEditProfile = () => {
    if (!userInfo.value) {
        message.warning('用户信息不存在')
        return
    }
    editProfileForm.nickname = userInfo.value.nickname || ''
    editProfileForm.fullName = userInfo.value.fullName || ''
    editProfileForm.email = userInfo.value.email || ''
    editProfileForm.gender = userInfo.value.gender || 'M'
    editProfileForm.birthday = userInfo.value.birthday ? dayjs(userInfo.value.birthday) : null
    editProfileForm.avatar = userInfo.value.avatar || ''
    editProfileVisible.value = true
}

// 保存编辑的资料
const handleSaveProfile = async () => {
    try {
        await editProfileFormRef.value?.validate()
        if (!userInfo.value) {
            message.error('用户信息不存在')
            return
        }

        editProfileLoading.value = true

        const submitData: UpdateUserRequest = {
            nickname: editProfileForm.nickname,
            fullName: editProfileForm.fullName,
            email: editProfileForm.email,
            gender: editProfileForm.gender,
            avatar: editProfileForm.avatar,
        }

        if (editProfileForm.birthday) {
            submitData.birthday = editProfileForm.birthday.format('YYYY-MM-DD')
        }

        const response = await userApi.updateUser(userInfo.value.userId, submitData)
        if (response.code === 200 || response.code === 0) {
            message.success('更新资料成功')
            // 刷新用户信息
            await userStore.fetchUserInfo()
            editProfileVisible.value = false
        }
    } catch (error: any) {
        console.error('保存资料失败:', error)
        if (error.errorFields) {
            return
        }
        message.error(error.message || '保存失败')
    } finally {
        editProfileLoading.value = false
    }
}

// 取消编辑资料
const handleCancelEditProfile = () => {
    editProfileVisible.value = false
    editProfileFormRef.value?.resetFields()
}

// 打开修改密码弹窗
const handleChangePassword = () => {
    changePasswordForm.oldPassword = ''
    changePasswordForm.newPassword = ''
    changePasswordForm.confirmPassword = ''
    changePasswordVisible.value = true
}

// 保存新密码
const handleSavePassword = async () => {
    try {
        await changePasswordFormRef.value?.validate()
        if (!userInfo.value) {
            message.error('用户信息不存在')
            return
        }

        changePasswordLoading.value = true

        // TODO: 调用后端修改密码接口
        // 目前后端可能没有专门的修改密码接口，需要先验证旧密码
        // 这里先提示用户功能待实现
        message.warning('修改密码功能需要后端接口支持，请联系管理员')
        
        // 如果后端有接口，可以这样调用：
        // const response = await userApi.changePassword({
        //     oldPassword: changePasswordForm.oldPassword,
        //     newPassword: changePasswordForm.newPassword,
        // })
        
        changePasswordVisible.value = false
    } catch (error: any) {
        console.error('修改密码失败:', error)
        if (error.errorFields) {
            return
        }
        message.error(error.message || '修改密码失败')
    } finally {
        changePasswordLoading.value = false
    }
}

// 取消修改密码
const handleCancelChangePassword = () => {
    changePasswordVisible.value = false
    changePasswordFormRef.value?.resetFields()
}

// 头像上传成功回调（顶部头像上传）
const handleAvatarUploadSuccess = async (file: FileResponse) => {
    try {
        const avatarUrl = file.file_url || file.file_path || ''
        if (avatarUrl && userInfo.value) {
            // 更新头像
            const updateData: UpdateUserRequest = {
                avatar: avatarUrl,
            }
            const updateResponse = await userApi.updateUser(userInfo.value.userId, updateData)
            if (updateResponse.code === 200 || updateResponse.code === 0) {
                message.success('头像上传成功')
                await userStore.fetchUserInfo()
            } else {
                message.error(updateResponse.msg || '更新头像失败')
            }
        }
    } catch (error: any) {
        console.error('更新头像失败:', error)
        message.error(error.message || '更新头像失败')
    }
}

onMounted(async () => {
    // 如果store中没有用户信息，从服务器获取
    if (!userInfo.value) {
        await userStore.fetchUserInfo()
    }
})
</script>

<style scoped>
.user-profile-container {
    padding: 0;
}

.profile-header {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    margin-bottom: 24px;
}

.avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
}

.avatar-section :deep(.image-upload-single) {
    .preview-wrapper {
        width: 120px;
        padding-top: 120px;
        border: 4px solid rgba(255, 255, 255, 0.3);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        background: rgba(255, 255, 255, 0.1);
    }

    .upload-button {
        width: 120px;
        min-height: 120px;
        border: 2px dashed rgba(255, 255, 255, 0.5);
        background: rgba(255, 255, 255, 0.1);
    }

    .upload-button:hover {
        border-color: rgba(255, 255, 255, 0.8);
        background: rgba(255, 255, 255, 0.2);
    }
}

.profile-info {
    flex: 1;
    color: #fff;
}

.profile-name {
    margin: 0 0 8px 0;
    font-size: 24px;
    font-weight: 500;
    color: #fff;
}

.profile-role {
    margin: 0;
    font-size: 14px;
    color: rgba(255, 255, 255, 0.85);
}

.avatar-edit-wrapper {
    width: 100%;
}
</style>

