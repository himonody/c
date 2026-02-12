# Habit Admin 开发指南

## 🚀 快速开始

### 环境准备

1. **Node.js**: >= 16.0.0
2. **pnpm**: >= 8.0.0 (推荐使用 pnpm)
3. **Git**: 最新版本

### 安装步骤

```bash
# 克隆项目
git clone <repository-url>
cd habit/admin

# 使用脚本快速启动
./scripts/dev.sh

# 或手动启动
pnpm install
cp .env.example .env
pnpm dev
```

## 📁 项目结构详解

```
src/
├── api/                    # API 接口层
│   ├── auth.ts            # 认证相关接口
│   ├── challenge.ts       # 挑战管理接口
│   ├── config.ts          # 系统配置接口
│   └── user.ts            # 用户管理接口
├── assets/                # 静态资源
├── components/            # 公共组件
│   ├── Breadcrumb.tsx     # 面包屑导航
│   ├── ErrorBoundary.tsx  # 错误边界
│   ├── Header.tsx         # 页面头部
│   ├── LazyLoad.tsx       # 懒加载组件
│   ├── Loading.tsx        # 加载组件
│   ├── PageHeader.tsx     # 页面头部
│   ├── ProtectedRoute.tsx # 路由保护
│   └── Sidebar.tsx        # 侧边栏
├── config/                # 配置文件
│   └── index.ts           # 应用配置
├── enums/                 # 枚举定义
│   └── index.ts           # 状态、角色等枚举
├── hooks/                 # 自定义 Hooks
│   └── index.ts           # 通用 Hooks
├── language/              # 国际化文件
│   ├── zh-CN.json         # 中文语言包
│   └── en-US.json         # 英文语言包
├── layouts/               # 布局组件
│   ├── AppLayout.tsx      # 主应用布局
│   └── AuthLayout.tsx     # 认证页面布局
├── redux/                 # 状态管理
│   ├── modules/           # 状态模块
│   │   ├── app.ts         # 应用状态
│   │   ├── auth.ts        # 认证状态
│   │   └── user.ts        # 用户状态
│   └── index.ts           # Store 配置
├── routers/               # 路由配置
│   └── index.tsx          # 路由定义
├── styles/                # 样式文件
│   └── index.less         # 全局样式
├── test/                  # 测试文件
│   ├── setup.ts           # 测试配置
│   └── login.test.tsx     # 登录组件测试
├── typings/               # TypeScript 类型
│   └── index.ts           # 通用类型定义
├── utils/                 # 工具函数
│   ├── constants.ts       # 常量定义
│   ├── index.ts           # 通用工具
│   ├── permission.ts      # 权限工具
│   ├── performance.ts     # 性能监控
│   └── request.ts         # HTTP 请求封装
└── views/                 # 页面组件
    ├── dashboard/         # 仪表盘
    ├── login/             # 登录页
    ├── challenge/         # 挑战管理
    ├── config/            # 系统配置
    ├── user/              # 用户管理
    ├── settings/          # 系统设置
    └── 404/               # 404页面
```

## 🛠️ 开发规范

### 代码风格

项目使用 ESLint + Prettier + Stylelint 确保代码质量：

```bash
# 代码检查
pnpm lint

# 代码格式化
pnpm format

# 类型检查
pnpm type-check
```

### 组件开发规范

1. **函数式组件**: 使用 React.FC 类型
2. **Hooks 优先**: 优先使用 Hooks 而非 Class 组件
3. **TypeScript**: 所有组件必须有完整的类型定义
4. **Props 接口**: 组件 Props 必须定义接口

```tsx
// ✅ 正确示例
interface UserCardProps {
  user: UserInfo
  onUpdate?: (user: UserInfo) => void
}

const UserCard: React.FC<UserCardProps> = ({ user, onUpdate }) => {
  return <div>{user.name}</div>
}

export default UserCard
```

### API 接口规范

1. **统一使用 Axios**: 所有 API 调用通过 request 工具
2. **类型定义**: 每个 API 必须定义请求/响应类型
3. **错误处理**: 统一的错误处理机制
4. **接口分组**: 按功能模块分组 API

```tsx
// ✅ 正确示例
export interface CreateConfigRequest {
  configName: string
  configKey: string
  configValue?: string
}

export const createConfig = (data: CreateConfigRequest) => {
  return request.post<ConfigInfo>('/admin/config/create', data)
}
```

### 状态管理规范

1. **Redux Toolkit**: 使用 Redux Toolkit 简化状态管理
2. **模块化**: 按功能模块划分 state
3. **类型安全**: 所有 state 必须有 TypeScript 类型
4. **异步处理**: 使用 createAsyncThunk 处理异步操作

```tsx
// ✅ 正确示例
export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
    },
  },
})
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
pnpm test

# 运行测试并生成覆盖率报告
pnpm test:coverage

# 监听模式运行测试
pnpm test:watch
```

### 测试规范

1. **组件测试**: 使用 @testing-library/react
2. **Mock 策略**: 合理 Mock 外部依赖
3. **覆盖率**: 保持 80% 以上的测试覆盖率
4. **测试命名**: 使用描述性的测试名称

```tsx
// ✅ 正确示例
describe('Login Component', () => {
  it('renders login form correctly', () => {
    render(<Login />)
    expect(screen.getByText('登录')).toBeInTheDocument()
  })
})
```

## 📦 构建部署

### 构建命令

```bash
# 开发环境构建
pnpm build:dev

# 测试环境构建
pnpm build:test

# 生产环境构建
pnpm build:prod

# 使用构建脚本
./scripts/build.sh
```

### Docker 部署

```bash
# 构建镜像
docker build -t habit-admin .

# 运行容器
docker run -p 3000:80 habit-admin

# 使用 docker-compose
docker-compose up -d
```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| VITE_API_BASE_URL | API 基础地址 | http://localhost:8080 |
| VITE_APP_TITLE | 应用标题 | Habit Admin |
| VITE_APP_ENV | 应用环境 | development |

## 🔧 调试技巧

### 1. Redux DevTools

安装 Redux DevTools 浏览器扩展，可以方便地调试状态变化。

### 2. 性能监控

项目内置了性能监控，可以在控制台查看性能指标：

```javascript
// 查看性能指标
console.log(performance.getEntriesByType('navigation'))
```

### 3. 错误追踪

所有错误都会被 ErrorBoundary 捕获，并在控制台显示详细信息。

## 🐛 常见问题

### Q: 页面白屏
A: 检查浏览器控制台是否有错误，可能是路由配置或组件导入问题。

### Q: API 请求失败
A: 检查 `.env` 文件中的 API 地址配置，确保后端服务正常运行。

### Q: 样式不生效
A: 确保 Less 文件正确导入，检查 CSS Modules 配置。

### Q: 构建失败
A: 清除 node_modules 重新安装依赖，检查 Node.js 版本。

## 📚 学习资源

- [React 官方文档](https://react.dev/)
- [TypeScript 手册](https://www.typescriptlang.org/docs/)
- [Ant Design 组件库](https://ant.design/)
- [Redux Toolkit 文档](https://redux-toolkit.js.org/)
- [Vite 构建工具](https://vitejs.dev/)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码变更
4. 编写测试用例
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。
