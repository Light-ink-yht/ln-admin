<template>
    <div class="theme-switch-container">
        <!-- 主题模式切换 -->
        <a-dropdown placement="bottomRight" :trigger="['click']">
            <a-button type="text" class="theme-switch" :style="themeButtonStyle">
                <template #icon>
                    <BulbFilled v-if="!isDark" />
                    <BulbOutlined v-else />
                </template>
            </a-button>
            <template #overlay>
                <a-menu @click="handleThemeChange">
                    <a-menu-item key="light" :disabled="!isDark">
                        <BulbFilled />
                        <span>浅色模式</span>
                    </a-menu-item>
                    <a-menu-item key="dark" :disabled="isDark">
                        <BulbOutlined />
                        <span>深色模式</span>
                    </a-menu-item>
                    <a-menu-item key="system">
                        <DesktopOutlined />
                        <span>跟随系统</span>
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="color" @click.stop="colorPickerVisible = true">
                        <SettingOutlined />
                        <span>主题色</span>
                    </a-menu-item>
                </a-menu>
            </template>
        </a-dropdown>
        
        <!-- 颜色选择器弹窗 -->
        <a-modal
            v-model:open="colorPickerVisible"
            title="选择主题色"
            :footer="null"
            width="320px"
            :mask-closable="true"
        >
            <div class="color-picker-menu">
                <div class="preset-colors">
                    <div
                        v-for="color in presetColors"
                        :key="color.value"
                        class="color-item"
                        :class="{ active: colorPrimary === color.value }"
                        :style="{ backgroundColor: color.value }"
                        @click="handleColorChange(color.value)"
                    >
                        <CheckOutlined v-if="colorPrimary === color.value" class="color-check" />
                    </div>
                </div>
                <div class="custom-color-picker">
                    <div class="custom-color-label">自定义颜色：</div>
                    <a-color-picker
                        v-model:value="customColor"
                        :show-text="true"
                        size="small"
                        @change="handleCustomColorChange"
                    />
                </div>
            </div>
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { BulbOutlined, BulbFilled, DesktopOutlined, CheckOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { ThemeMode, PRESET_COLORS, COLOR_NAMES } from '@/theme/config'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.isDark)
const colorPrimary = computed(() => themeStore.colorPrimary)
const customColor = ref<string>(themeStore.colorPrimary)
const colorPickerVisible = ref(false)

// 预设颜色列表
const presetColors = computed(() => {
    return PRESET_COLORS.map((color) => ({
        value: color,
        name: COLOR_NAMES[color] || color,
    }))
})

const themeButtonStyle = computed(() => {
    return {
        color: isDark.value ? '#fff' : 'rgba(0, 0, 0, 0.85)',
    }
})

const handleThemeChange = ({ key }: { key: string }) => {
    if (key === 'light') {
        themeStore.setThemeMode(ThemeMode.LIGHT)
        message.success('已切换为浅色模式')
    } else if (key === 'dark') {
        themeStore.setThemeMode(ThemeMode.DARK)
        message.success('已切换为深色模式')
    } else if (key === 'system') {
        themeStore.followSystemTheme()
        message.success('已设置为跟随系统主题')
    }
}

const handleColorChange = (color: string) => {
    themeStore.setColorPrimary(color)
    customColor.value = color
    message.success('主题色已更新')
    colorPickerVisible.value = false
}

const handleCustomColorChange = (value: string) => {
    if (value) {
        themeStore.setColorPrimary(value)
        message.success('主题色已更新')
        colorPickerVisible.value = false
    }
}
</script>

<style scoped>
.theme-switch-container {
    display: flex;
    align-items: center;
}

.theme-switch {
    transition: all 0.3s;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
}

.theme-switch:hover {
    background-color: rgba(255, 255, 255, 0.1);
    transform: scale(1.05);
}

:deep(.ant-dropdown-menu-item) {
    display: flex;
    align-items: center;
    gap: 8px;
}

.color-picker-menu {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 8px 0;
}

.preset-colors {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 8px;
}

.color-item {
    width: 24px;
    height: 24px;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
    border: 2px solid transparent;
    position: relative;
}

.color-item:hover {
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.color-item.active {
    border-color: rgba(0, 0, 0, 0.2);
    box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.1);
}

.dark-layout .color-item.active {
    border-color: rgba(255, 255, 255, 0.3);
    box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
}

.color-check {
    color: #fff;
    font-size: 12px;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.custom-color-picker {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.custom-color-label {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    margin-bottom: 4px;
}

body.dark-mode {
    .custom-color-label {
        color: rgba(255, 255, 255, 0.85);
    }
}

:deep(.ant-color-picker) {
    width: 100%;
}

:deep(.ant-color-picker-trigger) {
    width: 100%;
}

</style>

<style>
/* 浅色模式下的悬停效果 */
.dark-layout .theme-switch:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.layout-container:not(.dark-layout) .theme-switch:hover {
    background-color: rgba(0, 0, 0, 0.06);
}
</style>

