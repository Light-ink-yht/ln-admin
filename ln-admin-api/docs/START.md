# 启动指南

本文档详细介绍如何配置和启动 LN Admin API 服务。

## 📋 前置要求

### 环境依赖

- **Go**: >= 1.21
- **数据库**: 
  - MySQL >= 5.7
  - 或 Oracle >= 11g
- **Redis**: >= 6.0（用于验证码存储和会话管理）
- **操作系统**: Linux、macOS、Windows

### 开发工具推荐

- IDE: VS Code、GoLand
- 数据库管理工具: Navicat、DBeaver、DataGrip
- API 测试工具: Postman、Apifox

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd ln-admin-api
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置文件设置

#### 3.1 复制配置文件

```bash
cp configs/config.yaml configs/config.yaml
```

如果 `config.yaml` 不存在，需要从示例文件创建：

```bash
# 创建配置文件
cat > configs/config.yaml <<EOF
# 应用配置
app:
  name: ln-admin
  version: 1.0.0
  port: 8080
  mode: debug # debug, release, test

# 数据库配置
database:
  type: mysql # mysql 或 oracle
  mysql:
    host: localhost
    port: 3306
    user: root
    password: your_password
    dbname: ln_admin
    charset: utf8mb4
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 3600
    conn_max_idle_time: 300

# JWT配置
jwt:
  secret: your-secret-key-change-in-production
  expire: 7200 # 秒，2小时
  refresh_expire: 604800 # 秒，7天
  issuer: ln-admin

# 日志配置
log:
  level: info
  format: json
  output: file
  path: ./logs
  filename: app.log
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true

# Redis配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3

# 短信配置
sms:
  enabled: false # false=开发模式（仅日志），true=生产模式（需要完整配置）
EOF
```

#### 3.2 配置说明

##### 应用配置（app）

```yaml
app:
  name: ln-admin          # 应用名称
  version: 1.0.0          # 版本号
  port: 8080              # 服务端口
  mode: debug             # 运行模式：debug（开发）、release（生产）、test（测试）
```

##### 数据库配置（database）

**MySQL 配置示例：**

```yaml
database:
  type: mysql
  mysql:
    host: localhost                  # 数据库主机
    port: 3306                       # 端口
    user: root                       # 用户名
    password: your_password          # 密码
    dbname: ln_admin                 # 数据库名
    charset: utf8mb4                 # 字符集
    max_open_conns: 100              # 最大打开连接数
    max_idle_conns: 10               # 最大空闲连接数
    conn_max_lifetime: 3600          # 连接最大存活时间（秒）
    conn_max_idle_time: 300          # 连接最大空闲时间（秒）
```

**Oracle 配置示例：**

```yaml
database:
  type: oracle
  oracle:
    host: localhost
    port: 1521
    user: system
    password: your_password
    service_name: ORCL               # Oracle 服务名
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 3600
    conn_max_idle_time: 300
```

##### JWT 配置（jwt）

```yaml
jwt:
  secret: your-secret-key-change-in-production  # JWT签名密钥（生产环境必须更改！）
  expire: 7200                                   # AccessToken过期时间（秒），默认2小时
  refresh_expire: 604800                        # RefreshToken过期时间（秒），默认7天
  issuer: ln-admin                              # Token发行者
```

⚠️ **重要**：生产环境必须修改 `secret`，使用强随机字符串，建议至少32字符。

##### Redis 配置（redis）

```yaml
redis:
  host: localhost        # Redis主机
  port: 6379            # 端口
  password: ""           # 密码（无密码留空）
  db: 0                  # 数据库编号（0-15）
  pool_size: 10          # 连接池大小
  min_idle_conns: 5      # 最小空闲连接数
  dial_timeout: 5        # 连接超时（秒）
  read_timeout: 3        # 读超时（秒）
  write_timeout: 3       # 写超时（秒）
```

##### 短信配置（sms）

```yaml
sms:
  enabled: false  # false=开发模式（验证码仅在日志中显示），true=生产模式（实际发送短信）
```

### 4. 数据库初始化

#### 4.1 创建数据库

**MySQL：**

```sql
CREATE DATABASE ln_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**Oracle：**

```sql
CREATE USER ln_admin IDENTIFIED BY your_password;
GRANT CONNECT, RESOURCE, DBA TO ln_admin;
```

#### 4.2 自动迁移（推荐）

项目支持 GORM AutoMigrate，启动时会自动创建表结构：

```bash
go run cmd/server/main.go
```

启动时会在日志中看到：
```
[INFO] 数据库表结构迁移成功
```

#### 4.3 手动执行 SQL（可选）

如果希望手动控制表结构，可以执行 SQL 脚本：

**MySQL：**

```bash
mysql -u root -p ln_admin < scripts/migrate.sql
mysql -u root -p ln_admin < scripts/migrate_v2.sql
```

### 5. 启动服务

#### 5.1 开发模式启动

```bash
# 直接运行
go run cmd/server/main.go

# 或使用 air（热重载，需要先安装：go install github.com/cosmtrek/air@latest）
air
```

#### 5.2 生产模式启动

```bash
# 编译
go build -o bin/ln-admin-api cmd/server/main.go

# 运行
./bin/ln-admin-api

# 或使用 systemd（Linux）
sudo systemctl start ln-admin-api
```

#### 5.3 Docker 启动（可选）

```dockerfile
# Dockerfile 示例
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ln-admin-api cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ln-admin-api .
COPY --from=builder /app/configs ./configs
CMD ["./ln-admin-api"]
```

### 6. 验证启动

#### 6.1 检查服务状态

```bash
# 检查端口是否监听
netstat -tuln | grep 8080

# 或使用 curl
curl http://localhost:8080/api/user/captcha
```

#### 6.2 查看日志

日志文件位于 `./logs/app.log`（如果配置了文件输出）

```bash
tail -f logs/app.log
```

#### 6.3 默认账号

系统启动时会自动创建默认管理员账号：

- **手机号**: `18797131041`
- **密码**: `hello#123world`
- **角色**: 超级管理员（拥有所有权限）

⚠️ **安全提醒**：首次登录后请立即修改密码！

## 📝 启动流程详解

### 启动步骤顺序

```
1. 加载配置文件
   ↓
2. 初始化日志系统
   ↓
3. 初始化数据库连接
   ↓
4. 初始化Redis连接
   ↓
5. 初始化Casbin权限控制
   ↓
6. 数据库表结构自动迁移
   ↓
7. 初始化仓库（Repository）
   ↓
8. 初始化应用服务（Service）
   ↓
9. 初始化默认数据（角色、权限、用户、短信模板）
   ↓
10. 注册HTTP路由
    ↓
11. 启动HTTP服务器
```

### Bootstrap 启动器

项目使用 `Bootstrap` 启动器封装所有初始化逻辑：

```go
bootstrap := NewBootstrap()

// 按顺序初始化各组件
bootstrap.InitConfig()        // 配置
bootstrap.InitLogger()         // 日志
bootstrap.InitDatabase()       // 数据库
bootstrap.InitRedis()          // Redis
bootstrap.InitCasbin()         // 权限控制
bootstrap.MigrateDatabase()    // 表迁移
repos := bootstrap.InitRepositories()      // 仓库
services := bootstrap.InitServices(repos)   // 服务
bootstrap.InitDefaultData(repos, services)  // 默认数据
```

### 默认数据初始化

启动时会自动创建：

1. **默认角色**：
   - 超级管理员（super_admin）
   - 管理员（admin）
   - 普通用户（user）

2. **默认权限**：
   - 用户列表、用户信息、创建用户、更新用户、删除用户

3. **默认用户**：
   - 手机号：18797131041
   - 密码：hello#123world
   - 角色：超级管理员

4. **默认短信模板**：
   - 注册模板（register）
   - 忘记密码模板（forgot）

## 🔧 常见问题

### Q1: 数据库连接失败

**错误信息：**
```
初始化数据库失败: dial tcp: connect: connection refused
```

**解决方案：**
1. 检查数据库服务是否启动
2. 验证配置文件中的数据库连接信息
3. 检查防火墙设置
4. 验证用户权限

### Q2: Redis 连接失败

**错误信息：**
```
初始化Redis失败: redis连接失败
```

**解决方案：**
1. 检查 Redis 服务是否启动：`redis-cli ping`
2. 验证配置文件中的 Redis 连接信息
3. 检查 Redis 密码是否正确

### Q3: 端口被占用

**错误信息：**
```
listen tcp :8080: bind: address already in use
```

**解决方案：**
1. 更改配置文件中的端口号
2. 或结束占用端口的进程：
   ```bash
   # Linux/macOS
   lsof -ti:8080 | xargs kill -9
   
   # Windows
   netstat -ano | findstr :8080
   taskkill /PID <进程ID> /F
   ```

### Q4: 表迁移失败

**错误信息：**
```
数据库迁移失败: Error 1071: Specified key was too long
```

**解决方案：**
1. 检查数据库字符集是否为 utf8mb4
2. 如果是 MySQL，检查 innodb_large_prefix 是否启用
3. 手动执行 SQL 脚本创建表结构

### Q5: 短信功能不可用

**问题：** 发送短信验证码报错

**开发模式（sms.enabled: false）：**
- 正常现象，验证码仅在日志中显示
- 查看日志文件获取验证码：`tail -f logs/app.log`

**生产模式（sms.enabled: true）：**
- 需要在数据库 AA02 表中配置腾讯云短信参数
- 参考 [短信配置文档](./SMS.md)

## 🔍 调试技巧

### 1. 开启详细日志

在配置文件中设置：

```yaml
log:
  level: debug  # 改为 debug 级别
```

### 2. 查看启动日志

```bash
# 实时查看日志
tail -f logs/app.log

# 搜索错误
grep -i error logs/app.log

# 查看最近的日志
tail -n 100 logs/app.log
```

### 3. 使用 Go 调试器

```bash
# 使用 delve
dlv debug cmd/server/main.go

# 设置断点
(dlv) b main.go:50
```

### 4. 检查配置加载

启动时会在日志中显示配置信息（debug 模式）：

```json
{
  "level": "info",
  "msg": "配置加载成功",
  "config": {
    "app": {...},
    "database": {...}
  }
}
```

## 📚 下一步

- 了解 [Casbin 权限控制](./CASBIN.md)
- 了解 [API 接口文档](./API.md)
- 了解 [短信服务配置](./SMS.md)
- 查看 [完整项目文档](../README.md)

