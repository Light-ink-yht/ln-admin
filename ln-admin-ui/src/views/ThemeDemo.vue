<template>
    <div class="theme-demo">
        <a-page-header title="主题演示" sub-title="动态主题系统和功能展示"
            :breadcrumb="{ items: [{ title: '首页' }, { title: '主题演示' }] }" />

        <div class="demo-container">
            <!-- 主题色选择 -->
            <a-card title="主题色设置" class="demo-card">
                <template #extra>
                    <a-tag :color="themeStore.colorPrimary">当前色: {{ themeStore.colorPrimary }}</a-tag>
                </template>

                <a-space wrap>
                    <div v-for="color in PRESET_COLORS" :key="color" class="color-item"
                        :class="{ active: themeStore.colorPrimary === color }" :style="{ backgroundColor: color }"
                        @click="handleColorChange(color)">
                        <check-outlined v-if="themeStore.colorPrimary === color" class="check-icon" />
                        <span class="color-name">{{ COLOR_NAMES[color] }}</span>
                    </div>
                </a-space>

                <a-divider>自定义颜色</a-divider>

                <a-space>
                    <a-color-picker v-model:value="customColor" show-text @change="handleCustomColorChange" />
                    <a-button type="primary" @click="applyCustomColor">应用自定义颜色</a-button>
                </a-space>
            </a-card>

            <!-- 主题模式切换 -->
            <a-card title="主题模式" class="demo-card">
                <a-space size="large">
                    <a-switch :checked="themeStore.isDark" @change="handleModeToggle" checked-children="暗色"
                        un-checked-children="亮色" />
                    <a-button @click="handleFollowSystem">跟随系统主题</a-button>
                    <a-tag :color="themeStore.isDark ? 'purple' : 'blue'">
                        当前模式: {{ themeStore.isDark ? '暗色' : '亮色' }}
                    </a-tag>
                </a-space>
            </a-card>

            <!-- 组件展示 -->
            <a-card title="组件样式展示" class="demo-card">
                <a-space direction="vertical" :size="20" style="width: 100%">
                    <!-- 按钮组 -->
                    <div>
                        <h4>按钮</h4>
                        <a-space wrap>
                            <a-button type="primary">主要按钮</a-button>
                            <a-button type="default">默认按钮</a-button>
                            <a-button type="dashed">虚线按钮</a-button>
                            <a-button type="link">链接按钮</a-button>
                            <a-button danger>危险按钮</a-button>
                        </a-space>
                    </div>

                    <!-- 标签 -->
                    <div>
                        <h4>标签</h4>
                        <a-space wrap>
                            <a-tag color="success">成功</a-tag>
                            <a-tag color="processing">进行中</a-tag>
                            <a-tag color="error">错误</a-tag>
                            <a-tag color="warning">警告</a-tag>
                            <a-tag>默认标签</a-tag>
                        </a-space>
                    </div>

                    <!-- 输入框 -->
                    <div>
                        <h4>输入框</h4>
                        <a-space direction="vertical" style="width: 100%">
                            <a-input placeholder="请输入内容" />
                            <a-input placeholder="禁用状态" disabled />
                            <a-textarea placeholder="多行输入" :rows="4" />
                        </a-space>
                    </div>

                    <!-- 选择器 -->
                    <div>
                        <h4>选择器</h4>
                        <a-space wrap>
                            <a-select v-model:value="selectValue" style="width: 200px" placeholder="请选择">
                                <a-select-option value="option1">选项 1</a-select-option>
                                <a-select-option value="option2">选项 2</a-select-option>
                                <a-select-option value="option3">选项 3</a-select-option>
                            </a-select>
                            <a-date-picker v-model:value="dateValue" />
                            <a-time-picker v-model:value="timeValue" />
                        </a-space>
                    </div>

                    <!-- 消息提示 -->
                    <div>
                        <h4>消息提示</h4>
                        <a-space wrap>
                            <a-button @click="showSuccess">成功消息</a-button>
                            <a-button @click="showInfo">信息消息</a-button>
                            <a-button @click="showWarning">警告消息</a-button>
                            <a-button danger @click="showError">错误消息</a-button>
                        </a-space>
                    </div>

                    <!-- 进度条 -->
                    <div>
                        <h4>进度条</h4>
                        <a-space direction="vertical" style="width: 100%">
                            <a-progress :percent="30" />
                            <a-progress :percent="50" status="active" />
                            <a-progress :percent="70" :status="'exception'" />
                            <a-progress :percent="100" />
                        </a-space>
                    </div>

                    <!-- 开关 -->
                    <div>
                        <h4>开关</h4>
                        <a-space>
                            <a-switch v-model:checked="switchValue1" />
                            <a-switch v-model:checked="switchValue2" checked-children="开" un-checked-children="关" />
                            <a-switch v-model:checked="switchValue3" :loading="true" />
                        </a-space>
                    </div>

                    <!-- 表格 -->
                    <div>
                        <h4>表格</h4>
                        <a-table :columns="tableColumns" :data-source="tableData" :pagination="false" size="small" />
                    </div>

                    <!-- 卡片 -->
                    <div>
                        <h4>卡片</h4>
                        <a-row :gutter="16">
                            <a-col :span="8">
                                <a-card>
                                    <template #title>卡片标题 1</template>
                                    <template #extra><a href="#">更多</a></template>
                                    <p>这是卡片的内容</p>
                                </a-card>
                            </a-col>
                            <a-col :span="8">
                                <a-card title="卡片标题 2">
                                    <p>这是卡片的内容</p>
                                    <p>支持多行内容</p>
                                </a-card>
                            </a-col>
                            <a-col :span="8">
                                <a-card>
                                    <p>这是卡片的内容</p>
                                </a-card>
                            </a-col>
                        </a-row>
                    </div>
                </a-space>
            </a-card>

            <!-- 主题信息 -->
            <a-card title="主题信息" class="demo-card">
                <a-descriptions :column="2" bordered>
                    <a-descriptions-item label="当前主题色">
                        <a-tag :color="themeStore.colorPrimary">{{ themeStore.colorPrimary }}</a-tag>
                    </a-descriptions-item>
                    <a-descriptions-item label="主题模式">
                        {{ themeStore.isDark ? '🌙 暗色模式' : '☀️ 亮色模式' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="持久化状态">
                        <a-tag color="success">已启用</a-tag>
                    </a-descriptions-item>
                    <a-descriptions-item label="Ant Design 版本">
                        4.2.6
                    </a-descriptions-item>
                </a-descriptions>
            </a-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { CheckOutlined } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { PRESET_COLORS, COLOR_NAMES, ThemeMode } from '@/theme/config'

const themeStore = useThemeStore()

const customColor = ref('#1890ff')
const selectValue = ref(undefined)
const dateValue = ref(undefined)
const timeValue = ref(undefined)
const switchValue1 = ref(true)
const switchValue2 = ref(false)
const switchValue3 = ref(false)

// 表格数据
const tableColumns = [
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '年龄', dataIndex: 'age', key: 'age' },
    { title: '地址', dataIndex: 'address', key: 'address' },
    { title: '操作', key: 'action' },
]

const tableData = ref([
    { key: '1', name: '张三', age: 32, address: '北京市朝阳区' },
    { key: '2', name: '李四', age: 42, address: '上海市浦东新区' },
    { key: '3', name: '王五', age: 32, address: '广州市天河区' },
])

const handleColorChange = (color: string) => {
    themeStore.setColorPrimary(color)
    message.success(`已切换到${COLOR_NAMES[color]}`)
}

const handleCustomColorChange = (color: any) => {
    customColor.value = color.toHexString()
}

const applyCustomColor = () => {
    themeStore.setColorPrimary(customColor.value)
    message.success(`已应用自定义颜色: ${customColor.value}`)
}

const handleModeToggle = (checked: boolean) => {
    themeStore.setThemeMode(checked ? ThemeMode.DARK : ThemeMode.LIGHT)
    message.success(`已切换到${checked ? '暗色' : '亮色'}模式`)
}

const handleFollowSystem = () => {
    themeStore.followSystemTheme()
    message.success('已跟随系统主题')
}

const showSuccess = () => {
    message.success('这是一条成功消息')
}

const showInfo = () => {
    message.info('这是一条信息消息')
}

const showWarning = () => {
    message.warning('这是一条警告消息')
}

const showError = () => {
    message.error('这是一条错误消息')
}
</script>

<style scoped lang="less">
.theme-demo {
    padding: 24px;

    .demo-container {
        display: flex;
        flex-direction: column;
        gap: 24px;
        margin-top: 24px;
    }

    .demo-card {
        :deep(.ant-card-head) {
            background: transparent;
        }
    }
}

.color-item {
    position: relative;
    width: 80px;
    height: 80px;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

    &:hover {
        transform: scale(1.1);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
    }

    &.active {
        border: 3px solid #fff;
        box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.3);
    }

    .check-icon {
        font-size: 24px;
        color: #fff;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
    }

    .color-name {
        position: absolute;
        bottom: 4px;
        left: 50%;
        transform: translateX(-50%);
        font-size: 12px;
        color: #fff;
        font-weight: 500;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
    }
}
</style>
