# Casbin 权限控制指南

本文档详细介绍如何使用 Casbin 进行权限和资源控制。

## 📖 什么是 Casbin

Casbin 是一个强大和高效的开源访问控制库，支持 ACL、RBAC、ABAC 等多种访问控制模型。本项目使用 RBAC（基于角色的访问控制）模型。

### RBAC 模型

RBAC（Role-Based Access Control）基于角色的访问控制：

- **用户（User）**：系统的使用人员
- **角色（Role）**：权限的集合，如“管理员”、“普通用户”
- **权限（Permission）**：对资源的操作权限，如“查看用户列表”
- **资源（Resource）**：系统中的资源，如 API 路径

### 权限模型关系

```
用户 → 拥有 → 角色 → 拥有 → 权限 → 控制 → 资源
```

例如：
- 用户张三 → 拥有 → 管理员角色 → 拥有 → 查看用户列表权限 → 控制 → `/api/user/list` 资源

## 🏗️ 架构设计

### Casbin 策略存储

本项目使用 **GORM 适配器**将 Casbin 策略存储在数据库中：

- **策略表（casbin_rule）**：自动创建，存储所有权限规则
- **策略类型**：
  - `p`（Policy）：权限策略，格式：`角色, 资源路径, 请求方法`
  - `g`（Grouping）：角色分配，格式：`用户ID, 角色Key`

### RBAC 模型配置

模型文件：`configs/rbac_model.conf`

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

**配置说明：**
- `r = sub, obj, act`：请求定义为（用户、资源、操作）
- `p = sub, obj, act`：策略定义为（角色、资源、操作）
- `g = _, _`：角色定义为（用户，角色）的继承关系
- `e = ...`：策略效果，允许通过或拒绝访问
- `m = ...`：匹配规则，检查用户是否拥有角色，角色是否有权限

## 🔧 初始化

### 自动初始化

服务启动时会自动初始化 Casbin：

```go
// cmd/server/bootstrap.go
func (b *Bootstrap) InitCasbin() error {
    return casbin.Init()
}
```

### 初始化流程

1. 创建 GORM 适配器（连接数据库）
2. 加载 RBAC 模型配置
3. 从数据库加载策略规则
4. 启用自动保存（策略变更自动持久化）

### 策略表结构

Casbin 会自动创建 `casbin_rule` 表：

```sql
CREATE TABLE casbin_rule (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ptype VARCHAR(100) NOT NULL,      -- 策略类型：p 或 g
    v0 VARCHAR(100),                  -- 第1个参数
    v1 VARCHAR(100),                  -- 第2个参数
    v2 VARCHAR(100),                  -- 第3个参数
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

**策略存储格式：**

| ptype | v0 | v1 | v2 | 说明 |
|-------|----|----|----|------|
| g | user_id_123 | super_admin | | 用户拥有角色 |
| p | super_admin | /api/user/list | GET | 角色拥有权限 |

## 📝 使用方式

### 1. 角色管理

#### 1.1 创建角色

```go
// internal/application/service/permission_service.go
func (s *PermissionService) CreateRole(ctx context.Context, role *entity.Role) error {
    // 创建角色记录
    if err := s.roleRepo.Create(ctx, role); err != nil {
        return err
    }
    // Casbin 策略会自动通过角色表管理，不需要手动添加 g 规则
    return nil
}
```

#### 1.2 为用户分配角色

```go
// 通过权限服务分配角色（推荐）
func (s *PermissionService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
    // 1. 查询角色信息
    role, err := s.roleRepo.FindByID(ctx, roleID)
    
    // 2. 在数据库中添加用户角色关联
    err := s.userRoleRepo.AssignRole(ctx, userID, roleID)
    
    // 3. 更新 Casbin 策略：添加 g(user_id, role_key)
    err = s.casbinSvc.AddRoleForUser(ctx, userID, role.RoleKey)
    
    return nil
}
```

**Casbin 策略变更：**
```
g, user_id_123, super_admin  ← 用户 user_id_123 拥有 super_admin 角色
```

#### 1.3 移除用户角色

```go
func (s *PermissionService) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
    // 1. 查询角色
    role, err := s.roleRepo.FindByID(ctx, roleID)
    
    // 2. 删除数据库关联
    err := s.userRoleRepo.RemoveRole(ctx, userID, roleID)
    
    // 3. 更新 Casbin 策略：删除 g(user_id, role_key)
    err = s.casbinSvc.DeleteRoleForUser(ctx, userID, role.RoleKey)
    
    return nil
}
```

### 2. 权限管理

#### 2.1 创建权限

```go
func (s *PermissionService) CreatePermission(ctx context.Context, permission *entity.Permission) error {
    // 创建权限记录
    if err := s.permissionRepo.Create(ctx, permission); err != nil {
        return err
    }
    // 注意：创建权限后，需要手动为角色分配权限
    return nil
}
```

#### 2.2 为角色分配权限

```go
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
    // 1. 查询角色和权限
    role, err := s.roleRepo.FindByID(ctx, roleID)
    permission, err := s.permissionRepo.FindByID(ctx, permissionID)
    
    // 2. 添加 Casbin 策略：p(role_key, resource_path, method)
    err = s.casbinSvc.AddPolicy(ctx, role.RoleKey, permission.ResourcePath, permission.Method)
    
    return nil
}
```

**Casbin 策略变更：**
```
p, super_admin, /api/user/list, GET  ← super_admin 角色可以 GET /api/user/list
```

#### 2.3 移除角色权限

```go
func (s *PermissionService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
    role, err := s.roleRepo.FindByID(ctx, roleID)
    permission, err := s.permissionRepo.FindByID(ctx, permissionID)
    
    // 移除 Casbin 策略
    err = s.casbinSvc.RemovePolicy(ctx, role.RoleKey, permission.ResourcePath, permission.Method)
    
    return nil
}
```

### 3. 权限验证

#### 3.1 中间件验证（推荐）

使用 `CasbinMiddleware` 中间件自动验证权限：

```go
// internal/interfaces/http/middleware/casbin_middleware.go
func CasbinMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取用户ID（从JWT中间件设置）
        userID, exists := c.Get("user_id")
        
        // 2. 获取请求资源和操作方法
        obj := c.Request.URL.Path  // 资源路径：/api/user/list
        act := c.Request.Method     // 操作方法：GET
        
        // 3. 使用 Casbin 验证权限
        enforcer := casbin.GetEnforcer()
        ok, err := enforcer.Enforce(userID.(string), obj, act)
        
        // 4. 权限不足时拒绝访问
        if !ok {
            response.Forbidden(c, "权限不足")
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

#### 3.2 注册中间件

```go
// internal/interfaces/http/router/router.go
func SetupRouter(...) *gin.Engine {
    // 需要认证的路由组
    auth := api.Group("/user")
    auth.Use(middleware.Auth())              // JWT认证
    auth.Use(middleware.CasbinMiddleware())  // Casbin权限验证
    
    {
        auth.GET("/userinfo", userHandler.GetUserInfo)
        auth.GET("/list", userHandler.ListUsers)
    }
}
```

#### 3.3 手动验证（代码中验证）

```go
// 在业务代码中手动验证
func (h *UserHandler) DeleteUser(c *gin.Context) {
    userID := c.GetString("user_id")
    resource := "/api/user/delete"
    method := "DELETE"
    
    // 获取 Casbin Enforcer
    enforcer := casbin.GetEnforcer()
    
    // 验证权限
    allowed, err := enforcer.Enforce(userID, resource, method)
    if err != nil {
        response.InternalError(c, "权限检查失败")
        return
    }
    
    if !allowed {
        response.Forbidden(c, "您没有删除用户的权限")
        return
    }
    
    // 继续业务逻辑
    // ...
}
```

## 🎯 完整示例

### 示例：创建管理员角色并分配权限

```go
// 1. 创建角色
role := &entity.Role{
    RoleID:      "role_admin_001",
    RoleKey:     "admin",
    RoleName:    "管理员",
    Description: "系统管理员",
    Status:      "1",
    CreatorID:   "system",
}
permissionService.CreateRole(ctx, role)

// 2. 创建权限
permission := &entity.Permission{
    PermissionID:   "perm_user_list",
    PermissionKey:  "user:list",
    PermissionName: "查看用户列表",
    ResourcePath:   "/api/user/list",
    Method:         "GET",
    Status:         "1",
}
permissionService.CreatePermission(ctx, permission)

// 3. 为角色分配权限
permissionService.AssignPermissionToRole(ctx, role.RoleID, permission.PermissionID)

// 4. 为用户分配角色
permissionService.AssignRoleToUser(ctx, "user_id_123", role.RoleID)

// 结果：
// - 数据库：用户角色关联表中有记录
// - Casbin策略：
//   * g, user_id_123, admin
//   * p, admin, /api/user/list, GET
// - 用户 user_id_123 现在可以访问 GET /api/user/list
```

### 示例：权限检查流程

```
用户请求：GET /api/user/list
  ↓
1. JWT中间件验证Token，设置 user_id = "user_id_123"
  ↓
2. Casbin中间件拦截请求
  ↓
3. 提取请求信息：
   - sub (用户): "user_id_123"
   - obj (资源): "/api/user/list"
   - act (操作): "GET"
  ↓
4. Casbin查找策略：
   a) 查找 g 策略：g(user_id_123, ?)
      找到：g, user_id_123, admin
   b) 查找 p 策略：p(admin, /api/user/list, GET)
      找到：p, admin, /api/user/list, GET
  ↓
5. 匹配成功，允许访问
  ↓
6. 继续处理请求
```

## 🔍 策略查询

### 查看所有策略

```go
enforcer := casbin.GetEnforcer()

// 获取所有策略（p规则）
allPolicies := enforcer.GetPolicy()
for _, policy := range allPolicies {
    fmt.Printf("策略: %v\n", policy)  // [角色, 资源, 方法]
}

// 获取所有角色分配（g规则）
allGroupings := enforcer.GetGroupingPolicy()
for _, grouping := range allGroupings {
    fmt.Printf("角色分配: %v\n", grouping)  // [用户ID, 角色]
}
```

### 查询用户拥有的角色

```go
userID := "user_id_123"
enforcer := casbin.GetEnforcer()

// 获取用户的所有角色
roles, err := enforcer.GetRolesForUser(userID)
// 返回: ["super_admin", "admin"]
```

### 查询角色拥有的权限

```go
roleKey := "admin"
enforcer := casbin.GetEnforcer()

// 获取角色的所有权限策略
policies := enforcer.GetPermissionsForUser(roleKey)
for _, policy := range policies {
    fmt.Printf("权限: %s %s\n", policy[1], policy[2])  // 资源路径 和 方法
}
```

### 查询资源被哪些角色访问

```go
resource := "/api/user/list"
method := "GET"
enforcer := casbin.GetEnforcer()

// 查询哪些角色可以访问该资源
policies := enforcer.GetFilteredPolicy(1, resource, method)
for _, policy := range policies {
    fmt.Printf("角色 %s 可以 %s %s\n", policy[0], policy[2], policy[1])
}
```

## 🛠️ 数据库直接操作

### 查看策略表数据

```sql
-- 查看所有策略
SELECT * FROM casbin_rule;

-- 查看角色分配（g规则）
SELECT v0 as user_id, v1 as role_key 
FROM casbin_rule 
WHERE ptype = 'g';

-- 查看权限策略（p规则）
SELECT v0 as role_key, v1 as resource_path, v2 as method 
FROM casbin_rule 
WHERE ptype = 'p';
```

### 手动添加策略（不推荐）

```sql
-- 为用户分配角色
INSERT INTO casbin_rule (ptype, v0, v1) 
VALUES ('g', 'user_id_123', 'admin');

-- 为角色分配权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) 
VALUES ('p', 'admin', '/api/user/list', 'GET');
```

⚠️ **注意**：手动修改数据库后，需要重启服务或调用 `enforcer.LoadPolicy()` 重新加载策略。

## 📊 默认权限体系

### 默认角色

系统启动时自动创建：

1. **超级管理员（super_admin）**
   - 拥有所有权限
   - 自动分配所有现有和未来的权限

2. **管理员（admin）**
   - 基础管理权限（需要手动分配）

3. **普通用户（user）**
   - 基础查看权限（需要手动分配）

### 默认权限

系统启动时自动创建：

| 权限Key | 权限名称 | 资源路径 | 方法 |
|---------|---------|---------|------|
| user:list | 用户列表 | /api/user/list | GET |
| user:info | 用户信息 | /api/user/userinfo | GET |
| user:create | 创建用户 | /api/user/create | POST |
| user:update | 更新用户 | /api/user/update | PUT |
| user:delete | 删除用户 | /api/user/delete | DELETE |

### 超级管理员权限分配

启动时自动为超级管理员角色分配所有权限：

```go
// internal/infrastructure/defaultdata/init_data.go
func assignAllPermissionsToAdmin(...) {
    // 获取所有权限
    permissions, _ := permissionRepo.List(ctx, 1, 1000, nil)
    
    // 为超级管理员角色分配所有权限
    for _, permission := range permissions {
        casbinSvc.AddPolicy(ctx, "super_admin", permission.ResourcePath, permission.Method)
    }
}
```

## ⚠️ 注意事项

### 1. 策略同步

- **自动保存**：项目配置了 `enabler.EnableAutoSave(true)`，策略变更会自动保存到数据库
- **缓存刷新**：服务重启时会自动加载最新策略
- **多实例部署**：多实例环境下，策略变更可能延迟生效，建议使用分布式配置中心

### 2. 性能考虑

- Casbin 策略会加载到内存，访问速度快
- 大量策略时（>10万），考虑使用分布式缓存
- 定期清理无效策略，提升匹配效率

### 3. 权限粒度

- **粗粒度**：按模块控制（如：用户管理、订单管理）
- **细粒度**：按具体操作控制（如：查看列表、查看详情、创建、编辑、删除）

建议：
- API 级别使用粗粒度控制
- 数据级别使用细粒度控制（在业务代码中实现）

### 4. 权限继承

当前模型不支持角色继承，如果需要：

1. 修改 RBAC 模型文件，添加继承关系
2. 在策略中添加角色继承规则：`g, admin, super_admin`
3. 修改匹配器支持继承查询

## 🔧 故障排查

### Q1: 权限验证总是失败

**检查步骤：**
1. 确认用户已分配角色：查询 `casbin_rule` 表，`ptype='g'`
2. 确认角色有权限：查询 `casbin_rule` 表，`ptype='p'`
3. 检查资源路径是否匹配（注意斜杠、大小写）
4. 检查请求方法是否匹配（GET、POST、PUT、DELETE）

### Q2: 策略未生效

**可能原因：**
1. 策略已添加但未保存（检查数据库是否有记录）
2. 服务未重启，策略未重新加载
3. 缓存问题（重启服务清除内存缓存）

**解决方案：**
```go
// 手动重新加载策略
enforcer := casbin.GetEnforcer()
enforcer.LoadPolicy()
```

### Q3: 权限表查询慢

**优化方案：**
1. 为 `ptype`, `v0`, `v1`, `v2` 添加索引
2. 定期清理无效策略
3. 使用 Redis 缓存热门策略

## 📚 参考资源

- [Casbin 官方文档](https://casbin.org/)
- [RBAC 模型详解](https://casbin.org/docs/en/rbac)
- [GORM 适配器文档](https://github.com/casbin/gorm-adapter)

## 🔄 下一步

- 查看 [API 接口文档](./API.md)
- 查看 [权限管理API](./PERMISSION_API.md)
- 查看 [启动指南](./START.md)

