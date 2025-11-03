<template>
    <div class="welcome-banner">
        <!-- 背景装饰元素 -->
        <div class="banner-decoration">
            <div class="decoration-circle circle-1"></div>
            <div class="decoration-circle circle-2"></div>
            <div class="decoration-circle circle-3"></div>
            <div class="decoration-wave wave-1"></div>
            <div class="decoration-wave wave-2"></div>
            <div class="decoration-grid"></div>
        </div>
        
        <!-- 内容区域 -->
        <div class="welcome-content">
            <div class="welcome-text">
                <h1 class="welcome-title">{{ welcomeText }}</h1>
                <p class="welcome-subtitle">{{ currentTime }}</p>
            </div>
            <div class="welcome-avatar">
                <a-avatar :size="80" :src="userAvatar">
                    <template v-if="!userAvatar">
                        <UserOutlined />
                    </template>
                </a-avatar>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { UserOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()

const userAvatar = computed(() => userStore.userAvatar)

// 欢迎文本
const welcomeText = computed(() => {
    const hour = new Date().getHours()
    const userName = userStore.userName
    if (hour < 6) return `深夜好，${userName}`
    if (hour < 9) return `早上好，${userName}`
    if (hour < 12) return `上午好，${userName}`
    if (hour < 14) return `中午好，${userName}`
    if (hour < 18) return `下午好，${userName}`
    if (hour < 22) return `晚上好，${userName}`
    return `夜深了，${userName}`
})

const currentTime = ref('')
let timeInterval: number | null = null

// 更新时间
const updateTime = () => {
    const now = new Date()
    currentTime.value = now.toLocaleString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'long',
        hour: '2-digit',
        minute: '2-digit',
    })
}

onMounted(() => {
    updateTime()
    timeInterval = window.setInterval(updateTime, 1000)
})

onUnmounted(() => {
    if (timeInterval) {
        clearInterval(timeInterval)
    }
})
</script>

<style scoped lang="less">
.welcome-banner {
    position: relative;
    background: linear-gradient(135deg, 
        var(--color-primary) 0%, 
        var(--color-primary-hover) 50%,
        var(--color-primary) 100%);
    border-radius: 16px;
    padding: 32px;
    margin-bottom: 24px;
    color: #fff;
    box-shadow: 
        0 8px 24px rgba(var(--color-primary-rgb, 24, 144, 255), 0.25),
        0 2px 8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    isolation: isolate;

    // 渐变光晕效果
    &::before {
        content: '';
        position: absolute;
        top: -50%;
        right: -20%;
        width: 200%;
        height: 200%;
        background: radial-gradient(
            circle,
            rgba(255, 255, 255, 0.15) 0%,
            rgba(255, 255, 255, 0.05) 30%,
            transparent 70%
        );
        animation: glow-rotate 20s linear infinite;
        pointer-events: none;
    }

    &::after {
        content: '';
        position: absolute;
        bottom: -30%;
        left: -10%;
        width: 150%;
        height: 150%;
        background: radial-gradient(
            circle,
            rgba(255, 255, 255, 0.1) 0%,
            transparent 50%
        );
        animation: glow-pulse 8s ease-in-out infinite;
        pointer-events: none;
    }

    .welcome-content {
        position: relative;
        z-index: 2;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .welcome-text {
        flex: 1;
        position: relative;
        z-index: 3;
    }

    .welcome-title {
        font-size: 32px;
        font-weight: 600;
        margin-bottom: 8px;
        color: #fff;
        text-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    .welcome-subtitle {
        font-size: 16px;
        opacity: 0.95;
        margin: 0;
        text-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
    }

    .welcome-avatar {
        position: relative;
        z-index: 3;
        :deep(.ant-avatar) {
            background: rgba(255, 255, 255, 0.25);
            border: 3px solid rgba(255, 255, 255, 0.4);
            box-shadow: 
                0 4px 12px rgba(0, 0, 0, 0.15),
                0 0 0 4px rgba(255, 255, 255, 0.1);
            transition: all 0.3s ease;
            
            &:hover {
                transform: scale(1.05);
                box-shadow: 
                    0 6px 16px rgba(0, 0, 0, 0.2),
                    0 0 0 6px rgba(255, 255, 255, 0.15);
            }
        }
    }
}

// 背景装饰
.banner-decoration {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1;
    pointer-events: none;
    overflow: hidden;
}

// 装饰圆圈
.decoration-circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    animation: float 15s ease-in-out infinite;
}

.circle-1 {
    width: 200px;
    height: 200px;
    top: -100px;
    right: 10%;
    animation-delay: 0s;
}

.circle-2 {
    width: 150px;
    height: 150px;
    bottom: -75px;
    left: 15%;
    background: rgba(255, 255, 255, 0.08);
    animation-delay: 2s;
    animation-duration: 18s;
}

.circle-3 {
    width: 120px;
    height: 120px;
    top: 50%;
    right: 5%;
    background: rgba(255, 255, 255, 0.06);
    animation-delay: 4s;
    animation-duration: 20s;
}

// 装饰波浪
.decoration-wave {
    position: absolute;
    width: 200%;
    height: 100%;
    opacity: 0.1;
    background: repeating-linear-gradient(
        45deg,
        transparent,
        transparent 20px,
        rgba(255, 255, 255, 0.1) 20px,
        rgba(255, 255, 255, 0.1) 22px
    );
    animation: wave-move 25s linear infinite;
}

.wave-1 {
    top: 0;
    left: -100%;
    animation-duration: 25s;
}

.wave-2 {
    bottom: 0;
    right: -100%;
    animation-duration: 30s;
    animation-direction: reverse;
    opacity: 0.08;
}

// 网格装饰
.decoration-grid {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: 
        linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 40px 40px;
    opacity: 0.3;
    animation: grid-shift 20s linear infinite;
}

// 动画定义
@keyframes glow-rotate {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}

@keyframes glow-pulse {
    0%, 100% {
        opacity: 0.3;
        transform: scale(1);
    }
    50% {
        opacity: 0.5;
        transform: scale(1.1);
    }
}

@keyframes float {
    0%, 100% {
        transform: translate(0, 0) scale(1);
        opacity: 0.6;
    }
    33% {
        transform: translate(20px, -20px) scale(1.1);
        opacity: 0.8;
    }
    66% {
        transform: translate(-15px, 15px) scale(0.9);
        opacity: 0.5;
    }
}

@keyframes wave-move {
    0% {
        transform: translateX(0);
    }
    100% {
        transform: translateX(50%);
    }
}

@keyframes grid-shift {
    0% {
        background-position: 0 0;
    }
    100% {
        background-position: 40px 40px;
    }
}

// 响应式设计
@media (max-width: 768px) {
    .welcome-banner {
        padding: 24px;

        .welcome-title {
            font-size: 24px;
        }

        .welcome-subtitle {
            font-size: 14px;
        }

        .circle-1,
        .circle-2,
        .circle-3 {
            display: none;
        }
    }
}

// 暗色模式适配
body.dark-mode {
    .welcome-banner {
        &::before,
        &::after {
            opacity: 0.5;
        }

        .decoration-circle {
            background: rgba(255, 255, 255, 0.05);
        }

        .decoration-grid {
            opacity: 0.2;
        }
    }
}
</style>

