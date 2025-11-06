# 权限层级授权优化方案

## 需求概述

实现权限的层级授权机制：
- **超级管理员**：拥有所有权限，可以给管理员和普通用户授权
- **管理员**：只能看到超级管理员授予的权限，并且可以将这些权限赋予给普通用户
- **普通用户**：只能看到被授予的权限

## 实现方案

### 1. 数据库设计

创建角色权限授予关系表 `AA13`：
- `AAM001`：授予者角色ID（grantor_role_id）
- `AAM002`：被授予者角色ID（grantee_role_id）
- `AAM003`：权限ID（permission_id）

### 2. 核心逻辑

#### 2.1 权限查询逻辑
- 超级管理员：可以看到所有权限
- 其他角色：只能看到被授予的权限（通过 `AA13` 表查询）

#### 2.2 授权逻辑
- 超级管理员：可以授予任何权限给任何角色
- 其他角色：只能授予自己拥有的权限给其他角色
- 授权时需要：
  1. 检查当前用户是否有权限授予这些权限
  2. 记录授予关系到 `AA13` 表
  3. 更新 Casbin 策略

#### 2.3 角色权限分配逻辑
- 修改 `AssignPermissionToRole` 方法，支持层级授权
- 检查授予者是否有权限授予该权限
- 记录授予关系

### 3. 实现步骤

1. ✅ 创建角色权限授予关系表实体（`role_permission_grant.go`）
2. ✅ 创建 Repository 接口和实现
3. ✅ 修改 PermissionService，添加层级授权支持
4. ✅ 修改权限查询逻辑
5. ✅ 修改授权逻辑
6. ✅ 修改角色权限分配逻辑
7. ✅ 更新初始化数据，为超级管理员初始化授予关系
8. ✅ 创建数据库迁移脚本（`migrate_v5.sql`）

### 4. 关键方法

#### 4.1 获取用户可授予的权限列表
```go
GetGrantablePermissions(ctx context.Context, userID string) ([]*entity.Permission, error)
```

#### 4.2 检查用户是否有权限授予某个权限
```go
CanGrantPermission(ctx context.Context, grantorUserID, permissionID string) (bool, error)
```

#### 4.3 层级授权（为角色分配权限）
```go
GrantPermissionToRole(ctx context.Context, grantorUserID, roleID, permissionID string) error
```

#### 4.4 获取用户可见的权限列表
```go
GetVisiblePermissions(ctx context.Context, userID string) ([]*entity.Permission, error)
```

## 注意事项

1. 超级管理员角色标识：`super_admin`
2. 需要保持与现有 Casbin 策略的兼容性
3. 授权关系需要与 Casbin 策略同步
4. 删除权限或角色时，需要清理相关的授予关系

