<template>
    <div class="data-list-container">
        <a-card :bordered="false" class="page-card">
            <template #title>
                <h2 class="page-title">{{ title }}</h2>
            </template>
            <template #extra>
                <div class="header-actions">
                    <!-- 顶部操作按钮 -->
                    <template v-if="extraActions && extraActions.length">
                        <a-button
                            v-for="action in extraActions"
                            :key="action.key"
                            :type="action.type || 'primary'"
                            :danger="action.danger"
                            @click="action.onClick"
                            class="action-btn"
                        >
                            <template v-if="action.icon" #icon>
                                <component :is="action.icon" />
                            </template>
                            {{ action.label }}
                        </a-button>
                    </template>
                    <!-- 批量删除按钮 -->
                    <a-button
                        v-if="showBatchDelete && selectedRowKeys.length > 0"
                        :disabled="selectedRowKeys.length === 0"
                        @click="handleBatchDelete"
                        class="action-btn"
                    >
                        <template #icon>
                            <DeleteOutlined />
                        </template>
                        批量删除
                    </a-button>
                    <!-- 刷新按钮 -->
                    <a-button
                        type="default"
                        :loading="loading"
                        @click="fetchData"
                        class="action-btn"
                        title="刷新"
                    >
                        <template #icon>
                            <ReloadOutlined :spin="loading" />
                        </template>
                        刷新
                    </a-button>
                    <!-- 列设置按钮 -->
                    <a-button
                        type="default"
                        @click="showColumnSettings = !showColumnSettings"
                        class="action-btn"
                        title="列设置"
                    >
                        <template #icon>
                            <SettingOutlined />
                        </template>
                        列设置
                    </a-button>
                </div>
            </template>

            <!-- 搜索表单 -->
            <div v-if="searchFields && searchFields.length > 0" class="search-form-wrapper">
                <div class="search-form-header" @click="toggleSearchForm">
                    <span class="search-form-title">搜索条件</span>
                    <component :is="searchFormCollapsed ? DownOutlined : UpOutlined" class="collapse-icon" />
                </div>
                <a-collapse v-model:activeKey="searchFormCollapseKey" :bordered="false" class="search-form-collapse">
                    <a-collapse-panel key="search" :show-arrow="false">
                        <a-form :model="searchForm" layout="inline" @submit.prevent="handleSearch" class="search-form">
                    <a-form-item
                        v-for="field in searchFields"
                        :key="field.key"
                        :label="field.label"
                    >
                        <!-- 输入框（支持单个和多搜索词） -->
                        <a-input
                            v-if="field.type === 'input'"
                            v-model:value="searchForm[field.key]"
                            :placeholder="field.multiple ? `${field.placeholder || '请输入' + field.label}（支持多个，用逗号分隔）` : (field.placeholder || `请输入${field.label}`)"
                            allow-clear
                            :style="{ width: field.width || '180px' }"
                            @pressEnter="handleSearch"
                        />
                        <!-- 选择器 -->
                        <a-select
                            v-else-if="field.type === 'select'"
                            v-model:value="searchForm[field.key]"
                            :placeholder="field.placeholder || `请选择${field.label}`"
                            allow-clear
                            :style="{ width: field.width || '120px' }"
                        >
                            <a-select-option
                                v-for="option in field.options"
                                :key="option.value"
                                :value="option.value"
                            >
                                {{ option.label }}
                            </a-select-option>
                        </a-select>
                        <!-- 日期选择器 -->
                        <a-date-picker
                            v-else-if="field.type === 'date'"
                            v-model:value="searchForm[field.key]"
                            :placeholder="field.placeholder || `请选择${field.label}`"
                            :style="{ width: field.width || '180px' }"
                        />
                        <!-- 日期范围选择器 -->
                        <a-range-picker
                            v-else-if="field.type === 'dateRange'"
                            v-model:value="searchForm[field.key]"
                            :style="{ width: field.width || '240px' }"
                            :show-time="field.showTime !== false"
                            :format="field.format || 'YYYY-MM-DD HH:mm:ss'"
                            :placeholder="field.placeholder || ['开始时间', '结束时间']"
                        />
                    </a-form-item>
                    <a-form-item>
                        <a-button type="primary" @click="handleSearch">
                            <template #icon>
                                <SearchOutlined />
                            </template>
                            搜索
                        </a-button>
                        <a-button style="margin-left: 8px" @click="handleReset">
                            <template #icon>
                                <ReloadOutlined />
                            </template>
                            重置
                        </a-button>
                    </a-form-item>
                </a-form>
                    </a-collapse-panel>
                </a-collapse>
            </div>

            <!-- 操作栏（底部按钮，用于自定义内容） -->
            <div v-if="$slots.extra" class="extra-actions">
                <slot name="extra"></slot>
            </div>

            <!-- 列设置弹窗 -->
            <a-drawer
                v-model:open="showColumnSettings"
                title="列设置"
                placement="right"
                :width="360"
                :mask-closable="true"
            >
                <div class="column-settings">
                    <div class="column-settings-header">
                        <a-button type="link" size="small" @click="handleSelectAllColumns">
                            全选
                        </a-button>
                        <a-button type="link" size="small" @click="handleUnselectAllColumns">
                            取消全选
                        </a-button>
                        <a-button type="link" size="small" @click="handleResetColumns">
                            重置
                        </a-button>
                    </div>
                    <div class="column-settings-list">
                        <a-checkbox-group v-model:value="visibleColumnKeys" class="column-checkbox-group">
                            <div
                                v-for="column in availableColumns"
                                :key="column.key"
                                class="column-item"
                            >
                                <a-checkbox
                                    :value="column.key"
                                    :disabled="column.fixed || column.key === 'action'"
                                >
                                    {{ column.title }}
                                </a-checkbox>
                                <span v-if="column.fixed || column.key === 'action'" class="column-fixed-tag">
                                    (固定)
                                </span>
                            </div>
                        </a-checkbox-group>
                    </div>
                </div>
            </a-drawer>

            <!-- 数据表格 -->
            <a-table
                :columns="displayColumns"
                :data-source="tableData"
                :pagination="paginationConfig"
                :loading="loading"
                :scroll="scroll"
                :size="tableSize"
                :row-key="rowKey"
                :row-selection="showRowSelection ? rowSelection : undefined"
                class="data-table"
                @change="handleTableChange"
            >
                <!-- 自定义列插槽和操作列 -->
                <template #bodyCell="{ column, record, index }">
                    <!-- 自定义列插槽 -->
                    <template v-if="column.key && $slots[`column-${column.key}`]">
                        <slot :name="`column-${column.key}`" :column="column" :record="record" :index="index"></slot>
                    </template>
                    <!-- 操作列插槽 -->
                    <template v-else-if="column.key === 'action' && $slots.action">
                        <slot name="action" :record="record"></slot>
                    </template>
                    <!-- 默认操作列（如果配置了 actions） -->
                    <template v-else-if="column.key === 'action' && actions && actions.length">
                        <a-space :size="8">
                            <template v-for="action in actions" :key="action.key">
                                <a-popconfirm
                                    v-if="action.confirm"
                                    :title="action.confirm"
                                    :ok-text="action.okText || '确定'"
                                    :cancel-text="action.cancelText || '取消'"
                                    @confirm="() => action.onClick(record)"
                                >
                                    <a-button
                                        :type="getActionButtonType(action)"
                                        :size="action.size || 'small'"
                                        :danger="action.danger"
                                        :disabled="typeof action.disabled === 'function' ? action.disabled(record) : action.disabled"
                                        class="action-button"
                                    >
                                        {{ action.label }}
                                    </a-button>
                                </a-popconfirm>
                                <a-button
                                    v-else
                                    :type="getActionButtonType(action)"
                                    :size="action.size || 'small'"
                                    :danger="action.danger"
                                    :disabled="typeof action.disabled === 'function' ? action.disabled(record) : action.disabled"
                                    @click="() => action.onClick(record)"
                                    class="action-button"
                                >
                                    {{ action.label }}
                                </a-button>
                            </template>
                        </a-space>
                    </template>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { SearchOutlined, ReloadOutlined, UpOutlined, DownOutlined, DeleteOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import type { TableColumnsType, TableProps } from 'ant-design-vue'

/**
 * 搜索字段配置
 */
export interface SearchField {
    key: string
    label: string
    type: 'input' | 'select' | 'date' | 'dateRange'
    placeholder?: string
    width?: string
    options?: Array<{ label: string; value: string | number }>
    multiple?: boolean // 是否支持多个搜索词（用逗号或换行分隔）
    separator?: string // 分隔符，默认为逗号和换行
}

/**
 * 操作按钮配置
 */
export interface ActionButton {
    key: string
    label: string
    type?: 'default' | 'primary' | 'dashed' | 'link' | 'text'
    size?: 'small' | 'middle' | 'large'
    danger?: boolean
    confirm?: string
    okText?: string
    cancelText?: string
    onClick: (record: any) => void | Promise<void>
    disabled?: boolean | ((record: any) => boolean)
}

/**
 * 顶部操作按钮配置
 */
export interface ExtraAction {
    key: string
    label: string
    type?: 'default' | 'primary' | 'dashed' | 'link' | 'text'
    danger?: boolean
    icon?: any
    onClick: () => void | Promise<void>
}

/**
 * 分页响应数据
 */
export interface PageResponse<T> {
    code: number
    msg?: string
    message?: string
    data: T
    total?: number
    page?: number
    page_size?: number
    pageSize?: number
}

/**
 * Props
 */
interface Props {
    // 标题
    title?: string
    // 搜索字段配置
    searchFields?: SearchField[]
    // 表格列配置
    columns: TableColumnsType
    // 数据获取函数
    fetchData: (params: any) => Promise<PageResponse<any[]>>
    // 初始搜索参数
    initialSearchParams?: Record<string, any>
    // 操作按钮配置
    actions?: ActionButton[]
    // 是否显示操作列
    showActions?: boolean
    // 顶部额外操作按钮
    extraActions?: ExtraAction[]
    // 表格滚动配置
    scroll?: { x?: number | string; y?: number | string }
    // 表格尺寸
    tableSize?: 'default' | 'middle' | 'small'
    // 行键
    rowKey?: string | ((record: any) => string)
    // 默认每页数量
    defaultPageSize?: number
    // 每页数量选项
    pageSizeOptions?: string[]
    // 是否显示行选择（复选框）
    showRowSelection?: boolean
    // 是否显示批量删除
    showBatchDelete?: boolean
    // 批量删除回调
    onBatchDelete?: (selectedKeys: string[]) => Promise<void>
    // 列设置存储的key（用于localStorage）
    columnSettingsKey?: string
}

const props = withDefaults(defineProps<Props>(), {
    title: '数据列表',
    searchFields: () => [],
    showActions: true,
    tableSize: 'middle',
    rowKey: (record: any) => record.id || record.key || record.userId || record._id || Math.random(),
    defaultPageSize: 10,
    pageSizeOptions: () => ['10', '20', '50', '100'],
    scroll: () => ({ x: 1200 }),
    showRowSelection: false,
    showBatchDelete: false,
})

// 搜索表单
const searchForm = reactive<Record<string, any>>({})
const searchFormCollapsed = ref(false)
const searchFormCollapseKey = ref<string[]>(['search'])

// 切换搜索表单显示/隐藏
const toggleSearchForm = () => {
    searchFormCollapsed.value = !searchFormCollapsed.value
    searchFormCollapseKey.value = searchFormCollapsed.value ? [] : ['search']
}

// 表格数据
const tableData = ref<any[]>([])
const loading = ref(false)
const total = ref(0)

// 行选择
const selectedRowKeys = ref<string[]>([])
const showColumnSettings = ref(false)

// 列设置
const visibleColumnKeys = ref<string[]>([])
const originalColumns = computed(() => props.columns)

// 获取所有可用的列（排除固定的列）
const availableColumns = computed(() => {
    return originalColumns.value.filter((col) => {
        const key = col.key as string
        // 操作列和固定的列必须显示，但也要在列表中显示（禁用状态）
        return key
    })
})

// 初始化可见列
const initVisibleColumns = () => {
    const storageKey = props.columnSettingsKey || `column_settings_${props.title}`
    
    // 从localStorage读取保存的列设置
    try {
        const saved = localStorage.getItem(storageKey)
        if (saved) {
            const savedKeys = JSON.parse(saved)
            // 验证保存的key是否在现有列中存在
            const validKeys = savedKeys.filter((key: string) =>
                originalColumns.value.some((col) => col.key === key)
            )
            // 确保固定列和操作列始终包含在内
            const fixedKeys = originalColumns.value
                .filter((col) => col.fixed || col.key === 'action')
                .map((col) => col.key as string)
                .filter((key) => key)
            
            // 合并固定列和保存的列
            const allKeys = [...new Set([...fixedKeys, ...validKeys])]
            
            if (allKeys.length > fixedKeys.length) {
                visibleColumnKeys.value = allKeys
                return
            }
        }
    } catch (error) {
        console.error('读取列设置失败:', error)
    }
    
    // 如果没有保存的设置，默认显示所有列
    visibleColumnKeys.value = originalColumns.value
        .map((col) => col.key as string)
        .filter((key) => key)
}

// 保存列设置到localStorage
const saveColumnSettings = () => {
    const storageKey = props.columnSettingsKey || `column_settings_${props.title}`
    try {
        localStorage.setItem(storageKey, JSON.stringify(visibleColumnKeys.value))
    } catch (error) {
        console.error('保存列设置失败:', error)
    }
}

// 根据可见列过滤显示的列
const displayColumns = computed(() => {
    return originalColumns.value.filter((col) => {
        const key = col.key as string
        // 操作列和固定的列始终显示
        if (col.fixed || key === 'action') {
            return true
        }
        // 复选框列始终显示
        if (key === 'selection') {
            return true
        }
        // 其他列根据用户选择显示
        return visibleColumnKeys.value.includes(key)
    })
})

// 监听可见列变化，自动保存（延迟保存，避免频繁写入）
let saveTimer: ReturnType<typeof setTimeout> | null = null
watch(visibleColumnKeys, () => {
    if (saveTimer) {
        clearTimeout(saveTimer)
    }
    saveTimer = setTimeout(() => {
        saveColumnSettings()
    }, 300)
}, { deep: true })

// 全选列
const handleSelectAllColumns = () => {
    visibleColumnKeys.value = originalColumns.value
        .map((col) => col.key as string)
        .filter((key) => key)
}

// 取消全选列（保留固定的列）
const handleUnselectAllColumns = () => {
    visibleColumnKeys.value = originalColumns.value
        .filter((col) => col.fixed || col.key === 'action')
        .map((col) => col.key as string)
        .filter((key) => key)
}

// 重置列设置
const handleResetColumns = () => {
    const storageKey = props.columnSettingsKey || `column_settings_${props.title}`
    localStorage.removeItem(storageKey)
    initVisibleColumns()
    message.success('已重置为默认设置')
}

// 行选择配置
const rowSelection = computed(() => {
    if (!props.showRowSelection) return undefined
    
    return {
        selectedRowKeys: selectedRowKeys.value,
        onChange: (keys: string[]) => {
            selectedRowKeys.value = keys
        },
        onSelectAll: (selected: boolean, selectedRows: any[], changeRows: any[]) => {
            if (selected) {
                const keys = changeRows.map((row) => {
                    const key = typeof props.rowKey === 'function' ? props.rowKey(row) : row[props.rowKey]
                    return String(key)
                })
                selectedRowKeys.value = [...selectedRowKeys.value, ...keys]
            } else {
                const keys = changeRows.map((row) => {
                    const key = typeof props.rowKey === 'function' ? props.rowKey(row) : row[props.rowKey]
                    return String(key)
                })
                selectedRowKeys.value = selectedRowKeys.value.filter((key) => !keys.includes(key))
            }
        },
    }
})

// 获取操作按钮类型（根据按钮key自动设置）
const getActionButtonType = (action: ActionButton): 'default' | 'primary' | 'dashed' | 'link' | 'text' => {
    // 如果已经有type，使用原有type
    if (action.type) return action.type
    
    // 根据key自动设置类型
    if (action.key === 'edit') return 'primary'
    if (action.key === 'delete') return 'default'
    
    return 'default'
}

// 批量删除
const handleBatchDelete = async () => {
    if (selectedRowKeys.value.length === 0) return
    
    Modal.confirm({
        title: '确认批量删除',
        content: `确定要删除选中的 ${selectedRowKeys.value.length} 条记录吗？此操作不可恢复。`,
        okText: '确定',
        cancelText: '取消',
        okType: 'danger',
        onOk: async () => {
            if (props.onBatchDelete) {
                try {
                    await props.onBatchDelete(selectedRowKeys.value)
                    selectedRowKeys.value = []
                    fetchData()
                } catch (error) {
                    console.error('批量删除失败:', error)
                }
            }
        },
    })
}

// 初始化搜索表单
const initSearchForm = () => {
    if (props.searchFields && props.searchFields.length > 0) {
        props.searchFields.forEach((field) => {
            if (props.initialSearchParams && props.initialSearchParams[field.key] !== undefined) {
                searchForm[field.key] = props.initialSearchParams[field.key]
            } else {
                searchForm[field.key] = field.type === 'select' ? undefined : ''
            }
        })
    }
}

// 分页配置
const paginationConfig = reactive({
    current: 1,
    pageSize: props.defaultPageSize,
    total: 0,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number, range: [number, number]) => {
        if (total === 0) return '暂无数据'
        return `共 ${total} 条，第 ${range[0]}-${range[1]} 条`
    },
    pageSizeOptions: props.pageSizeOptions,
    onChange: (page: number, pageSize: number) => {
        paginationConfig.current = page
        paginationConfig.pageSize = pageSize
        fetchData()
    },
    onShowSizeChange: (current: number, size: number) => {
        paginationConfig.current = 1
        paginationConfig.pageSize = size
        fetchData()
    },
})

// 处理多个搜索词
const parseMultipleKeywords = (value: string, separator: string = ','): string[] => {
    if (!value || typeof value !== 'string') {
        return []
    }
    // 支持逗号分隔，去除空字符串和首尾空格
    return value
        .split(separator)
        .map((item) => item.trim())
        .filter((item) => item !== '')
}

// 构建请求参数
const buildParams = () => {
    const params: Record<string, any> = {
        page: paginationConfig.current,
        page_size: paginationConfig.pageSize,
    }

    // 添加搜索参数
    Object.keys(searchForm).forEach((key) => {
        const value = searchForm[key]
        if (value !== '' && value !== undefined && value !== null) {
            // 查找对应的字段配置
            const field = props.searchFields?.find((f) => f.key === key)
            
            // 如果是支持多个搜索词的字段，转换为数组
            if (field?.multiple && typeof value === 'string') {
                const keywords = parseMultipleKeywords(value, field.separator)
                if (keywords.length > 0) {
                    // 对于多个搜索词，始终使用数组格式
                    // axios会自动将数组转换为 key[]=value1&key[]=value2 格式
                    // 即使只有一个值，也使用数组格式以保持一致性
                    params[key] = keywords
                }
            } else {
                params[key] = value
            }
        }
    })

    return params
}

// 获取数据
const fetchData = async () => {
    loading.value = true
    try {
        const params = buildParams()
        console.log('请求参数:', params)
        const response = await props.fetchData(params)
        console.log('响应数据:', response)
        
        // 后端返回 code=200 表示成功，code=0 也表示成功
        if (response.code === 0 || response.code === 200) {
            // 确保 data 是数组
            const data = Array.isArray(response.data) ? response.data : []
            console.log('处理后的数据:', data)
            tableData.value = data
            paginationConfig.total = response.total || 0
            total.value = response.total || 0
            
            // 同步后端返回的分页信息
            if (response.page) {
                paginationConfig.current = response.page
            }
            if (response.page_size || response.pageSize) {
                paginationConfig.pageSize = response.page_size || response.pageSize || props.defaultPageSize
            }
        } else if (response.code === 403) {
            // 权限不足，显示友好提示
            console.error('权限不足:', response.msg || response.message)
            tableData.value = []
            // 错误消息已在 request.ts 中显示，这里不需要重复显示
        } else {
            console.error('获取数据失败:', response.msg || response.message || '未知错误', response)
            tableData.value = []
        }
    } catch (error: any) {
        console.error('获取数据失败:', error)
        console.error('错误详情:', error.response || error.message)
        tableData.value = []
        // 不抛出错误，避免阻断UI
    } finally {
        loading.value = false
    }
}

// 搜索
const handleSearch = () => {
    paginationConfig.current = 1
    fetchData()
}

// 重置
const handleReset = () => {
    initSearchForm()
    paginationConfig.current = 1
    fetchData()
}

// 表格变化处理（用于排序和筛选，分页由分页组件自己处理）
const handleTableChange = (pag: any, filters: any, sorter: any) => {
    // 分页已由 paginationConfig 的 onChange 处理，这里只处理筛选和排序
    if (filters || sorter) {
        fetchData()
    }
}

// 刷新
const refresh = () => {
    selectedRowKeys.value = []
    fetchData()
}

// 暴露方法给父组件
defineExpose({
    refresh,
    fetchData,
    reset: handleReset,
    search: handleSearch,
    selectedRowKeys,
})

// 初始化
onMounted(() => {
    initSearchForm()
    initVisibleColumns()
    fetchData()
})

// 监听columns变化，重新初始化可见列
watch(() => props.columns, () => {
    initVisibleColumns()
}, { deep: true })
</script>

<style scoped lang="less">
.data-list-container {
    padding: 0;
    background: transparent;
    min-height: 100%;

    .page-card {
        border-radius: 12px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        overflow: hidden;
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
            padding: 16px 20px;
        }

        .page-title {
            margin: 0;
            font-size: 20px;
            font-weight: 600;
            color: #262626;
            letter-spacing: 0.5px;
        }
        
        .header-actions {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .action-btn {
            display: flex;
            align-items: center;
            gap: 4px;
            border-radius: 6px;
            transition: all 0.3s;
        }
        
        .action-btn:hover {
            transform: translateY(-1px);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
    }

    .search-form-wrapper {
        margin-bottom: 12px;
    }

    .search-form-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px 12px;
        background: linear-gradient(135deg, #fafafa 0%, #f5f5f5 100%);
        border-radius: 6px 6px 0 0;
        border: 1px solid #f0f0f0;
        border-bottom: none;
        cursor: pointer;
        user-select: none;
        transition: all 0.2s;

        &:hover {
            background: linear-gradient(135deg, #f5f5f5 0%, #f0f0f0 100%);
        }

        .search-form-title {
            font-size: 14px;
            font-weight: 500;
            color: #595959;
        }

        .collapse-icon {
            color: #8c8c8c;
            font-size: 12px;
            transition: transform 0.2s;
        }
    }

    .search-form-collapse {
        :deep(.ant-collapse-item) {
            border: 1px solid #f0f0f0;
            border-top: none;
            border-radius: 0 0 6px 6px;
            background: linear-gradient(135deg, #fafafa 0%, #f5f5f5 100%);
        }

        :deep(.ant-collapse-header) {
            display: none;
        }

        :deep(.ant-collapse-content) {
            border: none;
            background: transparent;
        }

        :deep(.ant-collapse-content-box) {
            padding: 12px;
        }
    }

    .search-form {
        margin: 0;
        padding: 0;
        background: transparent;

        :deep(.ant-form-item) {
            margin-bottom: 8px;
        }

        :deep(.ant-form-item-label) {
            padding-bottom: 4px;
        }

        :deep(.ant-form-item-label > label) {
            font-weight: 500;
            color: #595959;
            font-size: 13px;
        }

        :deep(.ant-form-item-control) {
            min-height: auto;
        }

        :deep(.ant-input),
        :deep(.ant-select-selector) {
            border-radius: 6px;
            transition: all 0.2s;

            &:hover {
                border-color: var(--color-primary-hover);
            }

            &:focus {
                border-color: var(--color-primary);
                box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb, 24, 144, 255), 0.1);
            }
        }

        :deep(.ant-btn-primary) {
            border-radius: 6px;
            background: var(--color-primary);
            border-color: var(--color-primary);
            box-shadow: 0 2px 4px rgba(var(--color-primary-rgb, 24, 144, 255), 0.2);
            transition: all 0.2s;

            &:hover {
                background: var(--color-primary-hover);
                border-color: var(--color-primary-hover);
                transform: translateY(-1px);
                box-shadow: 0 4px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.3);
            }

            &:active {
                background: var(--color-primary-active);
                border-color: var(--color-primary-active);
            }
        }
    }

    .extra-actions {
        margin-bottom: 12px;
        display: flex;
        gap: 8px;
        align-items: center;
    }

    .data-table {
        :deep(.ant-table) {
            .ant-table-container {
                border-radius: 4px;
                overflow: hidden;
            }

            .ant-table-thead > tr > th {
                background: linear-gradient(180deg, #fafafa 0%, #f5f5f5 100%);
                font-weight: 600;
                border-bottom: 2px solid #e8e8e8;
                padding: 10px 12px;
                text-align: center;
            }

            .ant-table-tbody > tr {
                transition: all 0.2s;

                &:hover > td {
                    background: #f5f7fa !important;
                    color: #262626 !important;
                }

                > td {
                    padding: 10px 12px;
                    border-bottom: 1px solid #f0f0f0;
                    color: #262626;
                    text-align: center;
                }
            }

            .ant-table-tbody > tr:last-child > td {
                border-bottom: none;
            }
        }

        :deep(.ant-pagination) {
            margin-top: 16px;
            padding-top: 12px;
            border-top: 1px solid #f0f0f0;
            text-align: right;
        }

        // 美化操作按钮
        .action-button {
            border-radius: 4px;
            transition: all 0.2s;
            font-size: 13px;
            padding: 4px 12px;
            height: auto;
            line-height: 1.5;
            
            &:hover {
                transform: translateY(-1px);
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            }
        }
        
        :deep(.ant-btn-link) {
            padding: 0 8px;
            height: auto;
            font-size: 14px;
            color: var(--color-primary);

            &:hover {
                color: var(--color-primary-hover);
                transform: translateY(-1px);
            }

            &:active {
                color: var(--color-primary-active);
            }

            &.ant-btn-dangerous {
                color: var(--ant-error-color, #ff4d4f);

                &:hover {
                    color: var(--ant-error-color-hover, #ff7875);
                }
            }
        }
        
        :deep(.ant-table-selection-column) {
            text-align: center;
        }
    }
    
    .column-settings {
        .column-settings-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-bottom: 16px;
            border-bottom: 1px solid #f0f0f0;
            margin-bottom: 16px;
        }
        
        .column-settings-list {
            .column-checkbox-group {
                width: 100%;
                display: flex;
                flex-direction: column;
                gap: 12px;
            }
            
            .column-item {
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 8px 12px;
                border-radius: 4px;
                transition: background-color 0.2s;
                
                &:hover {
                    background-color: #f5f5f5;
                }
                
                .column-fixed-tag {
                    font-size: 12px;
                    color: #8c8c8c;
                    margin-left: 8px;
                }
            }
        }
    }
}

// 暗色模式
body.dark-mode {
    .data-list-container {
        background: transparent;

        .page-card {
            background: #1f1f1f;
            border-color: #303030;

            :deep(.ant-card-head) {
                border-bottom-color: #303030;
                background: linear-gradient(135deg, #1f1f1f 0%, #262626 100%);
            }

            .page-title {
                color: #fff;
            }
        }
        
        .column-settings {
            .column-settings-header {
                border-bottom-color: #434343;
            }
            
            .column-item {
                &:hover {
                    background-color: #2a2a2a;
                }
                
                .column-fixed-tag {
                    color: rgba(255, 255, 255, 0.45);
                }
            }
        }

        .search-form-wrapper {
            .search-form-header {
                background: linear-gradient(135deg, #262626 0%, #1f1f1f 100%);
                border-color: #434343;

                .search-form-title {
                    color: rgba(255, 255, 255, 0.85);
                }

                .collapse-icon {
                    color: rgba(255, 255, 255, 0.65);
                }

                &:hover {
                    background: linear-gradient(135deg, #303030 0%, #262626 100%);
                }
            }

            .search-form-collapse {
                :deep(.ant-collapse-item) {
                    border-color: #434343;
                    background: linear-gradient(135deg, #262626 0%, #1f1f1f 100%);
                }

                :deep(.ant-input),
                :deep(.ant-select-selector) {
                    background: #1f1f1f;
                    border-color: #434343;
                    color: rgba(255, 255, 255, 0.85);

                    &::placeholder {
                        color: rgba(255, 255, 255, 0.3);
                    }

                    &:hover {
                        border-color: var(--color-primary-hover);
                    }

                    &:focus {
                        border-color: var(--color-primary);
                    }
                }

                :deep(.ant-form-item-label > label) {
                    color: rgba(255, 255, 255, 0.65);
                }
            }
        }

        .search-form {
            background: transparent;
        }

        .data-table {
            :deep(.ant-table-thead > tr > th) {
                background: linear-gradient(180deg, #262626 0%, #1f1f1f 100%);
                border-color: #434343;
                border-bottom: 2px solid #434343;
            }

            :deep(.ant-table-tbody > tr) {
                > td {
                    color: rgba(255, 255, 255, 0.85);
                    background: #1f1f1f;
                }

                &:hover > td {
                    background: #2a2a2a !important;
                    color: rgba(255, 255, 255, 0.95) !important;
                }

                &:nth-child(even) > td {
                    background: #1f1f1f;
                }

                &:nth-child(even):hover > td {
                    background: #2a2a2a !important;
                    color: rgba(255, 255, 255, 0.95) !important;
                }
            }

            :deep(.ant-table-tbody > tr > td) {
                border-color: #434343;
            }

            :deep(.ant-pagination) {
                border-top-color: #434343;

                .ant-pagination-total-text,
                .ant-pagination-item,
                .ant-pagination-prev,
                .ant-pagination-next {
                    color: rgba(255, 255, 255, 0.85);
                }

                .ant-pagination-item {
                    background: #1f1f1f;
                    border-color: #434343;

                    a {
                        color: rgba(255, 255, 255, 0.85);
                    }

                    &:hover {
                        background: #262626;
                        border-color: var(--color-primary);

                        a {
                            color: var(--color-primary);
                        }
                    }

                    &.ant-pagination-item-active {
                        background: var(--color-primary);
                        border-color: var(--color-primary);

                        a {
                            color: #fff;
                        }
                    }
                }

                .ant-select-selector {
                    background: #1f1f1f;
                    border-color: #434343;
                    color: rgba(255, 255, 255, 0.85);
                }
            }

            // 暗色模式下的操作按钮
            :deep(.ant-btn-link) {
                color: var(--color-primary);

                &:hover {
                    color: var(--color-primary-hover);
                }

                &:active {
                    color: var(--color-primary-active);
                }
            }
        }
    }
}
</style>

