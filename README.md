Hello from Cursor cloud

# Go Arco Admin

通用后台管理系统：Go（Gin + GORM + SQLite + JWT）后端 + 官方 [Arco Design Pro Vue](https://github.com/arco-design/arco-design-pro-vue) 前端。

## 功能

- 登录鉴权：JWT，bcrypt 密码，账号启用/禁用
- 系统管理（仅 `admin` 角色可见）
  - 系统概览：统计卡片、近 7 日操作趋势、最新动态
  - 用户管理：搜索、分页、新增、编辑、重置密码、删除
  - 角色管理：新增、编辑、删除；角色标识直接对应前端路由 `meta.roles`
  - 操作日志：自动记录登录和所有写操作
- 保留 Arco Pro 全部示例页面（仪表盘、列表、表单、详情、结果、异常、可视化、个人中心），这些页面仍使用 mock 数据

## 目录

```
server/   Go 后端（接口 /api/*，同时托管 web/dist）
web/      Arco Design Pro Vue 前端
```

## 本地开发

```bash
# 后端，默认 :8080，首次启动自动建库并写入种子数据
cd server && go run .

# 前端，默认 :3000，/api/user 与 /api/system 代理到 :8080
cd web && pnpm install && pnpm dev
```

默认账号：`admin / admin`（超级管理员），`user / user`（普通用户，看不到系统管理菜单）。

## 生产部署

```bash
cd web && pnpm build
cd ../server && go build -o bin/server . && JWT_SECRET=<随机字符串> ./bin/server
```

访问 `http://localhost:8080` 即可，Go 服务同时提供接口和前端页面。

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `ADDR` | `:8080` | 监听地址 |
| `DB_PATH` | `data/admin.db` | SQLite 文件路径 |
| `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥，生产环境务必修改 |
| `WEB_DIR` | `../web/dist` | 前端构建产物目录 |

## 接口约定

响应沿用 Arco Pro 的格式 `{ code, msg, data }`：`20000` 成功，`50000` 业务错误，`50008`/`50014` 登录失效（前端会提示重新登录），`40300` 无权限。
