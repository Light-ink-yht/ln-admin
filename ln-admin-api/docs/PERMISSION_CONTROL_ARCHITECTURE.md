# LN Admin 权限控制架构图

## 系统权限控制架构

```mermaid
graph TB
    subgraph "用户层"
        U1[用户1<br/>User ID: user_001]
        U2[用户2<br/>User ID: user_002]
        U3[用户3<br/>User ID: user_003]
    end

    subgraph "角色层"
        SA[超级管理员<br/>super_admin]
        AD[管理员<br/>admin]
        UR[普通用户<br/>user]
    end

    subgraph "权限层"
        P1[权限1<br/>user:list]
        P2[权限2<br/>user:create]
        P3[权限3<br/>sms:send]
        P4[权限4<br/>file:upload]
        P_ALL[所有权限<br/>All Permissions]
    end

    subgraph "资源层"
        R1[API资源<br/>/api/user/list]
        R2[API资源<br/>/api/user/create]
        R3[API资源<br/>/api/sms/send]
        R4[API资源<br/>/api/file/upload]
    end

    subgraph "菜单层"
        M1[菜单1<br/>用户管理]
        M2[菜单2<br/>短信管理]
        M3[菜单3<br/>文件管理]
    end

    subgraph "数据库表"
        T1[AA07<br/>用户角色关联表]
        T2[AA08<br/>角色表]
        T3[AA09<br/>权限表]
        T4[AA10<br/>菜单表]
        T5[AA12<br/>角色菜单关联表]
        T6[AA13<br/>角色权限授予关系表]
        T7[casbin_rule<br/>Casbin策略表]
    end

    %% 用户-角色关联
    U1 -->|拥有| SA
    U2 -->|拥有| AD
    U3 -->|拥有| UR
    U1 -.->|存储| T1
    U2 -.->|存储| T1
    U3 -.->|存储| T1

    %% 角色-权限授权关系（层级授权）
    SA -->|授予所有权限| P_ALL
    SA -->|授予权限1,2| AD
    AD -->|授予权限1| UR
    SA -.->|记录授予关系| T6
    AD -.->|记录授予关系| T6

    %% 权限-资源映射
    P1 -->|控制| R1
    P2 -->|控制| R2
    P3 -->|控制| R3
    P4 -->|控制| R4
    P1 -.->|存储| T3
    P2 -.->|存储| T3
    P3 -.->|存储| T3
    P4 -.->|存储| T3

    %% Casbin策略
    SA -->|策略: p, super_admin, *, *| T7
    AD -->|策略: p, admin, /api/user/list, GET| T7
    UR -->|策略: p, user, /api/user/list, GET| T7

    %% 角色-菜单关联
    SA -->|关联所有菜单| M1
    SA -->|关联所有菜单| M2
    SA -->|关联所有菜单| M3
    AD -->|关联菜单1| M1
    AD -->|关联菜单2| M2
    AD -.->|存储| T5
    SA -.->|存储| T5

    %% 角色存储
    SA -.->|存储| T2
    AD -.->|存储| T2
    UR -.->|存储| T2

    %% 菜单存储
    M1 -.->|存储| T4
    M2 -.->|存储| T4
    M3 -.->|存储| T4

    style SA fill:#ff6b6b,stroke:#c92a2a,stroke-width:3px,color:#fff
    style AD fill:#4ecdc4,stroke:#2d9cdb,stroke-width:2px,color:#fff
    style UR fill:#95e1d3,stroke:#6c5ce7,stroke-width:2px,color:#fff
    style P_ALL fill:#feca57,stroke:#ff9ff3,stroke-width:2px
    style T6 fill:#ffd93d,stroke:#f6b93b,stroke-width:2px
    style T7 fill:#a29bfe,stroke:#6c5ce7,stroke-width:2px
```

## 权限控制流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant API as API接口
    participant Middleware as Casbin中间件
    participant Casbin as Casbin引擎
    participant DB as 数据库

    User->>Frontend: 1. 登录
    Frontend->>API: 2. POST /api/user/login
    API->>DB: 3. 验证用户信息
    DB-->>API: 4. 返回用户信息
    API->>DB: 5. 查询用户角色(AA07)
    DB-->>API: 6. 返回角色列表
    API->>API: 7. 生成JWT Token(包含用户ID和角色)
    API-->>Frontend: 8. 返回Token和用户信息
    Frontend->>Frontend: 9. 存储Token

    User->>Frontend: 10. 访问资源
    Frontend->>API: 11. 请求API(携带Token)
    API->>Middleware: 12. 进入Casbin中间件
    Middleware->>Middleware: 13. 解析Token获取用户ID和角色
    Middleware->>Casbin: 14. 检查权限: Enforce(userID, resource, method)
    Casbin->>DB: 15. 查询策略表(casbin_rule)
    DB-->>Casbin: 16. 返回策略规则
    Casbin->>Casbin: 17. 匹配策略: g(userID, role) && p(role, resource, method)
    
    alt 权限验证通过
        Casbin-->>Middleware: 18. 允许访问
        Middleware->>API: 19. 继续处理请求
        API->>DB: 20. 执行业务逻辑
        DB-->>API: 21. 返回数据
        API-->>Frontend: 22. 返回响应
    else 权限验证失败
        Casbin-->>Middleware: 18. 拒绝访问
        Middleware-->>Frontend: 19. 返回403 Forbidden
    end
```

## 层级授权机制

```mermaid
graph LR
    subgraph "层级授权流程"
        SA[超级管理员<br/>super_admin]
        AD[管理员<br/>admin]
        UR[普通用户<br/>user]
        
        P1[权限1]
        P2[权限2]
        P3[权限3]
        
        AA13[AA13表<br/>角色权限授予关系]
    end

    SA -->|1. 授予权限1,2给管理员| AD
    SA -->|2. 记录到AA13表| AA13
    AD -->|3. 授予权限1给普通用户| UR
    AD -->|4. 记录到AA13表| AA13
    
    SA -.->|拥有| P1
    SA -.->|拥有| P2
    SA -.->|拥有| P3
    AD -.->|被授予| P1
    AD -.->|被授予| P2
    UR -.->|被授予| P1

    style SA fill:#ff6b6b,stroke:#c92a2a,stroke-width:3px,color:#fff
    style AD fill:#4ecdc4,stroke:#2d9cdb,stroke-width:2px,color:#fff
    style UR fill:#95e1d3,stroke:#6c5ce7,stroke-width:2px,color:#fff
    style AA13 fill:#ffd93d,stroke:#f6b93b,stroke-width:2px
```

## 数据库表关系

```mermaid
erDiagram
    AA01 ||--o{ AA07 : "用户拥有角色"
    AA08 ||--o{ AA07 : "角色分配给用户"
    AA08 ||--o{ AA12 : "角色关联菜单"
    AA10 ||--o{ AA12 : "菜单关联角色"
    AA08 ||--o{ AA13 : "授予者角色"
    AA08 ||--o{ AA13 : "被授予者角色"
    AA09 ||--o{ AA13 : "权限被授予"
    AA08 ||--o{ casbin_rule : "角色策略"
    AA01 ||--o{ casbin_rule : "用户角色分配"

    AA01 {
        string AAB001 PK "用户ID"
        string AAB002 "手机号"
        string AAB003 "昵称"
    }

    AA07 {
        string AAG001 PK "用户ID"
        string AAG002 PK "角色ID"
    }

    AA08 {
        string AAH001 PK "角色ID"
        string AAH002 "角色Key"
        string AAH003 "角色名称"
    }

    AA09 {
        string AAI001 PK "权限ID"
        string AAI002 "权限Key"
        string AAI003 "权限名称"
        string AAI004 "资源路径"
        string AAI005 "请求方法"
    }

    AA10 {
        string AAI001 PK "菜单ID"
        string AAI002 "菜单Key"
        string AAI003 "菜单标题"
        string AAI004 "菜单路径"
    }

    AA12 {
        string AAL001 PK "角色ID"
        string AAL002 PK "菜单ID"
    }

    AA13 {
        string AAM001 "授予者角色ID"
        string AAM002 "被授予者角色ID"
        string AAM003 "权限ID"
    }

    casbin_rule {
        string ptype "策略类型"
        string v0 "角色Key/用户ID"
        string v1 "资源路径"
        string v2 "请求方法"
    }
```

## 权限检查逻辑

```mermaid
flowchart TD
    Start([用户请求API]) --> CheckToken{检查JWT Token}
    CheckToken -->|Token无效| Reject1[返回401未授权]
    CheckToken -->|Token有效| ParseToken[解析Token获取用户ID和角色]
    
    ParseToken --> CheckSuperAdmin{是否超级管理员?}
    CheckSuperAdmin -->|是| Allow[直接放行]
    CheckSuperAdmin -->|否| GetRoles[从Casbin获取用户角色]
    
    GetRoles --> CheckCasbin{Casbin检查权限}
    CheckCasbin -->|策略匹配| Allow
    CheckCasbin -->|策略不匹配| Reject2[返回403禁止访问]
    
    Allow --> ProcessRequest[处理业务请求]
    ProcessRequest --> ReturnResponse[返回响应]
    
    Reject1 --> End([结束])
    Reject2 --> End
    ReturnResponse --> End

    style CheckSuperAdmin fill:#ff6b6b,stroke:#c92a2a,stroke-width:2px,color:#fff
    style Allow fill:#51cf66,stroke:#2f9e44,stroke-width:2px,color:#fff
    style Reject1 fill:#ff8787,stroke:#c92a2a,stroke-width:2px,color:#fff
    style Reject2 fill:#ff8787,stroke:#c92a2a,stroke-width:2px,color:#fff
```

## 菜单权限控制

```mermaid
graph TB
    subgraph "菜单获取流程"
        User[用户登录] --> GetRoles[获取用户角色]
        GetRoles --> CheckSuperAdmin{是否超级管理员?}
        
        CheckSuperAdmin -->|是| ReturnAllMenus[返回所有菜单]
        CheckSuperAdmin -->|否| GetRoleMenus[查询角色菜单关联AA12]
        
        GetRoleMenus --> FilterMenus[根据角色过滤菜单]
        FilterMenus --> CheckParentMenu{父菜单是否授权?}
        
        CheckParentMenu -->|是| IncludeParent[包含父菜单]
        CheckParentMenu -->|否| CheckChildren{子菜单是否授权?}
        
        CheckChildren -->|是| IncludeParentWithChildren[包含父菜单和子菜单]
        CheckChildren -->|否| ExcludeMenu[排除菜单]
        
        IncludeParent --> ReturnFilteredMenus[返回过滤后的菜单]
        IncludeParentWithChildren --> ReturnFilteredMenus
        ReturnAllMenus --> End([结束])
        ReturnFilteredMenus --> End
        ExcludeMenu --> ReturnFilteredMenus
    end

    style CheckSuperAdmin fill:#ff6b6b,stroke:#c92a2a,stroke-width:2px,color:#fff
    style ReturnAllMenus fill:#51cf66,stroke:#2f9e44,stroke-width:2px,color:#fff
    style ReturnFilteredMenus fill:#51cf66,stroke:#2f9e44,stroke-width:2px,color:#fff
```

## 核心组件说明

### 1. 角色体系
- **超级管理员（super_admin）**：拥有所有权限，可以给任何角色授权
- **管理员（admin）**：只能看到被授予的权限，可以给普通用户授权
- **普通用户（user）**：只能看到被授予的权限，不能授权

### 2. 数据库表
- **AA07**：用户角色关联表，存储用户和角色的多对多关系
- **AA08**：角色表，存储角色基本信息
- **AA09**：权限表，存储权限定义
- **AA10**：菜单表，存储菜单定义
- **AA12**：角色菜单关联表，存储角色和菜单的多对多关系
- **AA13**：角色权限授予关系表，存储层级授权关系（授予者、被授予者、权限）
- **casbin_rule**：Casbin策略表，存储权限规则（p策略和g策略）

### 3. 权限控制流程
1. 用户登录 → 生成JWT Token（包含用户ID和角色）
2. 请求API → Casbin中间件验证权限
3. 权限检查 → Casbin引擎匹配策略
4. 授权决策 → 允许或拒绝访问

### 4. 层级授权机制
- 超级管理员可以授予任何权限给任何角色
- 管理员只能授予自己被授予的权限给普通用户
- 所有授权关系都记录在AA13表中，确保可追溯

### 5. 菜单权限控制
- 超级管理员可以看到所有菜单
- 其他角色只能看到被授权的菜单
- 如果子菜单被授权，父菜单会自动显示（作为容器）

## 技术栈

- **权限引擎**：Casbin (RBAC模型)
- **策略存储**：MySQL (GORM适配器)
- **认证方式**：JWT Token
- **中间件**：Gin Casbin中间件
- **数据库**：MySQL

