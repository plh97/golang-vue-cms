
# GoLang Vue CMS Project | GoLang Vue 内容管理系统

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

A full-stack Content Management System (CMS) built with Go (Gin, GORM, Casbin) for the backend and Vue 3 (Vite, TailwindCSS) for the frontend. Features robust RBAC (role-based access control), JWT authentication, and a modern, responsive UI. Project structure is inspired by [nunu-go](https://github.com/go-nunu/nunu).

### Features

- 🔐 **RBAC Authorization**: Role-based access control with Casbin
- 👤 **User Management**: Full CRUD with role assignment
- 🔑 **JWT Authentication**: Secure token-based login
- 📊 **RESTful API**: Clean API with Gin
- 🎨 **Modern UI**: Vue3 + TailwindCSS
- 🗄️ **Database**: MySQL + GORM

### Tech Stack

**Backend**
- Gin (HTTP framework)
- GORM (ORM)
- Casbin (RBAC)
- JWT
- Wire (DI)

**Frontend**
- Vue 3
- Vite
- TailwindCSS
- Composition API
- Fetch API

### Project Structure

```
.
├── api/         # API DTOs
├── cmd/         # Entrypoints (server, migration, task)
├── config/      # Config files (local.yml, prod.yml, model.conf)
├── internal/    # Main backend logic
│   ├── handler/     # HTTP handlers
│   ├── middleware/  # JWT, RBAC, CORS
│   ├── model/       # GORM models
│   ├── repository/  # Data access
│   ├── router/      # Routing
│   ├── service/     # Business logic
│   └── server/      # Server startup
├── pkg/         # Utilities (casbin, jwt, log, etc)
└── web/         # Frontend (Vue3 app)
    └── src/
        ├── api.ts
        ├── components/
        └── pages/
```

### Getting Started

#### Prerequisites
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- pnpm

#### Backend Setup
1. Clone repo
   ```bash
   git clone https://github.com/plh97/golang-tutorial.git
   cd golang-tutorial
   ```
2. Install Go dependencies
   ```bash
   go mod download
   ```
3. Configure database in `config/local.yml`
4. Run migrations
   ```bash
   make migration
   # or
   go run cmd/migration/main.go
   ```
5. Start server
   ```bash
   make server
   # or
   go run cmd/server/main.go
   ```
   Server: http://localhost:8291

#### Frontend Setup
1. Enter web directory
   ```bash
   cd web
   ```
2. Install dependencies
   ```bash
   pnpm install
   ```
3. Start dev server
   ```bash
   pnpm dev
   ```
   Frontend: http://localhost:8000

### RBAC Architecture

**Business Layer (MySQL):**
- Tables: users, roles, permissions, user_roles, role_permissions
- Used for: UI, management, assignments

**Access Layer (Casbin):**
- Table: casbin_rule
- Format: `p, role_key, api_path, method`
- Used for: Fast API permission checks

**Flow:**
```
Request → JWT Middleware → RBAC Middleware (Casbin) → Handler
```

**Permission Update:**
```
UI → Update role_permissions (GORM) → Sync casbin_rule (Casbin)
```

### API Endpoints

**Auth**
- POST /v1/login
- GET /v1/profile

**User**
- POST /v1/user/list
- POST /v1/user
- PUT /v1/user
- DELETE /v1/user/:id

**Role**
- GET /v1/role/list
- POST /v1/role
- PUT /v1/role
- DELETE /v1/role/:id

**Permission**
- GET /v1/permission/list

### Default Credentials

- Username: admin@gmail.com
- Password: admin123
- Role: Administrator

### Development

**Wire dependencies**
```bash
cd cmd/server/wire
wire
```

**Run tests**
```bash
go test ./...
```

**Build**
```bash
# Backend
go build -o bin/server cmd/server/main.go
# Frontend
cd web && pnpm build
```

### License

MIT License

---

<a name="中文"></a>
## 中文

基于 Go (Gin, GORM, Casbin) 后端和 Vue3 (Vite, TailwindCSS) 前端的全栈内容管理系统。支持 RBAC 权限、JWT 登录、现代响应式 UI。项目结构参考 [nunu-go](https://github.com/go-nunu/nunu)。

### 功能特性

- 🔐 **RBAC 权限控制**：基于 Casbin 的角色访问控制
- 👤 **用户管理**：完整的用户增删改查及角色分配
- 🔑 **JWT 认证**：安全的令牌登录
- 📊 **RESTful API**：Gin 框架
- 🎨 **现代 UI**：Vue3 + TailwindCSS
- 🗄️ **数据库**：MySQL + GORM

### 技术栈

**后端**
- Gin (HTTP 框架)
- GORM (ORM)
- Casbin (RBAC)
- JWT
- Wire (依赖注入)

**前端**
- Vue 3
- Vite
- TailwindCSS
- Composition API
- Fetch API

### 项目结构

```
.
├── api/         # API 定义和 DTO
├── cmd/         # 启动入口（server, migration, task）
├── config/      # 配置文件（local.yml, prod.yml, model.conf）
├── internal/    # 后端主逻辑
│   ├── handler/     # HTTP 处理器
│   ├── middleware/  # JWT、RBAC、CORS
│   ├── model/       # GORM 模型
│   ├── repository/  # 数据访问
│   ├── router/      # 路由
│   ├── service/     # 业务逻辑
│   └── server/      # 启动
├── pkg/         # 工具包（casbin, jwt, log 等）
└── web/         # 前端（Vue3 应用）
    └── src/
        ├── api.ts
        ├── components/
        └── pages/
```

### 快速开始

#### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- pnpm

#### 后端设置
1. 克隆仓库
   ```bash
   git clone https://github.com/plh97/golang-tutorial.git
   cd golang-tutorial
   ```
2. 安装 Go 依赖
   ```bash
   go mod download
   ```
3. 配置数据库 `config/local.yml`
4. 运行迁移
   ```bash
   make migration
   # 或
   go run cmd/migration/main.go
   ```
5. 启动服务
   ```bash
   make server
   # 或
   go run cmd/server/main.go
   ```
   地址: http://localhost:8291

#### 前端设置
1. 进入 web 目录
   ```bash
   cd web
   ```
2. 安装依赖
   ```bash
   pnpm install
   ```
3. 启动开发服务器
   ```bash
   pnpm dev
   ```
   地址: http://localhost:8000

### RBAC 架构

**业务层（MySQL）**
- 表：users, roles, permissions, user_roles, role_permissions
- 用于：UI、管理、分配

**访问层（Casbin）**
- 表：casbin_rule
- 格式：`p, role_key, api_path, method`
- 用于：API 权限快速校验

**流程：**
```
请求 → JWT 中间件 → RBAC 中间件（Casbin）→ 处理器
```

**权限更新：**
```
UI → 更新 role_permissions（GORM）→ 同步 casbin_rule（Casbin）
```

### API 接口

**认证**
- POST /v1/login
- GET /v1/profile

**用户**
- POST /v1/user/list
- POST /v1/user
- PUT /v1/user
- DELETE /v1/user/:id

**角色**
- GET /v1/role/list
- POST /v1/role
- PUT /v1/role
- DELETE /v1/role/:id

**权限**
- GET /v1/permission/list

### 默认账号

- 用户名：admin@gmail.com
- 密码：admin123
- 角色：管理员

### 开发

**生成 Wire 依赖**
```bash
cd cmd/server/wire
wire
```

**运行测试**
```bash
go test ./...
```

**构建**
```bash
# 后端
go build -o bin/server cmd/server/main.go
# 前端
cd web && pnpm build
```

### 许可证

MIT License
   ```
