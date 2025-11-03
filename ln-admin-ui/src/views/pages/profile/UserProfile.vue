<template>
    <div class="user-profile-container">
        <a-card title="个人资料" :bordered="false">
            <a-descriptions :column="2" bordered>
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
                <a-button type="primary">编辑资料</a-button>
                <a-button style="margin-left: 8px;">修改密码</a-button>
            </div>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()

// 从store获取用户信息
const userInfo = computed(() => userStore.userInfo)

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
</style>

