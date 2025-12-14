# GoLang Vue CMS Project | GoLang Vue 内容管理系统

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

A full-stack Content Management System built with Go backend and Vue3 frontend, featuring role-based access control (RBAC) using Casbin. Project skeleton comes from the [nunu-go](https://github.com/go-nunu/nunu) scaffold, with RBAC/CRUD flows layered on top.

## Features

- 🔐 **RBAC Authorization**: Role-based access control with Casbin
- 👤 **User Management**: Complete user CRUD operations with role assignment
- 🔑 **JWT Authentication**: Secure token-based authentication
- 📊 **RESTful API**: Clean API design with Gin framework
- 🎨 **Modern UI**: Responsive interface built with Vue3 and TailwindCSS
- 🗄️ **Database**: MySQL with GORM ORM

## Tech Stack

### Backend

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **ORM**: [GORM](https://gorm.io/) - Database ORM
- **Authorization**: [Casbin](https://casbin.org/) - RBAC access control
- **Authentication**: JWT tokens
- **Scaffolding**: [Nunu](https://github.com/go-nunu/nunu) - Go project generator
- **Dependency Injection**: Wire

### Frontend

- **Framework**: [Vue 3](https://vuejs.org/) - Progressive JavaScript framework
- **Build Tool**: [FesJS](https://fesjs.mumblefe.cn/) - Vue3 scaffolding tool
- **Styling**: [TailwindCSS](https://tailwindcss.com/) - Utility-first CSS framework
- **State Management**: Vue Composition API
- **HTTP Client**: Fetch API

## Project Structure

```
.
├── api/                    # API definitions and DTOs
├── cmd/
│   ├── migration/         # Database migration scripts
│   └── server/            # Main server entry point
├── config/                # Configuration files
│   ├── local.yml         # Local development config
│   ├── prod.yml          # Production config
│   └── model.conf        # Casbin RBAC model
├── internal/
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # Custom middleware (JWT, RBAC, CORS)
│   ├── model/            # Database models
│   ├── repository/       # Data access layer
│   ├── router/           # Route definitions
│   ├── service/          # Business logic layer
│   └── server/           # Server initialization
├── pkg/                  # Reusable packages
│   ├── casbin/          # Casbin initialization
│   ├── jwt/             # JWT utilities
│   └── log/             # Logging utilities
└── web/                 # Frontend Vue3 application
    └── src/
        ├── api.ts       # API client
        ├── components/  # Vue components
        └── pages/       # Page components
```

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- pnpm (for frontend)

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/plh97/golang-tutorial.git
   cd golang-tutorial
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure database**
   
   Edit `config/local.yml`:
   ```yaml
   data:
     mysql:
       user: root
       password: your_password
       host: 127.0.0.1
       port: 3306
       dbname: your_database
   ```

4. **Run migrations**
   ```bash
   make migration
   # or
   go run cmd/migration/main.go
   ```
   > This step will automatically create the table structure, admin user, roles, permissions, and other basic data.

5. **Start the server**
   ```bash
   make server
   # or
   go run cmd/server/main.go
   ```

   Server will start at `http://localhost:8291`

### Frontend Setup

1. **Navigate to web directory**
   ```bash
   cd web
   ```

2. **Install dependencies**
   ```bash
   pnpm install
   ```

3. **Start development server**
   ```bash
   pnpm dev
   ```

   Frontend will start at `http://localhost:8000`

## RBAC Architecture

The project implements a comprehensive RBAC system with two-layer data management:

### 1. Business Data Layer (MySQL)
- **Tables**: `users`, `roles`, `permissions`, `user_roles`, `role_permissions`
- **Purpose**: Store user/role/permission metadata and relationships
- **Used for**: Frontend display, permission management UI, role assignment

### 2. Access Control Layer (Casbin)
- **Table**: `casbin_rule`
- **Format**: `p, role_key, api_path, method` (e.g., `p, admin, /v1/user, GET`)
- **Purpose**: Fast permission checking during API requests
- **Used for**: Middleware authorization

### Authorization Flow

```
User Request → JWT Middleware (authenticate) 
            → RBAC Middleware (get user roles → Casbin.Enforce) 
            → Controller (authorized)
```

### Permission Update Flow

```
Update Role Permissions (UI) 
  → Update `role_permissions` table (GORM)
  → Sync to `casbin_rule` table (Casbin API)
```

## API Endpoints

### Authentication
- `POST /v1/login` - User login
- `GET /v1/profile` - Get current user profile

### User Management (Requires Auth)
- `POST /v1/user/list` - List users (paginated)
- `POST /v1/user` - Create user
- `PUT /v1/user` - Update user
- `DELETE /v1/user/:id` - Delete user

### Role Management (Requires Auth)
- `GET /v1/role/list` - List roles
- `POST /v1/role` - Create role
- `PUT /v1/role` - Update role permissions
- `DELETE /v1/role/:id` - Delete role

### Permission Management (Requires Auth)
- `GET /v1/permission/list` - List permissions

## Default Credentials

After running migrations, you can login with:

- **Username**: admin@gmail.com
- **Password**: admin123
- **Role**: Administrator (full access)

## Development

### Generate Wire Dependencies
```bash
cd cmd/server/wire
wire
```

### Run Tests
```bash
go test ./...
```

### Build for Production
```bash
# Backend
go build -o bin/server cmd/server/main.go

# Frontend
cd web
pnpm build
```

## License

MIT License

---

<a name="中文"></a>
## 中文

基于 Go 后端和 Vue3 前端构建的全栈内容管理系统，使用 Casbin 实现基于角色的访问控制（RBAC）。项目骨架来自 [nunu-go](https://github.com/go-nunu/nunu) 脚手架，并在其上扩展了 RBAC 与业务 CRUD 能力。

## 功能特性

- 🔐 **RBAC 权限控制**：基于 Casbin 的角色访问控制
- 👤 **用户管理**：完整的用户增删改查及角色分配
- 🔑 **JWT 认证**：安全的基于令牌的身份验证
- 📊 **RESTful API**：使用 Gin 框架的简洁 API 设计
- 🎨 **现代化 UI**：使用 Vue3 和 TailwindCSS 构建的响应式界面
- 🗄️ **数据库**：MySQL 配合 GORM ORM

## 技术栈

### 后端

- **框架**: [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- **ORM**: [GORM](https://gorm.io/) - 数据库 ORM
- **权限控制**: [Casbin](https://casbin.org/) - RBAC 访问控制
- **身份认证**: JWT 令牌
- **脚手架**: [Nunu](https://github.com/go-nunu/nunu) - Go 项目生成器
- **依赖注入**: Wire

### 前端

- **框架**: [Vue 3](https://vuejs.org/) - 渐进式 JavaScript 框架
- **构建工具**: [FesJS](https://fesjs.mumblefe.cn/) - Vue3 脚手架工具
- **样式**: [TailwindCSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- **状态管理**: Vue Composition API
- **HTTP 客户端**: Fetch API

## 项目结构

```
.
├── api/                    # API 定义和 DTO
├── cmd/
│   ├── migration/         # 数据库迁移脚本
│   └── server/            # 主服务器入口
├── config/                # 配置文件
│   ├── local.yml         # 本地开发配置
│   ├── prod.yml          # 生产环境配置
│   └── model.conf        # Casbin RBAC 模型
├── internal/
│   ├── handler/          # HTTP 请求处理器
│   ├── middleware/       # 自定义中间件（JWT、RBAC、CORS）
│   ├── model/            # 数据库模型
│   ├── repository/       # 数据访问层
│   ├── router/           # 路由定义
│   ├── service/          # 业务逻辑层
│   └── server/           # 服务器初始化
├── pkg/                  # 可复用包
│   ├── casbin/          # Casbin 初始化
│   ├── jwt/             # JWT 工具
│   └── log/             # 日志工具
└── web/                 # 前端 Vue3 应用
    └── src/
        ├── api.ts       # API 客户端
        ├── components/  # Vue 组件
        └── pages/       # 页面组件
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- pnpm（用于前端）

### 后端设置

1. **克隆仓库**
   ```bash
   git clone https://github.com/plh97/golang-tutorial.git
   cd golang-tutorial
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置数据库**
   
   编辑 `config/local.yml`：
   ```yaml
   data:
     mysql:
       user: root
       password: your_password
       host: 127.0.0.1
       port: 3306
       dbname: your_database
   ```

4. **运行迁移**
   ```bash
   make migration
   # 或者
   go run cmd/migration/main.go
   ```
   > 这一步会自动创建表结构、管理员、角色、权限等基础数据。

5. **启动服务器**
   ```bash
   make server
   # 或者
   go run cmd/server/main.go
   ```

   服务器将在 `http://localhost:8291` 启动

### 前端设置

1. **进入 web 目录**
   ```bash
   cd web
   ```

2. **安装依赖**
   ```bash
   pnpm install
   ```

3. **启动开发服务器**
   ```bash
   pnpm dev
   ```

   前端将在 `http://localhost:8000` 启动

## RBAC 架构

项目实现了具有双层数据管理的综合 RBAC 系统：

### 1. 业务数据层（MySQL）
- **数据表**：`users`、`roles`、`permissions`、`user_roles`、`role_permissions`
- **用途**：存储用户/角色/权限元数据和关系
- **应用场景**：前端展示、权限管理界面、角色分配

### 2. 访问控制层（Casbin）
- **数据表**：`casbin_rule`
- **格式**：`p, role_key, api_path, method`（例如：`p, admin, /v1/user, GET`）
- **用途**：API 请求期间快速权限检查
- **应用场景**：中间件授权

### 授权流程

```
用户请求 → JWT 中间件（认证）
        → RBAC 中间件（获取用户角色 → Casbin.Enforce）
        → 控制器（已授权）
```

### 权限更新流程

```
更新角色权限（UI）
  → 更新 `role_permissions` 表（GORM）
  → 同步到 `casbin_rule` 表（Casbin API）
```

## API 接口

### 身份认证
- `POST /v1/login` - 用户登录
- `GET /v1/profile` - 获取当前用户信息

### 用户管理（需要认证）
- `POST /v1/user/list` - 列出用户（分页）
- `POST /v1/user` - 创建用户
- `PUT /v1/user` - 更新用户
- `DELETE /v1/user/:id` - 删除用户

### 角色管理（需要认证）
- `GET /v1/role/list` - 列出角色
- `POST /v1/role` - 创建角色
- `PUT /v1/role` - 更新角色权限
- `DELETE /v1/role/:id` - 删除角色

### 权限管理（需要认证）
- `GET /v1/permission/list` - 列出权限

## 默认凭据

运行迁移后，您可以使用以下凭据登录：

- **用户名**：admin@gmail.com
- **密码**：admin123
- **角色**：管理员（完全访问权限）

## 开发

### 生成 Wire 依赖
```bash
cd cmd/server/wire
wire
```

### 运行测试
```bash
go test ./...
```

### 生产构建
```bash
# 后端
go build -o bin/server cmd/server/main.go

# 前端
cd web
pnpm build
```

## 许可证

MIT License
