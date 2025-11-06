<template>
    <div class="system-config-container">
        <div class="page-header-banner">
            <div class="banner-content">
                <div class="banner-left">
                    <div class="banner-icon-wrapper">
                        <SettingOutlined class="banner-icon" />
                    </div>
                    <div class="banner-text">
                        <h1 class="banner-title">系统配置</h1>
                        <p class="banner-subtitle">管理网站基本信息、联系方式等配置项</p>
                    </div>
                </div>
                <div class="banner-decorations">
                    <div class="decoration-circle circle-1"></div>
                    <div class="decoration-circle circle-2"></div>
                    <div class="decoration-circle circle-3"></div>
                </div>
            </div>
        </div>

        <a-card :bordered="false" class="page-card">
            <a-spin :spinning="loading">
                <div v-if="systemConfigs && systemConfigs.length > 0" class="config-content">
                    <!-- 基本信息卡片 -->
                    <div class="config-group">
                        <div class="group-header">
                            <div class="group-title">
                                <GlobalOutlined class="group-icon" />
                                <span>基本信息</span>
                            </div>
                            <a-button
                                type="primary"
                                size="large"
                                :loading="savingGroup['basic']"
                                @click="handleSaveGroup('basic')"
                            >
                                <template #icon>
                                    <SaveOutlined />
                                </template>
                                保存本组
                            </a-button>
                        </div>
                        
                        <!-- 第一行：网站名称 -->
                        <div class="config-row">
                            <div class="config-item config-item-full">
                                <div class="config-label">
                                    <span class="label-text">网站名称</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_name')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-input
                                    v-model:value="formData.site_name"
                                    placeholder="请输入网站名称"
                                    size="large"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_name')?.description }}</div>
                            </div>
                        </div>

                        <!-- 第二行：网站LOGO和图标 -->
                        <div class="config-row">
                            <div class="config-item">
                                <div class="config-label">
                                    <span class="label-text">网站LOGO</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_logo')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <ImageUploadSingle
                                    v-model="formData.site_logo"
                                    :max-size="5"
                                    placeholder="上传LOGO"
                                    preview-alt="网站LOGO"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_logo')?.description }}</div>
                            </div>

                            <div class="config-item">
                                <div class="config-label">
                                    <span class="label-text">网站图标</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_favicon')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <ImageUploadSingle
                                    v-model="formData.site_favicon"
                                    :max-size="2"
                                    placeholder="上传图标"
                                    preview-alt="网站图标"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_favicon')?.description }}</div>
                            </div>
                        </div>

                        <!-- 第三行：版权信息 -->
                        <div class="config-row">
                            <div class="config-item config-item-full">
                                <div class="config-label">
                                    <span class="label-text">版权信息</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_copyright')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-textarea
                                    v-model:value="formData.site_copyright"
                                    placeholder="请输入版权信息"
                                    :rows="3"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_copyright')?.description }}</div>
                            </div>
                        </div>

                        <!-- 第四行：网站描述和关键词 -->
                        <div class="config-row">
                            <div class="config-item">
                                <div class="config-label">
                                    <span class="label-text">网站描述</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_description')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-textarea
                                    v-model:value="formData.site_description"
                                    placeholder="请输入网站描述"
                                    :rows="3"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_description')?.description }}</div>
                            </div>

                            <div class="config-item">
                                <div class="config-label">
                                    <span class="label-text">网站关键词</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_keywords')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-input
                                    v-model:value="formData.site_keywords"
                                    placeholder="请输入网站关键词，多个用逗号分隔"
                                    size="large"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_keywords')?.description }}</div>
                            </div>
                        </div>

                        <!-- 第五行：备案号 -->
                        <div class="config-row">
                            <div class="config-item config-item-full">
                                <div class="config-label">
                                    <span class="label-text">备案号</span>
                                    <a-tag v-if="basicConfigs.find(c => c.configKey === 'site_beian')?.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-input
                                    v-model:value="formData.site_beian"
                                    placeholder="请输入备案号"
                                    size="large"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ basicConfigs.find(c => c.configKey === 'site_beian')?.description }}</div>
                            </div>
                        </div>
                    </div>

                    <!-- 联系信息卡片 -->
                    <div class="config-group" v-if="contactConfigs.length > 0">
                        <div class="group-header">
                            <div class="group-title">
                                <PhoneOutlined class="group-icon" />
                                <span>联系信息</span>
                            </div>
                            <a-button
                                type="primary"
                                size="large"
                                :loading="savingGroup['contact']"
                                @click="handleSaveGroup('contact')"
                            >
                                <template #icon>
                                    <SaveOutlined />
                                </template>
                                保存本组
                            </a-button>
                        </div>
                        <div class="config-row">
                            <div class="config-item" v-for="config in contactConfigs" :key="config.configKey" :class="{ 'config-item-full': isLongText(config.configKey) }">
                                <div class="config-label">
                                    <span class="label-text">{{ config.configName }}</span>
                                    <a-tag v-if="config.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-input
                                    v-if="!isLongText(config.configKey)"
                                    v-model:value="formData[config.configKey]"
                                    :placeholder="'请输入' + config.configName"
                                    size="large"
                                    allow-clear
                                    class="config-input"
                                />
                                <a-textarea
                                    v-else
                                    v-model:value="formData[config.configKey]"
                                    :placeholder="'请输入' + config.configName"
                                    :rows="3"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ config.description }}</div>
                            </div>
                        </div>
                    </div>

                    <!-- 其他配置卡片 -->
                    <div class="config-group" v-if="otherConfigs.length > 0">
                        <div class="group-header">
                            <div class="group-title">
                                <SettingOutlined class="group-icon" />
                                <span>其他配置</span>
                            </div>
                            <a-button
                                type="primary"
                                :loading="savingGroup['other']"
                                @click="handleSaveGroup('other')"
                            >
                                <template #icon>
                                    <SaveOutlined />
                                </template>
                                保存本组
                            </a-button>
                        </div>
                        <div class="config-grid">
                            <div class="config-item" v-for="config in otherConfigs" :key="config.configKey">
                                <div class="config-label">
                                    <span class="label-text">{{ config.configName }}</span>
                                    <a-tag v-if="config.status === '1'" color="success" size="small">启用</a-tag>
                                    <a-tag v-else color="error" size="small">禁用</a-tag>
                                </div>
                                <a-input
                                    v-if="!isLongText(config.configKey)"
                                    v-model:value="formData[config.configKey]"
                                    :placeholder="'请输入' + config.configName"
                                    size="large"
                                    allow-clear
                                    class="config-input"
                                />
                                <a-textarea
                                    v-else
                                    v-model:value="formData[config.configKey]"
                                    :placeholder="'请输入' + config.configName"
                                    :rows="3"
                                    allow-clear
                                    class="config-input"
                                />
                                <div class="config-desc">{{ config.description }}</div>
                            </div>
                        </div>
                    </div>
                </div>
                <a-empty v-else description="暂无配置" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
            </a-spin>
        </a-card>

        <!-- 编辑/添加配置弹窗 -->
        <a-modal
            v-model:open="editModalOpen"
            :title="editingConfig ? '编辑配置' : '添加配置'"
            :confirm-loading="saving"
            width="600px"
            @ok="handleSave"
            @cancel="handleCancel"
        >
            <a-form
                :model="editForm"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 18 }"
                :rules="editRules"
            >
                <a-form-item label="配置名称" name="configName" required>
                    <a-input
                        v-model:value="editForm.configName"
                        :disabled="!!editingConfig"
                        placeholder="请输入配置名称"
                    />
                </a-form-item>
                <a-form-item label="配置键" name="configKey" required>
                    <a-input
                        v-model:value="editForm.configKey"
                        :disabled="!!editingConfig"
                        placeholder="请输入配置键（英文，如：site_name）"
                    />
                </a-form-item>
                <a-form-item label="配置值" name="configValue" required>
                    <a-textarea
                        v-model:value="editForm.configValue"
                        :placeholder="'请输入' + editForm.configName"
                        :rows="4"
                        allow-clear
                    />
                </a-form-item>
                <a-form-item label="配置分组" name="configGroup" required>
                    <a-select v-model:value="editForm.configGroup" placeholder="请选择配置分组">
                        <a-select-option value="system">系统配置</a-select-option>
                        <a-select-option value="sms">短信配置</a-select-option>
                        <a-select-option value="jwt">JWT配置</a-select-option>
                        <a-select-option value="database">数据库配置</a-select-option>
                        <a-select-option value="redis">Redis配置</a-select-option>
                        <a-select-option value="other">其他配置</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="配置描述" name="description">
                    <a-textarea
                        v-model:value="editForm.description"
                        placeholder="请输入配置描述"
                        :rows="3"
                        allow-clear
                    />
                </a-form-item>
                <a-form-item label="状态" name="status" required>
                    <a-select v-model:value="editForm.status">
                        <a-select-option value="1">启用</a-select-option>
                        <a-select-option value="2">禁用</a-select-option>
                    </a-select>
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, reactive, h } from 'vue'
import { message, Empty } from 'ant-design-vue'
import {
    ReloadOutlined,
    PlusOutlined,
    SaveOutlined,
    SettingOutlined,
    GlobalOutlined,
    PhoneOutlined,
} from '@ant-design/icons-vue'
import { systemApi, type SystemConfig, type UpdateSystemConfigRequest, type CreateSystemConfigRequest } from '@/api/system'
import { useSystemConfigStore } from '@/stores/modules/systemConfig'
import { ImageUploadSingle } from '@/components/upload'

// 基本信息配置键
const basicConfigKeys = [
    'site_name',
    'site_logo',
    'site_favicon',
    'site_copyright',
    'site_description',
    'site_keywords',
    'site_beian',
]

// 联系信息配置键
const contactConfigKeys = ['site_contact_email', 'site_contact_phone', 'site_address']

const loading = ref(false)
const saving = ref(false)
const savingGroup = reactive<Record<string, boolean>>({})
const systemConfigs = ref<SystemConfig[]>([])
const editModalOpen = ref(false)
const editingConfig = ref<SystemConfig | null>(null)
const formData = reactive<Record<string, string>>({})
const systemConfigStore = useSystemConfigStore()
const editForm = ref<CreateSystemConfigRequest & { configKey: string; configName: string; configGroup: string; description?: string; status?: string }>({
    configKey: '',
    configName: '',
    configGroup: 'system',
    configValue: '',
    description: '',
    status: '1',
})

const editRules = {
    configName: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
    configKey: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
    configValue: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
    configGroup: [{ required: true, message: '请选择配置分组', trigger: 'change' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

// 判断是否为长文本配置
const isLongText = (key: string): boolean => {
    const longTextKeys = ['site_description', 'site_copyright', 'site_address']
    return longTextKeys.includes(key)
}


// 基本信息配置
const basicConfigs = computed(() => {
    return systemConfigs.value.filter((config) => basicConfigKeys.includes(config.configKey))
})

// 联系信息配置
const contactConfigs = computed(() => {
    return systemConfigs.value.filter((config) => contactConfigKeys.includes(config.configKey))
})

// 其他配置
const otherConfigs = computed(() => {
    return systemConfigs.value.filter(
        (config) => !basicConfigKeys.includes(config.configKey) && !contactConfigKeys.includes(config.configKey)
    )
})

// 加载配置列表
const loadConfigs = async () => {
    loading.value = true
    try {
        const response = await systemApi.getAllConfigs()
        if (response.code === 200 || response.code === 0) {
            systemConfigs.value = response.data || []
            // 初始化表单数据
            systemConfigs.value.forEach((config) => {
                formData[config.configKey] = config.configValue
            })
        }
    } catch (error: any) {
        console.error('加载配置列表失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        loading.value = false
    }
}

// 保存分组配置
const handleSaveGroup = async (group: string) => {
    let configsToSave: SystemConfig[] = []

    if (group === 'basic') {
        configsToSave = basicConfigs.value
    } else if (group === 'contact') {
        configsToSave = contactConfigs.value
    } else if (group === 'other') {
        configsToSave = otherConfigs.value
    }

    if (configsToSave.length === 0) {
        message.info('该组没有可保存的配置')
        return
    }

    savingGroup[group] = true
    let successCount = 0
    let failCount = 0

    try {
        for (const config of configsToSave) {
            const newValue = formData[config.configKey]
            if (newValue === undefined || newValue === config.configValue) {
                continue
            }

            try {
                const updateData: UpdateSystemConfigRequest = {
                    configValue: newValue,
                }
                await systemApi.updateConfig(config.configKey, updateData)
                // 更新本地配置
                const idx = systemConfigs.value.findIndex((c) => c.configKey === config.configKey)
                if (idx !== -1) {
                    systemConfigs.value[idx].configValue = newValue
                }
                // 更新store
                systemConfigStore.updateConfig(config.configKey, newValue)
                successCount++
            } catch (error) {
                failCount++
                console.error(`更新配置 ${config.configName} 失败:`, error)
            }
        }

        if (successCount > 0) {
            message.success(`成功保存 ${successCount} 个配置`)
        }
        if (failCount > 0) {
            message.warning(`有 ${failCount} 个配置保存失败`)
        }
        if (successCount === 0 && failCount === 0) {
            message.info('配置值未发生变化')
        }
    } finally {
        savingGroup[group] = false
    }
}

// 添加配置
const handleAdd = () => {
    editingConfig.value = null
    editForm.value = {
        configKey: '',
        configName: '',
        configGroup: 'system',
        configValue: '',
        description: '',
        status: '1',
    }
    editModalOpen.value = true
}

// 保存配置（弹窗）
const handleSave = async () => {
    if (!editForm.value.configName || !editForm.value.configKey || !editForm.value.configValue) {
        message.warning('请填写完整信息')
        return
    }

    saving.value = true
    try {
        if (editingConfig.value) {
            // 更新配置
            const updateData: UpdateSystemConfigRequest = {
                configValue: editForm.value.configValue,
                description: editForm.value.description,
                status: editForm.value.status,
            }
            await systemApi.updateConfig(editForm.value.configKey, updateData)
            message.success('更新配置成功')
            systemConfigStore.updateConfig(editForm.value.configKey, editForm.value.configValue)
        } else {
            // 创建配置
            const createData: CreateSystemConfigRequest = {
                configKey: editForm.value.configKey,
                configName: editForm.value.configName,
                configGroup: editForm.value.configGroup,
                configValue: editForm.value.configValue,
                description: editForm.value.description,
                status: editForm.value.status || '1',
            }
            await systemApi.createConfig(createData)
            message.success('创建配置成功')
        }
        editModalOpen.value = false
        loadConfigs()
        await systemConfigStore.loadConfigs()
    } catch (error: any) {
        console.error('保存配置失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        saving.value = false
    }
}

// 取消编辑
const handleCancel = () => {
    editModalOpen.value = false
    editingConfig.value = null
    editForm.value = {
        configKey: '',
        configName: '',
        configGroup: 'system',
        configValue: '',
        description: '',
        status: '1',
    }
}

onMounted(() => {
    loadConfigs()
})
</script>

<style scoped lang="less">
.system-config-container {
    padding: 0;
    background: transparent;
    min-height: 100%;

    .page-header-banner {
        background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-active) 100%);
        border-radius: 16px;
        padding: 32px 40px;
        margin-bottom: 24px;
        position: relative;
        overflow: hidden;
        box-shadow: 0 8px 32px rgba(var(--color-primary-rgb, 24, 144, 255), 0.25);

        .banner-content {
            position: relative;
            z-index: 2;
        }

        .banner-left {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .banner-icon-wrapper {
            width: 72px;
            height: 72px;
            background: rgba(255, 255, 255, 0.25);
            border-radius: 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.3);
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
            transition: all 0.3s;

            &:hover {
                transform: scale(1.05) rotate(5deg);
                background: rgba(255, 255, 255, 0.3);
            }

            .banner-icon {
                font-size: 32px;
                color: #fff;
            }
        }

        .banner-text {
            flex: 1;

            .banner-title {
                margin: 0 0 8px 0;
                font-size: 32px;
                font-weight: 700;
                color: #fff;
                line-height: 1.2;
                text-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
                letter-spacing: -0.5px;
            }

            .banner-subtitle {
                margin: 0;
                font-size: 15px;
                color: rgba(255, 255, 255, 0.9);
                line-height: 1.5;
                font-weight: 400;
            }
        }

        .banner-decorations {
            position: absolute;
            top: 0;
            right: 0;
            width: 100%;
            height: 100%;
            pointer-events: none;
            z-index: 1;

            .decoration-circle {
                position: absolute;
                border-radius: 50%;
                background: rgba(255, 255, 255, 0.1);
                backdrop-filter: blur(20px);
                animation: float 6s ease-in-out infinite;

                &.circle-1 {
                    width: 120px;
                    height: 120px;
                    top: -40px;
                    right: 80px;
                    animation-delay: 0s;
                }

                &.circle-2 {
                    width: 80px;
                    height: 80px;
                    top: 20px;
                    right: 20px;
                    animation-delay: 2s;
                }

                &.circle-3 {
                    width: 60px;
                    height: 60px;
                    bottom: -20px;
                    right: 120px;
                    animation-delay: 4s;
                }
            }
        }
    }

    @keyframes float {
        0%,
        100% {
            transform: translateY(0) translateX(0);
        }
        50% {
            transform: translateY(-20px) translateX(10px);
        }
    }

    .page-card {
        border-radius: 16px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
        overflow: hidden;
        background: #fff;

        :deep(.ant-card-head) {
            display: none;
        }

        :deep(.ant-card-body) {
            padding: 24px;
            background: #f8f9fa;
        }
    }

    .logo-preview-section {
        display: flex;
        gap: 24px;
        padding: 24px;
        margin-bottom: 24px;
        background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    .preview-item {
        flex: 1;
        background: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    .preview-label {
        font-size: 14px;
        font-weight: 500;
        color: rgba(0, 0, 0, 0.85);
        margin-bottom: 12px;
    }

    .preview-content {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 80px;
        background: #fafafa;
        border-radius: 8px;
        padding: 16px;
    }

    .logo-preview {
        max-width: 200px;
        max-height: 60px;
        object-fit: contain;
    }

    .favicon-preview {
        width: 32px;
        height: 32px;
        object-fit: contain;
    }

    .config-content {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .config-group {
        background: #fff;
        border-radius: 12px;
        padding: 24px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
        transition: all 0.3s;

        &:hover {
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
        }

        .group-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 24px;
            padding-bottom: 16px;
            border-bottom: 2px solid #f0f0f0;

            .group-title {
                display: flex;
                align-items: center;
                gap: 12px;
                font-size: 18px;
                font-weight: 600;
                color: #262626;

                .group-icon {
                    font-size: 20px;
                    color: var(--color-primary);
                }
            }
        }

        .config-row {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 20px;
            margin-bottom: 20px;

            &:last-child {
                margin-bottom: 0;
            }

            @media (max-width: 1200px) {
                grid-template-columns: 1fr;
            }

            @media (max-width: 768px) {
                grid-template-columns: 1fr;
            }
        }

        .config-item {
            background: #fff;
            border: 1px solid #e8e8e8;
            border-radius: 12px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            min-height: 160px;
            position: relative;
            overflow: hidden;

            &::before {
                content: '';
                position: absolute;
                top: 0;
                left: 0;
                width: 4px;
                height: 100%;
                background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-active) 100%);
                opacity: 0;
                transition: opacity 0.3s;
            }

            &:hover {
                border-color: var(--color-primary);
                box-shadow: 0 6px 20px rgba(var(--color-primary-rgb, 24, 144, 255), 0.2);
                transform: translateY(-4px);

                &::before {
                    opacity: 1;
                }
            }

            &.config-item-full {
                grid-column: 1 / -1;
            }
        }


        .image-config-wrapper {
            width: 100%;
            transition: all 0.3s;
            border: 1px solid #e8e8e8;

            &:hover {
                background: #fff;
                border-color: var(--color-primary);
                box-shadow: 0 2px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.1);
                transform: translateY(-2px);
            }

        .config-label {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
            padding-bottom: 8px;
            border-bottom: 1px solid #f0f0f0;

            .label-text {
                font-size: 15px;
                font-weight: 600;
                color: #262626;
            }
        }

        .config-input {
            margin-bottom: 8px;
            border-radius: 6px;
            transition: all 0.3s;

            &:hover {
                border-color: var(--color-primary);
            }

            &:focus,
            &:focus-within {
                border-color: var(--color-primary);
                box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb, 24, 144, 255), 0.1);
            }
        }

        .config-desc {
            font-size: 12px;
            color: #8c8c8c;
            line-height: 1.6;
            margin-top: auto;
            padding-top: 8px;
        }
        }
    }
}

/* 深色模式 */
body.dark-mode {
    .system-config-container {
        .page-header-banner {
            background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-active) 100%);
            box-shadow: 0 8px 32px rgba(var(--color-primary-rgb, 24, 144, 255), 0.3);
        }

        .page-card {
            :deep(.ant-card-body) {
                background: #141414;
            }
        }

        .config-group {
            background: #1f1f1f;
            border: 1px solid #434343;

            &:hover {
                box-shadow: 0 6px 24px rgba(0, 0, 0, 0.4);
                border-color: #595959;
            }

            .group-header {
                border-bottom-color: #434343;

                .group-title {
                    color: rgba(255, 255, 255, 0.85);

                    .group-icon {
                        color: var(--color-primary);
                    }
                }
            }

            .config-item {
                background: #141414;
                border-color: #434343;

                &:hover {
                    background: #1f1f1f;
                    border-color: var(--color-primary);
                    box-shadow: 0 6px 20px rgba(var(--color-primary-rgb, 24, 144, 255), 0.15);
                }

                &::before {
                    background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-active) 100%);
                }

                .config-label {
                    border-bottom-color: #434343;

                    .label-text {
                        color: rgba(255, 255, 255, 0.85);
                    }
                }

                .config-desc {
                    color: rgba(255, 255, 255, 0.45);
                }

                .config-input {
                    background: #1f1f1f;
                    border-color: #434343;
                    color: rgba(255, 255, 255, 0.85);

                    &:hover {
                        border-color: var(--color-primary-hover);
                    }

                    &:focus,
                    &:focus-within {
                        border-color: var(--color-primary);
                        background: #262626;
                        box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb, 24, 144, 255), 0.2);
                    }

                    &::placeholder {
                        color: rgba(255, 255, 255, 0.25);
                    }
                }

            }
        }
    }
}
</style>
