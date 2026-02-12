# Habit Admin Frontend

基于 React + TypeScript + Ant Design 的习惯打卡管理后台系统。

## 功能特性

- 🚀 基于 Vite + React 18 + TypeScript 构建
- 🎨 使用 Ant Design 5.x 组件库
- 📦 集成 Redux Toolkit 状态管理
- 🌐 支持国际化（i18n）
- 📱 响应式设计，支持移动端
- 🔐 完整的权限管理系统
- 📊 数据可视化图表
- 🛠️ 完整的代码规范和工具链

## 技术栈

- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **UI 库**: Ant Design 5.x
- **状态管理**: Redux Toolkit + Redux Persist
- **路由**: React Router v6
- **HTTP 客户端**: Axios
- **样式**: Less
- **图标**: Ant Design Icons
- **代码规范**: ESLint + Prettier + Stylelint

## 项目结构

```
src/
├── api/              # API 接口
│   ├── auth.ts       # 认证相关接口
│   ├── challenge.ts  # 挑战管理接口
│   ├── config.ts     # 配置管理接口
│   └── user.ts       # 用户管理接口
├── assets/           # 静态资源
├── components/       # 公共组件
│   ├── Header.tsx    # 头部组件
│   └── Sidebar.tsx   # 侧边栏组件
├── config/           # 配置文件
├── enums/            # 枚举定义
├── hooks/            # 自定义 Hooks
├── language/         # 国际化文件
├── layouts/          # 布局组件
│   ├── AppLayout.tsx    # 主布局
│   └── AuthLayout.tsx   # 认证布局
├── redux/            # Redux 状态管理
│   ├── modules/         # 状态模块
│   └── index.ts         # Store 配置
├── routers/          # 路由配置
├── styles/           # 全局样式
├── typings/          # TypeScript 类型定义
├── utils/            # 工具函数
└── views/            # 页面组件
    ├── dashboard/    # 仪表盘
    ├── login/        # 登录页
    ├── challenge/    # 挑战管理
    ├── config/       # 系统配置
    ├── user/         # 用户管理
    ├── settings/     # 系统设置
    └── 404/          # 404页面
```

## 快速开始

### 环境要求

- Node.js >= 16
- pnpm >= 8

### 安装依赖

```bash
pnpm install
```

### 环境配置

复制环境变量配置文件：

```bash
cp .env.example .env
```

根据实际情况修改 `.env` 文件中的配置。

### 开发环境

```bash
pnpm dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
# 开发环境构建
pnpm build:dev

# 测试环境构建
pnpm build:test

# 生产环境构建
pnpm build:prod
```

### 预览构建结果

```bash
pnpm preview
```

## 主要功能模块

### 1. 认证管理
- 管理员登录/登出
- JWT Token 认证
- 权限验证

### 2. 挑战管理
- 挑战列表查看
- 新增/编辑挑战配置
- 挑战参数设置（结算时间、奖池规则等）
- 挑战状态管理

### 3. 系统配置
- 系统参数配置
- 配置项增删改查
- 配置类型管理（字符串、数字、布尔、JSON）
- 前端/后端配置分离

### 4. 用户管理
- 用户列表查看
- 用户状态管理
- 用户角色分配
- 密码重置功能

### 5. 数据统计
- 实时数据仪表盘
- 挑战参与统计
- 收益分析图表
- 系统运行状态

## API 接口

前端通过 `/api` 前缀调用后端接口，具体接口文档请参考后端项目。

### 主要接口

#### 认证接口
- `POST /api/admin/auth/login` - 管理员登录
- `POST /api/admin/auth/logout` - 管理员登出
- `POST /api/admin/auth/me` - 获取当前管理员信息

#### 挑战管理接口
- `POST /api/admin/challenge/list` - 获取挑战列表
- `POST /api/admin/challenge/create` - 创建挑战
- `POST /api/admin/challenge/update` - 更新挑战
- `POST /api/admin/challenge/delete` - 删除挑战

#### 配置管理接口
- `POST /api/admin/config/list` - 获取配置列表
- `POST /api/admin/config/create` - 创建配置
- `POST /api/admin/config/update` - 更新配置
- `POST /api/admin/config/delete` - 删除配置

#### 用户管理接口
- `POST /api/admin/user/list` - 获取用户列表
- `POST /api/admin/user/create` - 创建用户
- `POST /api/admin/user/update` - 更新用户
- `POST /api/admin/user/delete` - 删除用户

## 开发规范

### 代码规范
- 使用 ESLint + Prettier 进行代码格式化
- 遵循 TypeScript 严格模式
- 组件使用函数式组件 + Hooks
- 使用语义化的变量和函数命名

### 文件命名规范
- 组件文件使用 PascalCase：`UserProfile.tsx`
- 工具文件使用 camelCase：`formatDate.ts`
- 常量文件使用 UPPER_CASE：`API_CONSTANTS.ts`

### 提交规范
- feat: 新功能
- fix: 修复问题
- docs: 文档更新
- style: 代码格式调整
- refactor: 代码重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

### Git 提交格式
```
<type>(<scope>): <subject>

<body>

<footer>
```

示例：
```
feat(auth): add user login functionality

- Add login form component
- Implement JWT token handling
- Add authentication middleware

Closes #123
```

## 部署

### 环境变量

创建 `.env.production` 文件：

```env
VITE_API_BASE_URL=https://your-api-domain.com
VITE_APP_TITLE=Habit Admin
VITE_APP_ENV=production
```

### Docker 部署

```bash
# 构建镜像
docker build -t habit-admin .

# 运行容器
docker run -p 3000:80 habit-admin
```

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name admin.habit.com;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 常见问题

### 1. 开发环境接口请求失败
检查 `vite.config.ts` 中的代理配置是否正确，确保后端服务已启动。

### 2. 构建失败
检查 Node.js 版本是否符合要求，清除 `node_modules` 重新安装依赖。

### 3. 样式不生效
确保 Less 文件已正确导入，检查 `vite.config.ts` 中的 Less 配置。

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证。
