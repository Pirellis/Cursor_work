# Go Arco Admin

通用后台管理系统：**Go（Gin + GORM + SQLite）** 后端 + 官方 **[Arco Design Pro Vue](https://github.com/arco-design/arco-design-pro-vue)** 前端。

- 保留 Arco Pro 全部官方示例页（工作台、实时监控、数据可视化、列表、表单、详情、结果、异常、个人中心），亮/暗主题、中英文、多标签页、全局配置开箱即用。
- 登录、用户信息由 Go 后端真实提供（JWT），示例页继续使用 Mock 数据。
- 新增 **系统管理** 模块（仅 `admin` 角色可见/可调用）：
  - 用户管理：数据概览卡片、关键词/角色/状态筛选、分页、新建/编辑/启停/删除
  - 角色管理：卡片式展示、用户数统计，内置角色受保护
  - 操作日志：自动记录所有写操作（含失败的业务码），按方法/结果筛选

## 目录

```
server/   Go 后端（cmd/server 入口，internal/{config,model,middleware,handler,router,resp}）
web/      Arco Design Pro Vue（Vite）前端
```

## 快速开始

依赖：Go 1.22+、Node 18+、pnpm 7（`packageManager` 已锁定为 pnpm@7，以使用官方 lockfile；可用 `corepack enable` 或 `npx pnpm@7`）。

```bash
# 后端（:8080，首次启动自动建库并写入默认账号）
cd server && go run ./cmd/server

# 前端开发（:3000，/api 代理到 :8080）
cd web && npx pnpm@7 install && npx pnpm@7 dev
```

单端口部署：`cd web && pnpm build`，后端会直接托管 `web/dist`（支持 history 路由刷新），访问 http://localhost:8080 即可。

默认账号：

| 账号 | 密码 | 角色 |
| --- | --- | --- |
| admin | admin | 超级管理员，可访问全部菜单与系统管理 |
| user | user | 普通用户，仅基础页面 |

## 配置（环境变量）

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `APP_ADDR` | `:8080` | 监听地址 |
| `APP_DB` | `data/admin.db` | SQLite 文件路径 |
| `APP_JWT_SECRET` | `change-me-in-production` | JWT 签名密钥，**生产必须修改** |
| `APP_STATIC_DIR` | `../web/dist` | 前端构建产物目录 |

前端开发时可用 `API_TARGET=http://host:port pnpm dev` 修改代理目标。

## 约定

- 响应格式与 Arco Pro 拦截器一致：`{ code, msg, data }`，`code = 20000` 为成功；`50008/50014` 触发前端重新登录；`40300` 为无权限。
- 角色编码与路由 `meta.roles` 对应（见 `web/src/router/routes/modules/system.ts`）。后端每次请求都会重新读取用户，禁用账号或修改角色立即生效。
- 新增业务模块：在 `server/internal/handler` 写接口并在 `router.go` 注册；前端在 `web/src/api` 加请求、`router/routes/modules` 加路由、`views/` 加页面即可自动出现在菜单中。

## 测试

```bash
cd server && go test ./...
cd web && npx pnpm@7 type:check && npx pnpm@7 build
```
