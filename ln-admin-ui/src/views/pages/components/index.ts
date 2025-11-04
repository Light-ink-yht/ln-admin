export { default as DataList } from './DataList.vue'

export type {
    SearchField,
    ActionButton,
    ExtraAction,
    PageResponse,
} from './DataList.vue'

// 用户相关组件
export {
    UserAvatar,
    UserStatus,
    UserGender,
    UserText,
    PermissionDeniedAlert,
    getGenderText,
    mapUserData,
    mapUserList,
    getUserSearchFields,
    getUserColumns,
    USER_STATUS_OPTIONS,
    USER_GENDER_OPTIONS,
    useUserList,
} from './user'

// 角色相关工具
export {
    getRoleSearchFields,
    getRoleColumns,
    ROLE_STATUS_OPTIONS,
    mapRoleData,
    mapRoleList,
} from './role'

// 权限相关工具
export {
    getPermissionSearchFields,
    getPermissionColumns,
    PERMISSION_STATUS_OPTIONS,
    METHOD_OPTIONS,
    mapPermissionData,
    mapPermissionList,
} from './permission'

// 系统日志相关工具
export {
    getSystemLogSearchFields,
    getSystemLogColumns,
} from './system/systemLogConfig'

