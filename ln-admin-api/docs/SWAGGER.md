# Swagger API 文档使用指南

本文档介绍如何使用和查看项目的 Swagger API 文档。

## 📖 什么是 Swagger

Swagger 是一个用于设计、构建、记录和使用 RESTful Web 服务的开放标准规范。本项目使用 `swaggo/swag` 来自动生成 API 文档。

**主要优势：**
- ✅ 自动生成美观的 API 文档
- ✅ 在线测试接口，无需 Postman
- ✅ 文档与代码同步，减少维护成本
- ✅ 团队协作时统一接口规范

## 🚀 快速开始

### 1. 安装 Swagger 工具

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 确保 $GOPATH/bin 或 $HOME/go/bin 在 PATH 环境变量中
# Windows: C:\Users\{username}\go\bin
# Linux/macOS: ~/go/bin
```

**验证安装：**
```bash
swag --version
```

### 2. 安装依赖包

项目已包含以下依赖（已在 `go.mod` 中）：

- `github.com/swaggo/swag` - Swagger 文档生成工具
- `github.com/swaggo/gin-swagger` - Gin 框架集成
- `github.com/swaggo/files` - Swagger UI 静态文件

如果缺失，请执行：

```bash
cd ln-admin-api
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go mod tidy
```

### 3. 生成 Swagger 文档

在项目根目录执行：

```bash
cd ln-admin-api
swag init -g cmd/server/main.go -o ./docs/swagger
```

**参数说明：**
- `-g`: 指定主文件路径（包含 `@title` 等注释的文件）
- `-o`: 指定输出目录

**成功后会输出：**
```
2024/xx/xx xx:xx:xx Generate swagger docs....
2024/xx/xx xx:xx:xx Generate general API Info
2024/xx/xx xx:xx:xx create docs.go at docs/swagger/docs.go
2024/xx/xx xx:xx:xx create swagger.json at docs/swagger/swagger.json
2024/xx/xx xx:xx:xx create swagger.yaml at docs/swagger/swagger.yaml
```

### 4. 启用文档导入（重要）

生成文档后，需要取消注释路由文件中的导入：

```go
// internal/interfaces/http/router/router.go
import (
    // ...
    _ "github.com/Light-ink-yht/ln-admin/docs/swagger" // 取消这行的注释
)
```

### 5. 查看文档

启动服务：

```bash
go run cmd/server/main.go
```

访问 Swagger UI：

```
http://localhost:8080/swagger/index.html
```

## 📝 当前接口列表

项目已为以下接口添加了 Swagger 注释：

### 认证相关（无需 Token）

| 接口 | 方法 | 说明 | 文档状态 |
|------|------|------|---------|
| `/api/user/captcha` | GET | 获取图片验证码 | ✅ |
| `/api/user/login` | POST | 用户登录 | ✅ |
| `/api/user/signup` | POST | 用户注册 | ✅ |
| `/api/user/password` | POST | 忘记密码 | ✅ |
| `/api/user/signup/code` | POST | 发送短信验证码 | ✅ |
| `/api/user/refresh-token` | POST | 刷新Token | ✅ |

### 用户管理（需要 Token）

| 接口 | 方法 | 说明 | 权限 | 文档状态 |
|------|------|------|------|---------|
| `/api/user/userinfo` | GET | 获取当前用户信息 | user:info | ✅ |
| `/api/user/list` | GET | 获取用户列表 | user:list | ✅ |

## 🔍 查看文档的方式

### 方法1：Swagger UI（推荐）

1. **启动服务**
   ```bash
   go run cmd/server/main.go
   ```

2. **浏览器访问**
   ```
   http://localhost:8080/swagger/index.html
   ```

3. **功能特性**
   - 📖 查看所有 API 接口
   - 🧪 在线测试接口（Try it out）
   - 📋 查看请求/响应示例
   - 🔐 在界面中输入 Token 进行认证
   - 📥 下载 OpenAPI 规范文件

### 方法2：ReDoc（可选）

如需使用 ReDoc，可在路由中添加：

```go
// internal/interfaces/http/router/router.go
r.StaticFile("/redoc", "./docs/swagger/swagger.yaml")
```

访问：
```
http://localhost:8080/redoc
```

### 方法3：直接查看文件

查看生成的 JSON 或 YAML 文件：

```bash
# 查看 JSON
cat docs/swagger/swagger.json

# 查看 YAML
cat docs/swagger/swagger.yaml
```

### 方法4：导入 Postman/Apifox

可以使用生成的 `swagger.json` 导入到 Postman 或 Apifox 中。

## 📚 Swagger 注释语法

### 主函数注释（main.go）

已在 `cmd/server/main.go` 中添加：

```go
// @title           LN Admin API
// @version         1.0
// @description     基于 Go + Gin + GORM 构建的现代化后台管理系统后端服务，采用 DDD（领域驱动设计）架构和 RBAC 权限控制
// @termsOfService  https://github.com/Light-ink-yht/ln-admin

// @contact.name   API Support
// @contact.url    https://github.com/Light-ink-yht/ln-admin
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT认证，格式：Bearer {token}。登录后获取 accessToken，在请求头中添加：Authorization: Bearer {accessToken}
```

### 接口注释（Handler）

所有接口已添加完整的 Swagger 注释，示例：

```go
// Login 用户登录
// @Summary      用户登录
// @Description  使用手机号和密码登录，需要先获取图片验证码
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "登录请求"
// @Success      200      {object}  response.Response{data=dto.LoginResponse}  "登录成功，返回双Token和用户信息"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "认证失败"
// @Failure      403      {object}  response.Response  "用户已禁用"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    // ...
}
```

**注释说明：**
- `@Summary`: 接口简短描述（必填）
- `@Description`: 接口详细描述（可选）
- `@Tags`: 接口分组标签
- `@Accept`: 接受的请求类型（json, form, multipart/form-data）
- `@Produce`: 返回的内容类型（json, xml）
- `@Param`: 请求参数（query/form/body/path）
- `@Success`: 成功响应示例
- `@Failure`: 失败响应示例
- `@Security`: 需要认证时使用 `@Security Bearer`
- `@Router`: 路由路径和HTTP方法

## 🎯 常用 Tag 分类

已使用的标签分类：

- `认证相关` - 登录、注册、忘记密码、Token刷新、验证码
- `用户管理` - 用户信息的查询和管理

未来可添加：
- `权限管理` - 角色、权限的管理
- `系统配置` - 系统配置的管理

## 🔐 Token 认证使用

### 在 Swagger UI 中测试

1. **获取 Token**
   - 先调用登录接口：`POST /api/user/login`
   - 复制返回的 `accessToken`

2. **设置 Token**
   - 点击 Swagger UI 右上角的 **`Authorize`** 按钮
   - 输入 Token（格式：`Bearer {token}`）或直接输入 token（会自动添加 Bearer 前缀）
   - 点击 **`Authorize`** 保存
   - 点击 **`Close`** 关闭对话框

3. **测试接口**
   - 现在可以测试所有需要认证的接口了
   - Token 会自动添加到请求头：`Authorization: Bearer {token}`

### 在代码中测试

```bash
# 1. 登录获取Token
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "18797131041",
    "password": "hello#123world",
    "captchaCode": "1234",
    "captchaId": "xxxxx"
  }'

# 响应示例：
# {
#   "code": 200,
#   "msg": "success",
#   "data": {
#     "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#     "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#     ...
#   }
# }

# 2. 使用Token访问接口
curl -X GET http://localhost:8080/api/user/userinfo \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📋 完整接口文档

### 认证相关

#### 1. 获取图片验证码

**接口：** `GET /api/user/captcha`

**说明：** 获取用于登录/注册的图片验证码

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "captchaId": "xxxx-xxxx-xxxx",
    "captcha_code": "data:image/png;base64,iVBORw0KG..."
  }
}
```

#### 2. 用户登录

**接口：** `POST /api/user/login`

**请求体：**
```json
{
  "phone": "18797131041",
  "password": "hello#123world",
  "captchaCode": "1234",
  "captchaId": "xxxx-xxxx-xxxx"
}
```

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 7200,
    "refreshExpiresIn": 604800,
    "user": {
      "userId": "uuid",
      "phone": "18797131041",
      "nickname": "超级管理员",
      ...
    }
  }
}
```

#### 3. 用户注册

**接口：** `POST /api/user/signup`

**请求体：**
```json
{
  "phone": "13800138000",
  "code": "123456",
  "password": "password123",
  "confirmPassword": "password123"
}
```

#### 4. 发送短信验证码

**接口：** `POST /api/user/signup/code`

**请求体：**
```json
{
  "phone": "13800138000",
  "type": "register"  // 或 "forgot"
}
```

#### 5. 忘记密码

**接口：** `POST /api/user/password`

**请求体：**
```json
{
  "phone": "13800138000",
  "code": "123456",
  "password": "newpassword123",
  "confirmPassword": "newpassword123"
}
```

#### 6. 刷新Token

**接口：** `POST /api/user/refresh-token`

**请求体：**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 用户管理

#### 1. 获取当前用户信息

**接口：** `GET /api/user/userinfo`

**认证：** 需要 Token

**权限：** `user:info`

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "userId": "uuid",
    "phone": "18797131041",
    "nickname": "超级管理员",
    "fullName": "系统管理员",
    "status": "1",
    ...
  }
}
```

#### 2. 获取用户列表

**接口：** `GET /api/user/list`

**认证：** 需要 Token

**权限：** `user:list`

**查询参数：**
- `page` (int, 可选): 页码，从1开始，默认1
- `page_size` (int, 可选): 每页数量，最大100，默认10
- `phone` (string, 可选): 手机号（模糊查询）
- `email` (string, 可选): 邮箱（模糊查询）
- `status` (string, 可选): 状态（1启用，2禁用）
- `gender` (string, 可选): 性别（1男，2女，3未知）
- `nickname` (string, 可选): 昵称（模糊查询）
- `full_name` (string, 可选): 姓名（模糊查询）

**请求示例：**
```
GET /api/user/list?page=1&page_size=20&status=1
```

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "userId": "uuid1",
      "phone": "18797131041",
      ...
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

## 🛠️ 问题排查

### Q1: swag 命令未找到

**错误：** `swag: command not found`

**解决：**
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 检查 PATH
# Windows PowerShell
$env:PATH -split ';' | Select-String go

# Linux/macOS
echo $PATH | grep go
```

**添加到 PATH（如未自动添加）：**
```bash
# Windows（PowerShell，临时）
$env:PATH += ";$env:USERPROFILE\go\bin"

# Windows（永久，系统环境变量）
# 添加到 PATH: C:\Users\{username}\go\bin

# Linux/macOS
export PATH=$PATH:~/go/bin
# 永久添加：写入 ~/.bashrc 或 ~/.zshrc
```

### Q2: 文档未更新

**问题：** 修改注释后文档未更新

**解决：**
1. 重新生成文档：
   ```bash
   swag init -g cmd/server/main.go -o ./docs/swagger
   ```
2. 重启服务
3. 清除浏览器缓存后刷新页面（Ctrl+F5）

### Q3: Swagger UI 无法访问

**错误：** 404 Not Found

**检查：**
1. 确认路由已注册：`r.GET("/swagger/*any", ginSwagger.WrapHandler(...))`
2. 确认服务已启动
3. 确认端口号正确（默认 8080）
4. 确认已生成文档文件（`docs/swagger/docs.go` 存在）

### Q4: 注释解析错误

**错误：** `parse swagger comment error`

**解决：**
1. 检查注释格式是否正确
2. 确保 `@Summary`、`@Router` 等标签拼写正确
3. 查看 swag 输出的详细错误信息
4. 常见错误：
   - 缺少 `@Router` 标签
   - `@Param` 参数类型错误
   - 响应类型格式错误

### Q5: 导入 docs/swagger 报错

**错误：** `could not import github.com/Light-ink-yht/ln-admin/docs/swagger`

**解决：**
1. 先执行 `swag init` 生成文档
2. 取消注释路由文件中的导入：
   ```go
   _ "github.com/Light-ink-yht/ln-admin/docs/swagger"
   ```
3. 如果仍有问题，检查 `docs/swagger/docs.go` 文件是否存在

## 📖 更多资源

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 规范](https://swagger.io/specification/)
- [Gin-Swagger 文档](https://github.com/swaggo/gin-swagger)
- [Swagger 注释示例](https://github.com/swaggo/swag/blob/master/example/celler/controller/product.go)

## 🔄 工作流程

### 开发新接口时

1. **编写 Handler 函数**

2. **添加 Swagger 注释**
   ```go
   // YourHandler 你的处理函数
   // @Summary      接口名称
   // @Description  接口描述
   // @Tags         标签
   // @Accept       json
   // @Produce      json
   // @Param        request  body  YourRequest  true  "请求参数"
   // @Success      200      {object}  response.Response{data=YourResponse}
   // @Router       /your/path [post]
   ```

3. **生成文档**
   ```bash
   swag init -g cmd/server/main.go -o ./docs/swagger
   ```

4. **测试接口**
   - 启动服务
   - 在 Swagger UI 中测试
   - 验证请求/响应格式

5. **提交代码**
   - 提交代码和生成的 `docs/swagger` 文件

## 💡 最佳实践

1. **及时更新文档**
   - 修改接口后立即更新注释
   - 提交前生成最新文档

2. **详细的注释**
   - `@Summary` 简洁明了（一行描述）
   - `@Description` 详细说明业务逻辑、注意事项
   - 参数说明清晰，包含类型、是否必填、默认值

3. **示例数据**
   - 提供真实的请求/响应示例
   - 说明每个字段的含义和取值范围

4. **错误码说明**
   - 明确各种错误场景
   - 说明错误原因和解决方案

5. **版本管理**
   - 及时更新 `@version`
   - 重大变更记录在 `@description` 中

## 📝 示例：完整接口注释

```go
// Login 用户登录
// @Summary      用户登录
// @Description  使用手机号和密码登录，需要先获取图片验证码。
// @Description  
// @Description  **注意事项：**
// @Description  - AccessToken有效期为2小时，过期后使用RefreshToken刷新
// @Description  - RefreshToken有效期为7天
// @Description  - 验证码有效期为5分钟，仅可使用一次
// @Description  - 连续5次登录失败将锁定账户30分钟
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "登录请求"
// @Success      200      {object}  response.Response{data=dto.LoginResponse}  "登录成功"
// @Failure      400      {object}  response.Response  "请求参数错误或验证码错误"
// @Failure      401      {object}  response.Response  "用户名或密码错误"
// @Failure      403      {object}  response.Response  "用户已禁用"
// @Failure      429      {object}  response.Response  "登录次数过多，账户已锁定"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    // 实现代码...
}
```

## ✅ 检查清单

在提交代码前，确保：

- [ ] 所有接口都添加了 Swagger 注释
- [ ] 已执行 `swag init` 生成最新文档
- [ ] Swagger UI 可以正常访问
- [ ] 所有接口可以在 Swagger UI 中测试
- [ ] Token 认证功能正常
- [ ] 响应示例正确
- [ ] 错误码说明完整

## 🎉 总结

使用 Swagger 可以：

✅ 自动生成美观的 API 文档  
✅ 在线测试接口，无需 Postman  
✅ 团队协作时统一接口规范  
✅ 前后端联调时减少沟通成本  
✅ 文档与代码同步更新，减少维护成本  

开始使用 Swagger，让 API 文档更清晰、测试更便捷！

---

**快速命令总结：**

```bash
# 生成文档
swag init -g cmd/server/main.go -o ./docs/swagger

# 启动服务
go run cmd/server/main.go

# 访问文档
# http://localhost:8080/swagger/index.html
```
