import type { SearchField } from '../DataList.vue'
import type { TableColumnsType } from 'ant-design-vue'

/**
 * 获取菜单搜索字段配置
 */
export function getMenuSearchFields(): SearchField[] {
    return [
        {
            key: 'menu_type',
            label: '菜单类型',
            type: 'select',
            placeholder: '请选择菜单类型',
            width: '150px',
            options: [
                { label: '侧边栏菜单', value: '1' },
                { label: '用户菜单', value: '2' },
            ],
        },
        {
            key: 'status',
            label: '状态',
            type: 'select',
            placeholder: '请选择状态',
            width: '120px',
            options: [
                { label: '启用', value: '1' },
                { label: '禁用', value: '2' },
            ],
        },
    ]
}

/**
 * 获取菜单表格列配置
 */
export function getMenuColumns(): TableColumnsType {
    return [
        {
            title: '菜单标识',
            key: 'menu_key',
            dataIndex: 'menu_key',
            width: 150,
            align: 'center',
        },
        {
            title: '菜单标题',
            key: 'title',
            dataIndex: 'title',
            width: 150,
            align: 'center',
        },
        {
            title: '路由路径',
            key: 'path',
            dataIndex: 'path',
            width: 200,
            align: 'center',
        },
        {
            title: '图标',
            key: 'icon',
            dataIndex: 'icon',
            width: 120,
            align: 'center',
        },
        {
            title: '菜单类型',
            key: 'menu_type',
            dataIndex: 'menu_type',
            width: 120,
            align: 'center',
        },
        {
            title: '排序',
            key: 'sort',
            dataIndex: 'sort',
            width: 80,
            align: 'center',
        },
        {
            title: '状态',
            key: 'status',
            dataIndex: 'status',
            width: 100,
            align: 'center',
        },
        {
            title: '创建时间',
            key: 'created_at',
            dataIndex: 'created_at',
            width: 180,
            align: 'center',
        },
        {
            title: '操作',
            key: 'action',
            width: 180,
            align: 'center',
            fixed: 'right',
        },
    ]
}

