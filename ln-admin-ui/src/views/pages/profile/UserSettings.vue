<template>
    <div class="user-settings-container">
        <a-card title="个人设置" :bordered="false">
            <a-tabs>
                <a-tab-pane key="basic" tab="基础设置">
                    <a-form :model="basicForm" layout="vertical" style="max-width: 600px;">
                        <a-form-item label="昵称">
                            <a-input v-model:value="basicForm.nickname" placeholder="请输入昵称" />
                        </a-form-item>
                        <a-form-item label="姓名">
                            <a-input v-model:value="basicForm.fullName" placeholder="请输入姓名" />
                        </a-form-item>
                        <a-form-item label="性别">
                            <a-radio-group v-model:value="basicForm.gender">
                                <a-radio value="M">男</a-radio>
                                <a-radio value="F">女</a-radio>
                            </a-radio-group>
                        </a-form-item>
                        <a-form-item>
                            <a-button type="primary">保存</a-button>
                        </a-form-item>
                    </a-form>
                </a-tab-pane>
                <a-tab-pane key="security" tab="安全设置">
                    <a-form :model="securityForm" layout="vertical" style="max-width: 600px;">
                        <a-form-item label="旧密码">
                            <a-input-password v-model:value="securityForm.oldPassword" placeholder="请输入旧密码" />
                        </a-form-item>
                        <a-form-item label="新密码">
                            <a-input-password v-model:value="securityForm.newPassword" placeholder="请输入新密码" />
                        </a-form-item>
                        <a-form-item label="确认新密码">
                            <a-input-password v-model:value="securityForm.confirmPassword" placeholder="请再次输入新密码" />
                        </a-form-item>
                        <a-form-item>
                            <a-button type="primary">修改密码</a-button>
                        </a-form-item>
                    </a-form>
                </a-tab-pane>
            </a-tabs>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()

// 从store获取用户信息
const userInfo = computed(() => userStore.userInfo)

const basicForm = ref({
    nickname: '',
    fullName: '',
    gender: 'M',
})

const securityForm = ref({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
})

onMounted(async () => {
    // 如果store中没有用户信息，从服务器获取
    if (!userInfo.value) {
        await userStore.fetchUserInfo()
    }
    
    // 填充表单
    if (userInfo.value) {
        basicForm.value.nickname = userInfo.value.nickname || ''
        basicForm.value.fullName = userInfo.value.fullName || ''
        basicForm.value.gender = userInfo.value.gender || 'M'
    }
})
</script>

<style scoped>
.user-settings-container {
    padding: 0;
}
</style>

