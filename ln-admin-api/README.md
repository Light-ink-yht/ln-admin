# LN Admin API

基于 Go + Gin + GORM 构建的现代化后台管理系统后端服务，采用 DDD（领域驱动设计）架构和 RBAC 权限控制。

## 📚 文档导航

- 📖 [快速开始指南](./docs/START.md) - 详细的配置和启动说明
- 🔐 [Casbin 权限控制](./docs/CASBIN.md) - 权限管理系统使用指南
- 📱 [短信服务配置](./docs/SMS.md) - 腾讯云短信服务集成指南（即将发布）
- 📖 [Swagger API 文档使用指南](./docs/SWAGGER.md) - Swagger 文档生成和使用说明

## ✨ 核心特性

### 架构特性

- 🏗️ **DDD 架构** - 清晰的领域驱动设计，业务与技术分离
- 📦 **模块化设计** - 高内聚低耦合，易于扩展和维护
- 🎯 **单一职责** - 每个模块职责明确，符合SOLID原则

### 安全特性

- 🔐 **双Token认证** - AccessToken + RefreshToken 机制，提升安全性
- 🛡️ **RBAC权限控制** - 基于 Casbin 的细粒度权限管理
- 🔒 **密码加密** - 使用 bcrypt 单向哈希，不可逆加密
- 🔐 **JWT认证** - 基于JWT的无状态认证机制
- 🚫 **防刷机制** - 验证码频率限制，Redis + Lua 原子操作

### 功能特性

- 📱 **腾讯云短信服务** - 集成短信验证码功能，支持注册、找回密码
- 🖼️ **图片验证码** - 防止机器人攻击
- 📝 **数据库日志系统** - 完整的操作日志记录，便于审计和运维
- ⚙️ **动态配置管理** - 系统配置存储在数据库，支持热更新无需重启
- 🔄 **统一响应格式** - 标准化的 API 响应结构
- 🛡️ **完善的中间件** - CORS、认证、日志、错误恢复等

### 技术特性

- 🗄️ **多数据库支持** - 兼容 MySQL 和 Oracle 数据库
- 📦 **GORM AutoMigrate** - 自动数据库表结构迁移
- 🚀 **Redis缓存** - 验证码和会话缓存，高性能
- 📊 **结构化日志** - 基于 Zap 的结构化日志系统
- 🔄 **优雅关闭** - 支持优雅关闭，不丢失请求

## 📁 项目结构

```
ln-admin-api/
├── cmd/                                 # 应用入口
│   └── server/
│       ├── main.go                     # 主函数，启动流程
│       └── bootstrap.go                # 启动器，封装初始化逻辑
├── internal/                            # 内部代码（不对外暴露）
│   ├── domain/                         # 领域层（核心业务逻辑）
│   │   ├── entity/                     # 领域实体
│   │   │   ├── user.go                 # 用户实体（AA01）
│   │   │   ├── system_config.go         # 系统配置实体（AA02）
│   │   │   ├── system_log.go           # 系统日志实体（AA03）
│   │   │   ├── sms_code.go             # 短信验证码实体（AA04）
│   │   │   ├── role.go                 # 角色实体（AA05）
│   │   │   ├── permission.go           # 权限实体（AA06）
│   │   │   ├── user_role.go            # 用户角色关联实体（AA07）
│   │   │   └── sms_template.go        # 短信模板实体（AA08）
│   │   └── repository/                 # 仓库接口（领域层定义）
│   │       ├── user_repository.go
│   │       ├── system_config_repository.go
│   │       ├── system_log_repository.go
│   │       ├── sms_code_repository.go
│   │       ├── role_repository.go
│   │       ├── permission_repository.go
│   │       ├── user_role_repository.go
│   │       └── sms_template_repository.go
│   ├── application/                     # 应用层（业务用例编排）
│   │   ├── dto/                        # 数据传输对象
│   │   │   └── user_dto.go            # 用户相关DTO
│   │   └── service/                   # 应用服务
│   │       ├── user_app_service.go    # 用户应用服务
│   │       ├── sms_app_service.go     # 短信应用服务
│   │       └── permission_service.go  # 权限管理服务
│   ├── infrastructure/                 # 基础设施层（技术实现）
│   │   ├── database/                  # 数据库连接
│   │   │   └── db.go
│   │   ├── redis/                      # Redis客户端
│   │   │   ├── redis.go               # Redis初始化
│   │   │   └── sms_code_store.go      # 短信验证码存储（Lua脚本）
│   │   ├── casbin/                     # Casbin权限控制
│   │   │   └── casbin.go              # Casbin初始化
│   │   ├── logger/                     # 日志系统
│   │   │   └── logger.go
│   │   ├── repository/                 # 仓库实现
│   │   │   ├── user_repository_impl.go
│   │   │   ├── system_config_repository_impl.go
│   │   │   ├── system_log_repository_impl.go
│   │   │   ├── sms_code_repository_impl.go
│   │   │   ├── role_repository_impl.go
│   │   │   ├── permission_repository_impl.go
│   │   │   ├── user_role_repository_impl.go
│   │   │   └── sms_template_repository_impl.go
│   │   └── defaultdata/                # 默认数据初始化
│   │       └── init_data.go           # 默认用户、角色、权限、模板
│   └── interfaces/                     # 接口层（对外接口）
│       └── http/                       # HTTP接口
│           ├── handler/                # 处理器
│           │   └── user_handler.go
│           ├── middleware/             # 中间件
│           │   ├── auth.go             # JWT认证
│   │   │   ├── casbin_middleware.go    # Casbin权限验证
│   │   │   ├── cors.go                 # 跨域
│   │   │   ├── logger.go              # 日志
│   │   │   └── recovery.go            # 错误恢复
│   │   └── router/                     # 路由配置
│   │       └── router.go
├── pkg/                                 # 公共包（可被外部引用）
│   ├── config/                         # 配置管理
│   │   └── config.go
│   ├── response/                       # 统一响应封装
│   │   └── response.go
│   ├── service/                        # 公共服务
│   │   ├── sms/                        # 短信服务
│   │   │   └── tencent_sms.go         # 腾讯云短信实现
│   │   └── casbin_service/            # Casbin服务封装
│   │       └── casbin_service.go
│   ├── utils/                          # 工具函数
│   │   ├── jwt.go                      # JWT工具（双Token）
│   │   └── password.go                 # 密码加密工具
│   └── validator/                      # 验证器
│       └── validator.go
├── configs/                             # 配置文件
│   ├── config.yaml                      # 应用配置
│   └── rbac_model.conf                  # Casbin RBAC模型
├── docs/                                # 文档目录
│   ├── START.md                         # 启动指南
│   ├── CASBIN.md                        # Casbin权限控制指南
│   └── ...                             # 其他文档
├── scripts/                              # 数据库脚本
│   ├── migrate.sql                      # 初始迁移脚本
│   └── migrate_v2.sql                   # V2迁移脚本
└── README.md                            # 项目说明
```

## 🏗️ 架构设计

### DDD（领域驱动设计）分层架构

本项目采用经典的 DDD 四层架构，清晰地分离业务逻辑和技术实现：

#### 1. 领域层（Domain Layer）

**职责**：包含核心业务逻辑和业务规则

- **实体（Entity）**：领域模型的核心，包含业务标识和业务行为
  - `User`：用户领域实体，包含用户状态判断、密码验证等业务方法
  - `SystemConfig`：系统配置实体，管理配置的生命周期
  - `SystemLog`：系统日志实体，记录操作审计信息
  - `SMSCode`：短信验证码实体，包含过期判断、使用状态等业务逻辑

- **仓库接口（Repository Interface）**：定义数据访问的抽象接口，不依赖具体实现
  - 领域层只定义接口，具体实现在基础设施层
  - 便于测试和切换数据源

**设计原则**：
- 领域实体不依赖任何技术框架
- 业务逻辑集中在领域实体中
- 通过接口定义与外部交互

#### 2. 应用层（Application Layer）

**职责**：编排领域对象完成业务用例，协调不同领域服务

- **应用服务（Application Service）**：
  - `UserAppService`：用户相关的业务用例
    - 登录、注册、密码重置等完整业务流程
    - 协调用户仓库、短信服务等
  - `SMSAppService`：短信业务用例
    - 发送验证码、验证验证码
    - 防刷机制、频率限制

- **DTO（数据传输对象）**：用于接口层和应用层之间的数据传输
  - 屏蔽领域实体的内部细节
  - 适配前端需求格式

**设计原则**：
- 无状态服务，方法幂等
- 一个应用服务方法对应一个业务用例
- 不包含业务规则，只编排流程

#### 3. 基础设施层（Infrastructure Layer）

**职责**：提供技术实现，支撑上层业务

- **数据库实现**：
  - 支持 MySQL 和 Oracle
  - 连接池管理、事务处理
  - 自动适配不同数据库的特性

- **仓库实现**：
  - 实现领域层定义的仓库接口
  - 使用 GORM 进行 ORM 映射
  - 处理数据库实体的转换

- **日志系统**：
  - 基于 Zap 的结构化日志
  - 支持文件输出和日志轮转
  - 数据库日志记录

- **外部服务集成**：
  - 腾讯云短信服务封装
  - 统一错误处理和重试机制

**设计原则**：
- 技术细节与业务逻辑隔离
- 支持多种实现方式切换
- 提供统一的抽象接口

#### 4. 接口层（Interfaces Layer）

**职责**：处理外部请求，适配外部协议

- **HTTP 处理器**：
  - 参数验证、请求解析
  - 调用应用服务
  - 格式化响应

- **中间件**：
  - 认证：JWT Token 验证
  - CORS：跨域处理
  - 日志：请求日志记录
  - 恢复：Panic 捕获

- **路由配置**：
  - RESTful API 路由定义
  - 路由分组和权限控制

**设计原则**：
- 薄控制器，只负责协议适配
- 业务逻辑在应用层和领域层
- 统一的错误处理和响应格式

## 🗄️ 数据库设计

### 表命名规范

项目采用统一的表命名规范，使用 `AA + 2位数字` 格式：

| 表名 | 说明 | 实体名称 |
|------|------|---------|
| AA01 | 用户表 | User |
| AA02 | 系统配置表 | SystemConfig |
| AA03 | 系统日志表 | SystemLog |
| AA04 | 短信验证码表 | SMSCode |
| AA05 | 角色表 | Role |
| AA06 | 权限表 | Permission |
| AA07 | 用户角色关联表 | UserRole |
| AA08 | 短信模板表 | SMSTemplate |
| casbin_rule | Casbin策略表 | （自动创建）|

**字段命名规范**：`AAA + 3位数字`，如 `AAA001`、`AAA002`，便于识别和索引管理。

### 数据表设计

#### AA01 - 用户表

```sql
AAA001  VARCHAR(50)  PRIMARY KEY  用户ID（业务主键）
AAA002  VARCHAR(100) UNIQUE       邮箱（全局唯一）
AAA003  VARCHAR(20)  UNIQUE       手机号（全局唯一）
AAA004  VARCHAR(255) NOT NULL     密码（加密存储）
AAA005  VARCHAR(50)               昵称
AAA006  VARCHAR(50)               姓名
AAA007  VARCHAR(255)              头像
AAA008  VARCHAR(10)               性别（1 男，2 女，3 未知）
AAA009  DATE                      生日
AAA010  VARCHAR(10)               状态（1 启用，2 禁用）
AAA011  VARCHAR(500)              备注
AAA012  INT                       登录次数
AAA013  DATETIME    INDEX         最后登录时间
AAA014  VARCHAR(50)               最后登录IP
AAA015  DATETIME                  密码修改时间
AAA016  VARCHAR(50)               创建人
AAA017  VARCHAR(50)               修改人
```

**设计要点**：
- 使用业务主键（UUID）而非自增ID，便于分布式部署
- 邮箱和手机号建立唯一索引，保证数据完整性
- 支持软删除，保留历史数据
- 记录完整的审计信息（创建人、修改人等）

#### AA02 - 系统配置表

```sql
AAA021  VARCHAR(50)  PRIMARY KEY  配置键（业务主键）
AAA022  VARCHAR(500)              配置值
AAA023  VARCHAR(100)              配置名称
AAA024  VARCHAR(50)   INDEX       配置分组（sms, jwt, system等）
AAA025  VARCHAR(500)              配置描述
AAA026  VARCHAR(10)               状态（1 启用，2 禁用）
AAA027  DATETIME                  最后修改时间
AAA028  VARCHAR(50)               创建人
AAA029  VARCHAR(50)               修改人
```

**设计要点**：
- 配置键作为主键，便于快速查找
- 按分组管理配置，支持配置分类
- 支持配置的启用/禁用，不影响业务逻辑
- 所有配置存储在数据库，支持热更新

**典型配置项**：
- 短信服务配置（SecretId、SecretKey、AppId等）
- JWT 配置（过期时间等）
- 系统参数配置

#### AA03 - 系统日志表

```sql
AAA031  VARCHAR(50)               日志ID（UUID）
AAA032  VARCHAR(20)   INDEX       日志级别（debug, info, warn, error）
AAA033  VARCHAR(200)              模块名称
AAA034  VARCHAR(100)              操作类型（login, register等）
AAA035  TEXT                      日志内容
AAA036  VARCHAR(50)               用户ID
AAA037  VARCHAR(50)               IP地址
AAA038  VARCHAR(500)              请求路径
AAA039  VARCHAR(20)               请求方法（GET, POST等）
AAA040  INT                       响应状态码
AAA041  VARCHAR(500)              用户代理
AAA042  VARCHAR(200)              错误信息
AAA043  DATETIME     INDEX        日志时间
AAA044  VARCHAR(50)               创建人
```

**设计要点**：
- 完整的请求上下文记录，便于问题追踪
- 日志级别和日志时间建立索引，支持快速查询
- 支持按模块、操作类型、用户等多维度查询
- 文本类型字段使用 TEXT，支持长内容记录

#### AA04 - 短信验证码表（数据库备份，主要使用Redis存储）

```sql
AAA041  VARCHAR(50)  PRIMARY KEY  验证码ID（UUID）
AAA042  VARCHAR(20)  INDEX        手机号
AAA043  VARCHAR(10)               验证码（6位数字）
AAA044  VARCHAR(20)               验证码类型（register, forgot, login）
AAA045  VARCHAR(10)               状态（1 未使用，2 已使用，3 已过期）
AAA046  DATETIME                  过期时间
AAA047  DATETIME                  使用时间
AAA048  VARCHAR(50)               IP地址
AAA049  INT                       发送次数
```

**设计要点**：
- **主要存储**：验证码优先存储在Redis中，10分钟过期
- **原子操作**：使用Redis Lua脚本保证验证码验证和标记的原子性
- **防重复使用**：通过Lua脚本原子性地标记已使用状态
- **防刷机制**：使用Lua脚本原子性检查发送频率（1分钟内只能发送一次）
- **记录IP地址**：记录发送IP，便于安全审计
- **数据库备份**：表结构保留用于历史查询和审计

## 🔐 安全设计

### 双Token机制

#### AccessToken（访问令牌）

- **用途**：用于API请求认证
- **有效期**：2小时（可配置）
- **存储位置**：客户端内存或本地存储
- **刷新策略**：过期后使用 RefreshToken 刷新

#### RefreshToken（刷新令牌）

- **用途**：用于刷新 AccessToken
- **有效期**：7天（可配置）
- **存储位置**：客户端安全存储（HttpOnly Cookie 或 SecureStorage）
- **安全策略**：
  - 一次性使用：每次刷新都生成新的 Token 对
  - 刷新后旧的 RefreshToken 失效
  - 支持 Token 撤销机制

#### Token生成流程

```
用户登录
  ↓
生成 AccessToken + RefreshToken
  ↓
返回给客户端
  ↓
客户端存储 Token
  ↓
API 请求携带 AccessToken
  ↓
AccessToken 过期
  ↓
使用 RefreshToken 刷新
  ↓
获得新的 Token 对
```

#### 安全特性

- **JWT签名**：使用 HS256 算法，密钥存储在配置文件
- **Token验证**：每次请求验证 Token 的有效性和签名
- **过期检查**：自动检查 Token 过期时间
- **用户绑定**：Token 中绑定用户ID，防止 Token 被他人使用

### 密码安全

- **加密算法**：使用 bcrypt 单向哈希
- **加密强度**：可配置的 cost 参数（默认10）
- **密码策略**：
  - 最小长度：6位
  - 最大长度：20位
  - 建议包含字母、数字、特殊字符

### 短信验证码安全

- **存储方式**：使用Redis存储，高性能，自动过期
- **验证码生成**：6位随机数字
- **有效期**：10分钟（600秒）
- **原子操作**：使用Redis Lua脚本保证验证和标记操作的原子性，防止并发问题
- **防刷机制**：
  - 使用Lua脚本原子性检查，同一手机号1分钟内只能发送一次
  - 验证码使用后立即失效（通过Lua脚本原子性标记）
  - 记录发送IP，便于安全审计
- **验证流程**：
  - 发送验证码前验证图片验证码（防止机器人）
  - 验证码存储在Redis，10分钟自动过期
  - 数据库表保留用于历史查询和审计
- **模板配置**：不同场景（注册、登录、忘记密码）使用不同的短信模板ID，通过数据库配置管理

## 📝 日志系统设计

### 日志分级

- **Debug**：调试信息，开发环境使用
- **Info**：一般信息，记录正常业务流程
- **Warn**：警告信息，潜在问题但不影响功能
- **Error**：错误信息，需要关注的问题

### 日志输出方式

1. **文件输出**：
   - 支持日志轮转（按大小、按时间）
   - 自动压缩历史日志
   - 配置保留天数

2. **数据库存储**：
   - 存储在 AA03 表
   - 记录完整的请求上下文
   - 支持多维度查询和分析

### 日志记录内容

- **请求信息**：路径、方法、参数
- **响应信息**：状态码、响应时间
- **用户信息**：用户ID、IP地址
- **错误信息**：错误堆栈、错误类型
- **业务信息**：操作类型、模块名称

### 日志中间件

自动记录所有HTTP请求：
- 记录请求开始时间
- 记录请求参数和响应结果
- 记录处理耗时
- 异常情况自动记录错误信息

## ⚙️ 配置管理设计

### 配置文件

使用 YAML 格式的配置文件，包含：
- 应用配置（端口、模式等）
- 数据库配置（支持MySQL和Oracle）
- JWT配置（密钥、过期时间等）
- 日志配置（级别、输出方式等）
- 中间件配置（CORS、认证等）

### 数据库配置

系统配置存储在 AA02 表，支持：

1. **动态配置**：
   - 无需重启服务即可更新配置
   - 配置修改后立即生效
   - 支持配置的启用/禁用

2. **配置分组**：
   - sms：短信服务配置
   - jwt：JWT配置
   - system：系统参数配置

3. **配置管理**：
   - 通过接口查询配置
   - 支持按分组查询
   - 配置修改记录审计信息

### 配置优先级

1. 数据库配置（最高优先级）
2. 环境变量
3. 配置文件（默认值）

## 🔄 业务流程设计

### 用户注册流程

```
1. 获取图片验证码
   ↓
2. 发送短信验证码（type: register）
   ↓
3. 提交注册信息（手机号、密码、短信验证码）
   ↓
4. 验证短信验证码
   ↓
5. 检查手机号是否已存在
   ↓
6. 创建用户记录
   ↓
7. 返回用户信息
```

### 用户登录流程

```
1. 获取图片验证码
   ↓
2. 发送登录短信验证码（type: login）
   ↓
3. 提交登录信息（手机号、密码、图片验证码、短信验证码）
   ↓
4. 验证图片验证码
   ↓
5. 验证短信验证码
   ↓
6. 验证用户密码
   ↓
7. 检查用户状态
   ↓
8. 生成双Token
   ↓
9. 更新登录信息（次数、时间、IP）
   ↓
10. 返回Token和用户信息
```

### Token刷新流程

```
1. 检测到 AccessToken 过期
   ↓
2. 使用 RefreshToken 调用刷新接口
   ↓
3. 验证 RefreshToken 有效性
   ↓
4. 生成新的双Token
   ↓
5. 返回新的Token对
   ↓
6. 客户端更新存储的Token
```

## 🗃️ 数据持久化设计

### 实体映射策略

使用领域实体和数据库实体分离的设计：

- **领域实体**（Domain Entity）：
  - 面向业务，字段名语义清晰
  - 包含业务方法（如 `IsActive()`、`IsExpired()`）
  - 不依赖数据库框架

- **数据库实体**（Database Entity）：
  - 面向数据库，字段名遵循数据库规范（AAA + 数字）
  - 包含 GORM 标签和索引定义
  - 负责数据库映射

### 转换机制

通过 `ToXXX()` 和 `FromXXX()` 方法实现双向转换：
- 业务层使用领域实体
- 数据层使用数据库实体
- 转换逻辑集中在实体中，便于维护

### 软删除

所有表支持软删除：
- 使用 GORM 的 `DeletedAt` 字段
- 删除操作只标记，不真正删除数据
- 查询时自动过滤已删除记录
- 支持数据恢复和历史追溯

## 🚀 快速开始

### 最小化启动步骤

1. **安装依赖**
   ```bash
   go mod download
   ```

2. **配置数据库和Redis**
   - 修改 `configs/config.yaml` 中的数据库和Redis配置

3. **启动服务**
   ```bash
   go run cmd/server/main.go
   ```

4. **使用默认账号登录**
   - 手机号：`18797131041`
   - 密码：`hello#123world`

### 详细启动指南

完整的环境配置、数据库初始化、启动流程和故障排查，请查看：

👉 **[启动指南文档](./docs/START.md)**

## 🎯 权限管理

### RBAC 权限体系

项目使用 Casbin 实现基于角色的访问控制（RBAC）：

- **用户（User）**：系统使用者
- **角色（Role）**：权限的集合（如：超级管理员、管理员、普通用户）
- **权限（Permission）**：对资源的操作（如：查看用户列表、删除用户）
- **资源（Resource）**：API路径（如：`/api/user/list`）

### 权限控制流程

```
用户登录 → JWT认证 → 获取用户角色 → Casbin验证权限 → 允许/拒绝访问
```

### 权限管理指南

完整的权限管理使用方法，包括：
- 角色创建和分配
- 权限创建和分配
- 权限验证机制
- 策略查询和管理

请查看：

👉 **[Casbin 权限控制指南](./docs/CASBIN.md)**

### 运维建议

1. **日志管理**：
   - 定期清理历史日志文件
   - 数据库日志按时间分区存储
   - 重要日志定期备份

2. **性能优化**：
   - 数据库连接池参数调优
   - 日志记录异步化（可选）
   - 验证码缓存优化（可选Redis）

3. **安全建议**：
   - 定期更换JWT密钥
   - 监控异常登录行为
   - 配置HTTPS协议
   - 限制API访问频率

4. **监控指标**：
   - 请求量和响应时间
   - 错误率和异常日志
   - 数据库连接池状态
   - Token使用情况

## 📖 设计原则和最佳实践

### 代码组织原则

1. **依赖方向**：接口层 → 应用层 → 领域层 ← 基础设施层
2. **单一职责**：每个服务、每个方法只做一件事
3. **开闭原则**：对扩展开放，对修改关闭
4. **依赖倒置**：依赖抽象而非具体实现

### 错误处理原则

1. **统一错误码**：定义清晰的错误码体系
2. **错误传播**：错误信息明确，便于定位问题
3. **优雅降级**：非关键功能失败不影响主流程
4. **错误日志**：所有错误都记录到日志系统

### 性能优化原则

1. **连接池管理**：合理配置数据库连接池参数
2. **索引优化**：为常用查询字段建立索引
3. **批量操作**：减少数据库交互次数
4. **缓存策略**：对热点数据使用缓存

### 安全设计原则

1. **最小权限**：每个用户只拥有必要的权限
2. **输入验证**：所有外部输入都进行验证
3. **密码安全**：使用强加密算法，不存储明文
4. **审计日志**：记录所有重要操作，便于追溯

## 🔧 扩展指南

### 添加新的业务模块

遵循 DDD 分层架构，按以下步骤添加新模块：

1. **领域层（Domain）**
   ```bash
   # 1. 创建领域实体
   internal/domain/entity/your_entity.go
   
   # 2. 定义仓库接口
   internal/domain/repository/your_repository.go
   ```

2. **基础设施层（Infrastructure）**
   ```bash
   # 3. 实现仓库
   internal/infrastructure/repository/your_repository_impl.go
   ```

3. **应用层（Application）**
   ```bash
   # 4. 创建应用服务
   internal/application/service/your_app_service.go
   
   # 5. 定义DTO
   internal/application/dto/your_dto.go
   ```

4. **接口层（Interfaces）**
   ```bash
   # 6. 创建处理器
   internal/interfaces/http/handler/your_handler.go
   
   # 7. 配置路由
   internal/interfaces/http/router/router.go
   ```

5. **数据库迁移**
   ```bash
   # 8. 在 bootstrap.go 中添加表迁移
   db.AutoMigrate(&entity.YourEntity{})
   ```

### 集成新的外部服务

1. **创建服务封装**
   ```go
   // pkg/service/your_service/your_service.go
   type YourService struct {
       // 依赖注入
   }
   
   func (s *YourService) DoSomething() error {
       // 实现逻辑
   }
   ```

2. **配置管理**
   - 在 `AA02` 系统配置表中存储服务配置
   - 使用配置分组管理：`configRepo.GetByGroup(ctx, "your_service")`

3. **错误处理**
   ```go
   var (
       ErrYourServiceError = errors.New("服务错误")
   )
   ```

4. **日志记录**
   - 使用 `logger` 包记录服务调用
   - 记录请求参数、响应结果、错误信息

### 自定义中间件

1. **创建中间件**
   ```go
   // internal/interfaces/http/middleware/your_middleware.go
   func YourMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // 中间件逻辑
           c.Next()
       }
   }
   ```

2. **注册中间件**
   ```go
   // internal/interfaces/http/router/router.go
   r.Use(middleware.YourMiddleware())
   ```

### 添加新的权限

1. **创建权限实体**
   ```go
   permission := &entity.Permission{
       PermissionID:   "perm_your_action",
       PermissionKey:  "your:action",
       PermissionName: "您的操作",
       ResourcePath:   "/api/your/resource",
       Method:         "GET",
       Status:         "1",
   }
   ```

2. **为角色分配权限**
   ```go
   permissionService.AssignPermissionToRole(ctx, roleID, permissionID)
   ```

3. **在路由中使用权限验证**
   ```go
   auth.Use(middleware.CasbinMiddleware())
   ```

## 📚 更多文档

- 📖 [启动指南](./docs/START.md) - 详细的配置和启动说明
- 🔐 [Casbin 权限控制](./docs/CASBIN.md) - 完整的权限管理指南
- 📖 [Swagger API 文档使用指南](./docs/SWAGGER.md) - Swagger 文档生成和使用说明
- 📱 [短信服务配置](./docs/SMS.md) - 腾讯云短信集成（即将发布）

### 在线 API 文档

启动服务后，访问 Swagger UI 查看完整的 API 文档：

```
http://localhost:8080/swagger/index.html
```

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📄 License

MIT License
