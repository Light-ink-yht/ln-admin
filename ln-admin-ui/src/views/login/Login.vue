<template>
    <div class="login-container">
        <div class="login-wrapper">
            <div class="login-left">
                <div class="brand">
                    <h1 class="brand-title">LN Admin</h1>
                    <p class="brand-subtitle">企业级管理系统</p>
                </div>
                <div class="features">
                    <div class="feature-item">
                        <div class="feature-icon">🎨</div>
                        <div class="feature-content">
                            <h3>动态主题</h3>
                            <p>支持多种主题切换</p>
                        </div>
                    </div>
                    <div class="feature-item">
                        <div class="feature-icon">⚡</div>
                        <div class="feature-content">
                            <h3>极速体验</h3>
                            <p>基于 Vue 3 + Vite</p>
                        </div>
                    </div>
                    <div class="feature-item">
                        <div class="feature-icon">🔒</div>
                        <div class="feature-content">
                            <h3>安全可靠</h3>
                            <p>完善的权限管理</p>
                        </div>
                    </div>
                </div>
            </div>

            <div class="login-right">
                <div class="form-wrapper">
                    <!-- 标签页 -->
                    <a-tabs v-model:activeKey="activeTab" :centered="true" class="login-tabs">
                        <a-tab-pane key="login" tab="登录">
                            <a-form ref="loginFormRef" :model="loginForm" :rules="loginRules" @finish="handleLogin"
                                layout="vertical" class="login-form">
                                <a-form-item label="手机号" name="phone">
                                    <a-input v-model:value="loginForm.phone" placeholder="请输入手机号" size="large"
                                        prefix="📱" />
                                </a-form-item>

                                <a-form-item label="密码" name="password">
                                    <a-input-password v-model:value="loginForm.password" placeholder="请输入密码"
                                        size="large" prefix="🔒" />
                                </a-form-item>

                                <a-form-item label="图片验证码" name="code">
                                    <a-row :gutter="8">
                                        <a-col :span="14">
                                            <a-input v-model:value="loginForm.code" placeholder="请输入验证码" size="large"
                                                @pressEnter="handleLogin" />
                                        </a-col>
                                        <a-col :span="10">
                                            <div class="captcha-wrapper" @click="refreshCaptcha">
                                                <img v-if="captchaImage" :src="captchaImage" alt="验证码"
                                                    class="captcha-image" />
                                                <div class="captcha-refresh">
                                                    <ReloadOutlined :class="{ 'icon-spin': loadingCaptcha }" />
                                                </div>
                                            </div>
                                        </a-col>
                                    </a-row>
                                </a-form-item>

                                <a-form-item>
                                    <a-button type="primary" html-type="submit" size="large" :loading="loginLoading"
                                        block :style="{ background: primaryColor }">
                                        {{ loginLoading ? '登录中...' : '立即登录' }}
                                    </a-button>
                                </a-form-item>

                                <!--                                <div class="form-footer">-->
                                <!--                                    <a @click="activeTab = 'forgot'">忘记密码？</a>-->
                                <!--                                    <a @click="activeTab = 'register'">还没有账号？立即注册</a>-->
                                <!--                                </div>-->
                            </a-form>
                        </a-tab-pane>

                        <a-tab-pane key="register" tab="注册">
                            <a-form ref="registerFormRef" :model="registerForm" :rules="registerRules"
                                @finish="handleRegister" layout="vertical" class="login-form">
                                <a-form-item label="手机号" name="phone">
                                    <a-input v-model:value="registerForm.phone" placeholder="请输入手机号" size="large"
                                        prefix="📱" />
                                </a-form-item>

                                <a-form-item label="短信验证码" name="code">
                                    <a-row :gutter="8">
                                        <a-col :span="14">
                                            <a-input v-model:value="registerForm.code" placeholder="请输入验证码"
                                                size="large" />
                                        </a-col>
                                        <a-col :span="10">
                                            <a-button size="large" block :disabled="smsLoading"
                                                @click="sendSms('register')">
                                                {{ smsButtonText }}
                                            </a-button>
                                        </a-col>
                                    </a-row>
                                </a-form-item>

                                <a-form-item label="密码" name="password">
                                    <a-input-password v-model:value="registerForm.password" placeholder="请设置密码"
                                        size="large" prefix="🔒" />
                                </a-form-item>

                                <a-form-item label="确认密码" name="confirmPassword">
                                    <a-input-password v-model:value="registerForm.confirmPassword" placeholder="请再次输入密码"
                                        size="large" prefix="🔒" />
                                </a-form-item>

                                <a-form-item>
                                    <a-button type="primary" html-type="submit" size="large" :loading="registerLoading"
                                        block :style="{ background: primaryColor }">
                                        {{ registerLoading ? '注册中...' : '立即注册' }}
                                    </a-button>
                                </a-form-item>

                                <!--                                <div class="form-footer">-->
                                <!--                                    <a @click="activeTab = 'login'">已有账号？立即登录</a>-->
                                <!--                                </div>-->
                            </a-form>
                        </a-tab-pane>

                        <a-tab-pane key="forgot" tab="忘记密码">
                            <a-form ref="forgotFormRef" :model="forgotForm" :rules="forgotRules"
                                @finish="handleForgotPassword" layout="vertical" class="login-form">
                                <a-form-item label="手机号" name="phone">
                                    <a-input v-model:value="forgotForm.phone" placeholder="请输入手机号" size="large"
                                        prefix="📱" />
                                </a-form-item>

                                <a-form-item label="短信验证码" name="code">
                                    <a-row :gutter="8">
                                        <a-col :span="14">
                                            <a-input v-model:value="forgotForm.code" placeholder="请输入验证码"
                                                size="large" />
                                        </a-col>
                                        <a-col :span="10">
                                            <a-button size="large" block :disabled="smsLoading"
                                                @click="sendSms('forgot')">
                                                {{ smsButtonText }}
                                            </a-button>
                                        </a-col>
                                    </a-row>
                                </a-form-item>

                                <a-form-item label="新密码" name="password">
                                    <a-input-password v-model:value="forgotForm.password" placeholder="请设置新密码"
                                        size="large" prefix="🔒" />
                                </a-form-item>

                                <a-form-item label="确认密码" name="confirmPassword">
                                    <a-input-password v-model:value="forgotForm.confirmPassword" placeholder="请再次输入新密码"
                                        size="large" prefix="🔒" />
                                </a-form-item>

                                <a-form-item>
                                    <a-button type="primary" html-type="submit" size="large" :loading="forgotLoading"
                                        block :style="{ background: primaryColor }">
                                        {{ forgotLoading ? '重置中...' : '重置密码' }}
                                    </a-button>
                                </a-form-item>

                                <!--                                <div class="form-footer">-->
                                <!--                                    <a @click="activeTab = 'login'">返回登录</a>-->
                                <!--                                </div>-->
                            </a-form>
                        </a-tab-pane>
                    </a-tabs>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { authApi } from '@/api/auth'
import type { FormInstance } from 'ant-design-vue'

const router = useRouter()
const themeStore = useThemeStore()

// 主题色
const primaryColor = computed(() => themeStore.colorPrimary)

// 当前标签页
const activeTab = ref('login')

// 表单引用
const loginFormRef = ref<FormInstance>()
const registerFormRef = ref<FormInstance>()
const forgotFormRef = ref<FormInstance>()

// 验证码
const captchaImage = ref('')
const captchaId = ref('')
const loadingCaptcha = ref(false)

// 短信验证码
const smsLoading = ref(false)
const smsCountdown = ref(0)
const smsButtonText = computed(() => {
    if (smsCountdown.value > 0) {
        return `${smsCountdown.value}秒后重试`
    }
    return '发送验证码'
})

// 加载状态
const loginLoading = ref(false)
const registerLoading = ref(false)
const forgotLoading = ref(false)

// 登录表单
const loginForm = ref({
    phone: '',
    password: '',
    code: '',
})

// 注册表单
const registerForm = ref({
    phone: '',
    code: '',
    password: '',
    confirmPassword: '',
})

// 忘记密码表单
const forgotForm = ref({
    phone: '',
    code: '',
    password: '',
    confirmPassword: '',
})

// 登录表单验证规则
const loginRules = {
    phone: [
        { required: true, message: '请输入手机号', trigger: 'blur' },
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' },
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度为8位以上，区分大小写。', trigger: 'blur' },
    ],
    code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

// 注册表单验证规则
const registerRules = {
    phone: [
        { required: true, message: '请输入手机号', trigger: 'blur' },
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' },
    ],
    code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度为8位以上，区分大小写。', trigger: 'blur' },
    ],
    confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
            validator: (_rule: any, value: string) => {
                if (value !== registerForm.value.password) {
                    return Promise.reject('两次输入的密码不一致')
                }
                return Promise.resolve()
            },
            trigger: 'blur',
        },
    ],
}

// 忘记密码表单验证规则
const forgotRules = {
    phone: [
        { required: true, message: '请输入手机号', trigger: 'blur' },
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' },
    ],
    code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
    password: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度为6-20位', trigger: 'blur' },
    ],
    confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
            validator: (_rule: any, value: string) => {
                if (value !== forgotForm.value.password) {
                    return Promise.reject('两次输入的密码不一致')
                }
                return Promise.resolve()
            },
            trigger: 'blur',
        },
    ],
}

// 获取图片验证码
const getCaptcha = async () => {
    try {
        loadingCaptcha.value = true
        const res = await authApi.getCaptcha()
        captchaImage.value = res.data.captcha_code
        captchaId.value = res.data.captchaId
        loginForm.value.code = ''
    } catch (error) {
        console.error('获取验证码失败:', error)
    } finally {
        loadingCaptcha.value = false
    }
}

// 刷新验证码
const refreshCaptcha = () => {
    getCaptcha()
}

// 发送短信验证码
const sendSms = async (type: 'register' | 'forgot') => {
    const phone = type === 'register' ? registerForm.value.phone : forgotForm.value.phone

    if (!/^1[3-9]\d{9}$/.test(phone)) {
        message.error('请输入正确的手机号')
        return
    }

    if (smsCountdown.value > 0) {
        return
    }

    try {
        smsLoading.value = true
        const res = await authApi.sendSms({ phone, type })

        // 无论成功还是警告，都开始倒计时
        smsCountdown.value = 60
        const timer = setInterval(() => {
            smsCountdown.value--
            if (smsCountdown.value <= 0) {
                clearInterval(timer)
            }
        }, 1000)
    } catch (error: any) {
        // code=1 或 code=2 都会进入这里
        // 但 code=1 是警告，验证码可能已经发送，所以也启动倒计时
        if (error.message && error.message.includes('验证码')) {
            // 如果是验证码相关的警告，也启动倒计时
            smsCountdown.value = 60
            const timer = setInterval(() => {
                smsCountdown.value--
                if (smsCountdown.value <= 0) {
                    clearInterval(timer)
                }
            }, 1000)
        }
        console.error('发送验证码失败:', error)
    } finally {
        smsLoading.value = false
    }
}

// 登录
const handleLogin = async () => {
    try {
        loginLoading.value = true
        const res = await authApi.login({
            phone: loginForm.value.phone,
            password: loginForm.value.password,
            code: loginForm.value.code,
            captchaId: captchaId.value,
        })
        setTimeout(() => {
            router.push('/')
        }, 500)
    } catch (error) {
        console.error('登录失败:', error)
        refreshCaptcha()
    } finally {
        loginLoading.value = false
    }
}

// 注册
const handleRegister = async () => {
    try {
        registerLoading.value = true
        const res = await authApi.register(registerForm.value)
        activeTab.value = 'login'
        // 清空表单
        registerFormRef.value?.resetFields()
    } catch (error) {
        console.error('注册失败:', error)
    } finally {
        registerLoading.value = false
    }
}

// 忘记密码
const handleForgotPassword = async () => {
    try {
        forgotLoading.value = true
        const res = await authApi.forgotPassword(forgotForm.value)
        activeTab.value = 'login'
        // 清空表单
        forgotFormRef.value?.resetFields()
    } catch (error) {
        console.error('重置密码失败:', error)
    } finally {
        forgotLoading.value = false
    }
}

onMounted(() => {
    getCaptcha()
})
</script>

<style scoped lang="less">
.login-container {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.login-wrapper {
    width: 100%;
    max-width: 1200px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: grid;
    grid-template-columns: 1fr 1fr;
    overflow: hidden;
}

.login-left {
    background: linear-gradient(135deg, var(--color-primary, #1890ff) 0%, #40a9ff 100%);
    padding: 60px 50px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    color: #fff;
}

.brand-title {
    font-size: 48px;
    font-weight: bold;
    margin-bottom: 10px;
}

.brand-subtitle {
    font-size: 18px;
    opacity: 0.9;
}

.features {
    margin-top: 40px;
}

.feature-item {
    display: flex;
    align-items: center;
    margin-bottom: 30px;
}

.feature-icon {
    font-size: 36px;
    margin-right: 20px;
}

.feature-content h3 {
    font-size: 18px;
    margin-bottom: 5px;
}

.feature-content p {
    font-size: 14px;
    opacity: 0.8;
}

.login-right {
    padding: 60px 50px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.form-wrapper {
    width: 100%;
    max-width: 400px;
}

.login-tabs {
    :deep(.ant-tabs-nav) {
        margin-bottom: 30px;
    }

    :deep(.ant-tabs-tab) {
        font-size: 16px;
        font-weight: 500;
        padding: 12px 24px;
    }

    :deep(.ant-tabs-content-holder) {
        min-height: 450px;
    }

    :deep(.ant-tabs-content) {
        transition: opacity 0.3s ease;
    }
}

.login-form {
    .form-footer {
        display: flex;
        justify-content: space-between;
        margin-top: 16px;
        font-size: 14px;
    }

    :deep(.ant-input-affix-wrapper) {
        padding-left: 12px;
    }
}

.captcha-wrapper {
    position: relative;
    width: 100%;
    height: 40px;
    cursor: pointer;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    overflow: hidden;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary, #1890ff);
    }
}

.captcha-image {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.captcha-refresh {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.02);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s;
    color: var(--color-primary, #1890ff);
    font-size: 16px;

    &:hover {
        opacity: 1;
        background: rgba(0, 0, 0, 0.05);
    }
}

.icon-spin {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

// 暗色模式适配
body.dark-mode {
    .login-container {
        background: linear-gradient(135deg, #1f1f1f 0%, #000 100%);
    }

    .login-wrapper {
        background: rgba(30, 30, 30, 0.95);
    }

    .form-footer a {
        color: var(--color-primary, #1890ff);
    }
}

// 响应式设计
@media (max-width: 768px) {
    .login-wrapper {
        grid-template-columns: 1fr;
    }

    .login-left {
        display: none;
    }

    .login-right {
        padding: 40px 30px;
    }

    .brand-title {
        font-size: 36px;
    }
}
</style>
