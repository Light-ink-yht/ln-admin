<template>
    <div class="header-breadcrumb">
        <a-breadcrumb>
            <a-breadcrumb-item v-for="(item, index) in breadcrumbItems" :key="item.path">
                <span v-if="index === breadcrumbItems.length - 1">
                    {{ item.breadcrumbName }}
                </span>
                <a v-else @click.prevent="handleBreadcrumbClick(item.path)">
                    {{ item.breadcrumbName }}
                </a>
            </a-breadcrumb-item>
        </a-breadcrumb>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// 路由名称映射（可根据实际路由配置修改）
const routeNameMap: Record<string, string> = {
    '/': '首页',
    '/dashboard': '仪表盘',
    '/dashboard/workbench': '工作台',
    '/dashboard/analysis': '分析页',
    '/user': '用户管理',
    '/user/list': '用户列表',
    '/user/profile': '个人资料',
    '/role': '角色管理',
    '/role/list': '角色列表',
    '/permission': '权限管理',
    '/permission/list': '权限列表',
    '/settings': '设置',
    '/settings/basic': '基础设置',
    '/settings/security': '安全设置',
    '/system': '系统管理',
    '/system/config': '系统配置',
    '/system/log': '操作日志',
    '/sms': '短信管理',
    '/sms/template': '短信模板',
    '/sms/code': '短信验证码',
    '/menu': '菜单管理',
    '/menu/list': '菜单列表',
    '/ops': '系统运维',
    '/ops/api-doc': '接口文档',
    '/ops/monitor': '系统监控',
}

// 将路径段转换为可读名称（将 kebab-case 或 camelCase 转换为中文）
const formatPathName = (segment: string): string => {
    // 如果已经有映射，使用映射
    return segment
        .replace(/-/g, ' ')
        .replace(/([A-Z])/g, ' $1')
        .trim()
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ')
}

// 根据当前路由生成面包屑
const breadcrumbItems = computed(() => {
    const items: Array<{ path: string; breadcrumbName: string }> = []
    
    // 如果路径是根路径，直接返回首页
    if (route.path === '/' || route.path === '') {
        items.push({
            path: '/',
            breadcrumbName: routeNameMap['/'] || '首页',
        })
        return items
    }
    
    // 使用路由的matched数组来生成面包屑（更准确）
    const matched = route.matched.filter((item) => item.meta && item.meta.title)
    
    // 如果没有匹配的路由，使用路径解析
    if (matched.length === 0) {
        const paths = route.path.split('/').filter(Boolean)
        let currentPath = ''
        
        paths.forEach((segment) => {
            currentPath += `/${segment}`
            const mappedName = routeNameMap[currentPath]
            const name = mappedName || formatPathName(segment)
            
            items.push({
                path: currentPath,
                breadcrumbName: name,
            })
        })
    } else {
        // 使用matched生成面包屑
        matched.forEach((match) => {
            const path = match.path
            const title = match.meta?.title as string
            const name = title || routeNameMap[path] || formatPathName(path.split('/').pop() || '')
            
            items.push({
                path: path === '/' ? '/' : path,
                breadcrumbName: name,
            })
        })
    }
    
    // 确保至少有一个首页
    if (items.length === 0 || items[0].path !== '/') {
        items.unshift({
            path: '/',
            breadcrumbName: routeNameMap['/'] || '首页',
        })
    }
    
    return items
})

const handleBreadcrumbClick = (path: string) => {
    if (path !== route.path) {
        router.push(path)
    }
}
</script>

<style scoped>
.header-breadcrumb {
    flex: 1;
    display: flex;
    align-items: center;
    min-width: 0;
    overflow: hidden;
    padding-left: 24px;
    box-sizing: border-box;
}

:deep(.ant-breadcrumb) {
    font-size: 14px;
    line-height: 1.5;
}

:deep(.ant-breadcrumb-link) {
    color: rgba(0, 0, 0, 0.65);
    transition: color 0.3s, transform 0.2s;
    display: inline-block;
}

:deep(.ant-breadcrumb-link:hover) {
    color: rgba(0, 0, 0, 0.85);
    transform: translateY(-1px);
}

:deep(.ant-breadcrumb-separator) {
    color: rgba(0, 0, 0, 0.35);
    margin: 0 8px;
    font-size: 12px;
}

:deep(.ant-breadcrumb-item:last-child span) {
    color: rgba(0, 0, 0, 0.85);
    font-weight: 500;
}

/* 暗色模式 */
.dark-layout :deep(.ant-breadcrumb-link) {
    color: rgba(255, 255, 255, 0.65);
}

.dark-layout :deep(.ant-breadcrumb-link:hover) {
    color: rgba(255, 255, 255, 0.85);
}

.dark-layout :deep(.ant-breadcrumb-separator) {
    color: rgba(255, 255, 255, 0.35);
}

.dark-layout :deep(.ant-breadcrumb-item:last-child span) {
    color: rgba(255, 255, 255, 0.85);
}

/* 移动端优化 */
@media (max-width: 768px) {
    .header-breadcrumb {
        display: none;
    }
}
</style>

