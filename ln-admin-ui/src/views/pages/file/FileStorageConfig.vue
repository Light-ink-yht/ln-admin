<template>
    <div class="storage-config-container">
        <a-card title="文件存储配置" :bordered="false">
            <a-form
                ref="formRef"
                :model="formData"
                :rules="rules"
                layout="vertical"
            >
                <!-- 存储类型选择 -->
                <div class="storage-type-section">
                    <div class="section-title">选择存储方式</div>
                    <a-radio-group v-model:value="formData.storage_type" @change="handleStorageTypeChange" class="storage-type-radio">
                        <a-radio-button value="local">
                            <div class="storage-option">
                                <div class="storage-icon">📁</div>
                                <div class="storage-info">
                                    <div class="storage-name">本地存储</div>
                                    <div class="storage-desc">文件存储在服务器本地目录</div>
                                </div>
                            </div>
                        </a-radio-button>
                        <a-radio-button value="s3">
                            <div class="storage-option">
                                <div class="storage-icon">☁️</div>
                                <div class="storage-info">
                                    <div class="storage-name">S3云存储</div>
                                    <div class="storage-desc">文件存储在S3兼容的云存储服务</div>
                                </div>
                            </div>
                        </a-radio-button>
                    </a-radio-group>
                </div>

                <!-- S3配置卡片（仅当选择S3时显示） -->
                <a-card
                    v-if="formData.storage_type === 's3'"
                    title="S3配置信息"
                    :bordered="true"
                    class="s3-config-card"
                    size="small"
                >
                    <a-form-item label="存储桶名称" name="bucket" :rules="[{ required: true, message: '请输入存储桶名称' }]">
                        <a-input v-model:value="formData.bucket" placeholder="请输入S3存储桶名称" size="large" />
                    </a-form-item>

                    <a-row :gutter="16">
                        <a-col :span="12">
                            <a-form-item label="区域" name="region" :rules="[{ required: true, message: '请输入区域' }]">
                                <a-input v-model:value="formData.region" placeholder="例如：us-east-1" size="large" />
                            </a-form-item>
                        </a-col>
                        <a-col :span="12">
                            <a-form-item label="端点地址" name="endpoint">
                                <a-input v-model:value="formData.endpoint" placeholder="可选，用于兼容服务" size="large" />
                            </a-form-item>
                        </a-col>
                    </a-row>

                    <a-row :gutter="16">
                        <a-col :span="12">
                            <a-form-item label="访问密钥ID" name="access_key_id" :rules="[{ required: true, message: '请输入访问密钥ID' }]">
                                <a-input v-model:value="formData.access_key_id" placeholder="Access Key ID" size="large" />
                            </a-form-item>
                        </a-col>
                        <a-col :span="12">
                            <a-form-item label="访问密钥" name="secret_key" :rules="[{ required: true, message: '请输入访问密钥' }]">
                                <a-input-password v-model:value="formData.secret_key" placeholder="Secret Key" size="large" />
                            </a-form-item>
                        </a-col>
                    </a-row>

                    <a-form-item label="基础访问URL" name="base_url">
                        <a-input v-model:value="formData.base_url" placeholder="可选，CDN或自定义域名" size="large" />
                        <div class="form-item-tip">
                            <a-typography-text type="secondary">如果提供，将使用此URL作为文件访问地址</a-typography-text>
                        </div>
                    </a-form-item>

                    <!-- S3启用开关 -->
                    <a-form-item label="启用状态" name="status">
                        <a-switch
                            v-model:checked="formData.statusSwitch"
                            checked-children="已启用"
                            un-checked-children="未启用"
                            @change="handleStatusChange"
                        />
                        <div class="form-item-tip">
                            <a-typography-text type="secondary">启用后，新上传的文件将使用S3存储</a-typography-text>
                        </div>
                    </a-form-item>
                </a-card>

                <!-- 操作按钮 -->
                <div class="action-section">
                    <a-space>
                        <a-button type="primary" size="large" :loading="loading" @click="handleSubmit">
                            <template #icon><SaveOutlined /></template>
                            保存配置
                        </a-button>
                        <a-button size="large" @click="handleReset">
                            <template #icon><ReloadOutlined /></template>
                            重置
                        </a-button>
                    </a-space>
                </div>
            </a-form>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { SaveOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'
import { fileApi } from '@/api/file'

const formRef = ref<FormInstance>()
const loading = ref(false)

const formData = reactive({
    storage_type: 'local',
    bucket: '',
    region: '',
    endpoint: '',
    access_key_id: '',
    secret_key: '',
    base_url: '',
    status: '1',
    statusSwitch: true, // 用于Switch组件
})

const rules = {
    bucket: [{ required: true, message: '请输入存储桶名称', trigger: 'blur' }],
    region: [{ required: true, message: '请输入区域', trigger: 'blur' }],
    access_key_id: [{ required: true, message: '请输入访问密钥ID', trigger: 'blur' }],
    secret_key: [{ required: true, message: '请输入访问密钥', trigger: 'blur' }],
}

// 存储类型改变时的处理
const handleStorageTypeChange = () => {
    // 切换存储类型时，清空S3相关字段（如果切换到本地）
    if (formData.storage_type === 'local') {
        formData.bucket = ''
        formData.region = ''
        formData.endpoint = ''
        formData.access_key_id = ''
        formData.secret_key = ''
        formData.base_url = ''
        formData.status = '1' // 本地存储默认启用
        formData.statusSwitch = true
    } else {
        // 切换到S3时，如果之前有配置，保持状态
        // 如果没有配置，默认启用
        if (formData.status === '') {
            formData.status = '1'
            formData.statusSwitch = true
        }
    }
}

// 状态开关改变
const handleStatusChange = (checked: boolean) => {
    formData.status = checked ? '1' : '2'
}

// 加载配置
const loadConfig = async () => {
    try {
        loading.value = true
        const response = await fileApi.getStorageConfig()
        if (response.code === 200 || response.code === 0) {
            if (response.data) {
                const status = response.data.status || '1'
                Object.assign(formData, {
                    storage_type: response.data.storage_type || 'local',
                    bucket: response.data.bucket || '',
                    region: response.data.region || '',
                    endpoint: response.data.endpoint || '',
                    access_key_id: response.data.access_key_id || '',
                    secret_key: '', // 密钥不返回，保持为空
                    base_url: response.data.base_url || '',
                    status: status,
                    statusSwitch: status === '1',
                })
            }
        } else {
            message.error(response.msg || '获取配置失败')
        }
    } catch (error: any) {
        console.error('获取配置失败:', error)
        message.error(error.message || '获取配置失败')
    } finally {
        loading.value = false
    }
}

// 提交表单
const handleSubmit = async () => {
    try {
        // 如果是S3存储，需要验证必填字段
        if (formData.storage_type === 's3') {
            await formRef.value?.validate()
        }

        loading.value = true

        // 如果是本地存储，只需要存储类型，状态固定为启用
        if (formData.storage_type === 'local') {
            const submitData = {
                storage_type: 'local',
                status: '1', // 本地存储总是启用
            }
            const response = await fileApi.saveStorageConfig(submitData)
            if (response.code === 200 || response.code === 0) {
                message.success('保存配置成功')
                await loadConfig()
            } else {
                message.error(response.msg || '保存配置失败')
            }
        } else {
            // S3存储需要所有字段
            const submitData = {
                storage_type: formData.storage_type,
                bucket: formData.bucket,
                region: formData.region,
                endpoint: formData.endpoint || '',
                access_key_id: formData.access_key_id,
                secret_key: formData.secret_key,
                base_url: formData.base_url || '',
                status: formData.status,
            }
            const response = await fileApi.saveStorageConfig(submitData)
            if (response.code === 200 || response.code === 0) {
                message.success('保存配置成功')
                await loadConfig()
            } else {
                message.error(response.msg || '保存配置失败')
            }
        }
    } catch (error: any) {
        console.error('保存配置失败:', error)
        if (error.errorFields) {
            return
        }
        message.error(error.message || '保存配置失败')
    } finally {
        loading.value = false
    }
}

// 重置表单
const handleReset = () => {
    formRef.value?.resetFields()
    loadConfig()
}

onMounted(() => {
    loadConfig()
})
</script>

<style scoped>
.storage-config-container {
    padding: 0;
}

.storage-type-section {
    margin-bottom: 24px;
}

.section-title {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 16px;
    color: rgba(0, 0, 0, 0.85);
}

.storage-type-radio {
    width: 100%;
    display: flex;
    gap: 16px;
}

.storage-type-radio :deep(.ant-radio-button-wrapper) {
    flex: 1;
    height: auto;
    padding: 20px;
    border-radius: 8px;
    border: 2px solid #d9d9d9;
    transition: all 0.3s;
}

.storage-type-radio :deep(.ant-radio-button-wrapper:hover) {
    border-color: #1890ff;
}

.storage-type-radio :deep(.ant-radio-button-wrapper-checked) {
    border-color: #1890ff;
    background: #e6f7ff;
    box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.storage-option {
    display: flex;
    align-items: center;
    gap: 12px;
}

.storage-icon {
    font-size: 32px;
    line-height: 1;
}

.storage-info {
    flex: 1;
}

.storage-name {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 4px;
    color: rgba(0, 0, 0, 0.85);
}

.storage-desc {
    font-size: 13px;
    color: rgba(0, 0, 0, 0.45);
}

.s3-config-card {
    margin-top: 24px;
    margin-bottom: 24px;
}

.s3-config-card :deep(.ant-card-head) {
    background: #fafafa;
    border-bottom: 1px solid #f0f0f0;
}

.s3-config-card :deep(.ant-card-body) {
    padding: 24px;
}

.action-section {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #f0f0f0;
    text-align: center;
}

.form-item-tip {
    margin-top: 4px;
    font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .storage-type-radio {
        flex-direction: column;
    }

    .storage-type-radio :deep(.ant-radio-button-wrapper) {
        width: 100%;
    }
}
</style>

