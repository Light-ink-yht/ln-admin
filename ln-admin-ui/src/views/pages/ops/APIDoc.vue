<template>
    <div class="api-doc-container">
        <a-card :bordered="false" class="api-doc-card">
            <template #title>
                <div class="api-doc-header">
                    <h2>接口文档</h2>
                    <a-button type="primary" @click="openSwagger">
                        <template #icon>
                            <GlobalOutlined />
                        </template>
                        在新窗口打开
                    </a-button>
                </div>
            </template>
            <div class="api-doc-content">
                <a-alert
                    message="API文档说明"
                    :description="systemConfigStore.apiDocUrl 
                        ? `接口文档地址：${systemConfigStore.apiDocUrl}。您可以在「系统配置」页面修改此地址。Swagger是一个规范和完整的框架，用于生成、描述、调用和可视化RESTful风格的Web服务。您可以通过Swagger文档查看所有可用的API接口，包括请求参数、响应格式等信息。在Swagger UI中测试接口时，请点击右上角的「Authorize」按钮，输入Bearer Token（格式：Bearer {您的token}）进行认证。`
                        : '接口文档地址未配置，使用默认地址。您可以在「系统配置」页面配置「接口文档地址」。Swagger是一个规范和完整的框架，用于生成、描述、调用和可视化RESTful风格的Web服务。您可以通过Swagger文档查看所有可用的API接口，包括请求参数、响应格式等信息。在Swagger UI中测试接口时，请点击右上角的「Authorize」按钮，输入Bearer Token（格式：Bearer {您的token}）进行认证。'"
                    type="info"
                    show-icon
                    :closable="false"
                    style="margin-bottom: 16px"
                />
                
                <!-- 跳转卡片 -->
                <a-card class="jump-card" :bordered="false">
                    <a-result
                        status="info"
                        title="接口文档"
                        sub-title="点击下方按钮跳转到接口文档页面"
                    >
                        <template #extra>
                            <a-space size="large">
                                <a-button type="primary" size="large" @click="jumpToSwagger">
                                    <template #icon>
                                        <GlobalOutlined />
                                    </template>
                                    跳转到接口文档
                                </a-button>
                                <a-button size="large" @click="openSwagger">
                                    <template #icon>
                                        <GlobalOutlined />
                                    </template>
                                    在新窗口打开
                                </a-button>
                            </a-space>
                        </template>
                    </a-result>
                    
                    <a-divider>或者</a-divider>
                    
                    <div class="url-display">
                        <div class="url-label">接口文档地址：</div>
                        <a-input
                            :value="swaggerUrl"
                            readonly
                            size="large"
                            class="url-input"
                        >
                            <template #addonAfter>
                                <a-button type="link" @click="copyUrl" :title="copied ? '已复制' : '复制地址'">
                                    <template #icon>
                                        <CopyOutlined v-if="!copied" />
                                        <CheckOutlined v-else />
                                    </template>
                                </a-button>
                            </template>
                        </a-input>
                    </div>
                </a-card>
            </div>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { GlobalOutlined, CopyOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { useSystemConfigStore } from '@/stores/modules/systemConfig'

const systemConfigStore = useSystemConfigStore()
const copied = ref(false)

const swaggerUrl = computed(() => {
    // 优先使用系统配置中的接口文档地址
    const configuredUrl = systemConfigStore.apiDocUrl
    if (configuredUrl && configuredUrl.trim()) {
        // 如果配置了地址，直接使用（不需要token，Swagger已开放访问）
        return configuredUrl.trim()
    }
    
    // 如果没有配置，使用默认的Swagger地址
    const protocol = window.location.protocol
    const host = window.location.host
    return `${protocol}//${host}/swagger/index.html`
})

// 跳转到接口文档
const jumpToSwagger = () => {
    window.location.href = swaggerUrl.value
}

// 在新窗口打开
const openSwagger = () => {
    window.open(swaggerUrl.value, '_blank')
}

// 复制URL
const copyUrl = async () => {
    try {
        await navigator.clipboard.writeText(swaggerUrl.value)
        copied.value = true
        message.success('地址已复制到剪贴板')
        setTimeout(() => {
            copied.value = false
        }, 2000)
    } catch (error) {
        message.error('复制失败，请手动复制')
    }
}
</script>

<style scoped lang="less">
.api-doc-container {
    padding: 0;
    background: transparent;
    min-height: 100%;

    .api-doc-card {
        border-radius: 12px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        transition: box-shadow 0.3s;

        &:hover {
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
        }

        :deep(.ant-card-head) {
            border-bottom: 2px solid #f0f0f0;
            padding: 12px 20px;
            background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
        }

        :deep(.ant-card-body) {
            padding: 24px;
        }
    }

    .api-doc-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        h2 {
            margin: 0;
            font-size: 20px;
            font-weight: 600;
            color: #262626;
        }
    }

    .api-doc-content {
        width: 100%;
        
        .jump-card {
            margin-top: 24px;
            border-radius: 12px;
            box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
            
            :deep(.ant-card-body) {
                padding: 40px;
            }
            
            .url-display {
                margin-top: 24px;
                
                .url-label {
                    margin-bottom: 8px;
                    font-weight: 500;
                    color: #262626;
                }
                
                .url-input {
                    :deep(.ant-input) {
                        font-family: monospace;
                        font-size: 13px;
                    }
                }
            }
        }
    }
}

// 暗色模式
body.dark-mode {
    .api-doc-container {
        .api-doc-card {
            background: #1f1f1f;
            border-color: #303030;

            :deep(.ant-card-head) {
                border-bottom-color: #303030;
                background: linear-gradient(135deg, #1f1f1f 0%, #262626 100%);
            }

            .api-doc-header {
                h2 {
                    color: #fff;
                }
            }

            .api-doc-content {
                .jump-card {
                    background: #1f1f1f;
                    border-color: #303030;
                    
                    .url-display {
                        .url-label {
                            color: rgba(255, 255, 255, 0.85);
                        }
                    }
                }
            }
        }
    }
}
</style>

