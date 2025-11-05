<template>
    <div class="file-preview">
        <a-image
            v-if="isImage"
            :src="file.file_url"
            :preview="false"
            :width="60"
            :height="60"
            style="object-fit: cover; border-radius: 4px; cursor: pointer;"
            @click="handlePreview"
        />
        <div
            v-else
            class="file-icon"
            @click="handlePreview"
        >
            <component :is="getFileIcon" :style="{ fontSize: '32px', color: getFileIconColor }" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import {
    FileImageOutlined,
    FileTextOutlined,
    FilePdfOutlined,
    FileExcelOutlined,
    FileWordOutlined,
    FilePptOutlined,
    FileZipOutlined,
    PlayCircleOutlined,
    SoundOutlined,
    FileOutlined,
} from '@ant-design/icons-vue'
import type { FileResponse } from '@/api/file'

const props = defineProps<{
    file: FileResponse
}>()

const emit = defineEmits<{
    preview: [file: FileResponse]
}>()

const isImage = computed(() => props.file.category === 'image')

const getFileIcon = computed(() => {
    const ext = props.file.extension.toLowerCase()
    const category = props.file.category

    if (category === 'document') {
        if (ext === 'pdf') return FilePdfOutlined
        if (['doc', 'docx'].includes(ext)) return FileWordOutlined
        if (['xls', 'xlsx'].includes(ext)) return FileExcelOutlined
        if (['ppt', 'pptx'].includes(ext)) return FilePptOutlined
        if (['txt', 'md'].includes(ext)) return FileTextOutlined
        return FileTextOutlined
    }
    if (category === 'video') return PlayCircleOutlined
    if (category === 'audio') return SoundOutlined
    if (category === 'archive') return FileZipOutlined
    return FileOutlined
})

const getFileIconColor = computed(() => {
    const category = props.file.category
    const colorMap: Record<string, string> = {
        document: '#1890ff',
        video: '#ff4d4f',
        audio: '#52c41a',
        archive: '#faad14',
        other: '#8c8c8c',
    }
    return colorMap[category] || '#8c8c8c'
})

const handlePreview = () => {
    emit('preview', props.file)
}
</script>

<style scoped>
.file-preview {
    display: flex;
    align-items: center;
    justify-content: center;
}

.file-icon {
    width: 60px;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    cursor: pointer;
    background: #fafafa;
    transition: all 0.3s;
}

.file-icon:hover {
    border-color: #1890ff;
    background: #e6f7ff;
}
</style>

