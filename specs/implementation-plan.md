# project-alpha 实现计划

## 1. 实施目标

本计划基于 `specs/spec.md`，目标是在第一阶段交付一个可本地运行、前后端联通、数据持久化完整的轻量级 ticket 管理工具。

第一阶段完成后应具备以下能力：

- 使用 PostgreSQL 持久化 ticket、tag 和 ticket-tag 关联数据。
- 后端通过 Go Gin 暴露 `/api/v1` REST API。
- 后端通过 GORM 完成模型映射、查询、关联维护和基础迁移。
- 前端使用 TypeScript、Vite、React、Tailwind CSS、shadcn/ui 构建单页应用。
- 用户可以在 Web 页面完成 ticket 和标签的主要操作。
- 支持 title 搜索、标签筛选、完成状态筛选和分页。

## 2. 实施原则

- 先搭建可运行骨架，再逐步补齐业务功能。
- 先完成后端数据闭环，再接入前端，避免 UI 先行导致接口反复调整。
- API、数据库模型和前端类型保持一致，字段命名以 JSON camelCase、数据库 snake_case 为准。
- 业务规则放在 service 层，handler 只负责 HTTP 入参、出参和状态码。
- 第一阶段优先使用 GORM `AutoMigrate`，不引入显式 migration 工具。
- 前端不引入复杂全局状态库，除非实现中确认手写请求状态明显变复杂。
- 每个阶段结束后都保留可运行状态，避免大量代码一次性集成。

## 3. 推荐目录结构

```text
project-alpha/
  backend/
    cmd/server/main.go
    internal/config/
    internal/database/
    internal/models/
    internal/repositories/
    internal/services/
    internal/handlers/
    internal/router/
    internal/response/
    internal/validation/
  frontend/
    src/
      app/
      components/
        layout/
        tags/
        tickets/
        ui/
      hooks/
      lib/
      types/
      main.tsx
      index.css
  specs/
    spec.md
    implementation-plan.md
  docker-compose.yml
  README.md
```

## 4. 里程碑拆分

### 4.1 M0：项目基础设施

目标：建立可启动的本地开发环境。

任务：

- 初始化后端 Go module。
- 初始化前端 Vite React TypeScript 项目。
- 添加 `docker-compose.yml`，提供 PostgreSQL 16 本地服务。
- 添加根目录 `README.md`，记录本地启动步骤。
- 准备 `.env.example`，分别说明后端和前端必需环境变量。
- 确认端口约定：
  - 后端：`localhost:8080`
  - 前端：`localhost:5173`
  - PostgreSQL：`localhost:5432`

建议命令：

```bash
cd backend
go mod init project-alpha/backend

cd ../frontend
npm create vite@latest . -- --template react-ts
```

交付物：

- `backend/go.mod`
- `frontend/package.json`
- `docker-compose.yml`
- `.env.example`
- `README.md`

验收标准：

- PostgreSQL 可以通过 Docker Compose 启动。
- 前端开发服务器可以启动。
- 后端空服务或健康检查可以启动。

### 4.2 M1：后端基础框架

目标：完成 Gin 服务骨架、配置加载、数据库连接和统一响应结构。

任务：

- 在 `backend/cmd/server/main.go` 中启动 Gin 服务。
- 在 `internal/config` 中读取配置：
  - `APP_ENV`
  - `HTTP_ADDR`
  - `DATABASE_URL`
  - `CORS_ALLOWED_ORIGINS`
- 在 `internal/database` 中初始化 GORM PostgreSQL 连接。
- 在 `internal/router` 中注册路由：
  - `GET /healthz`
  - `/api/v1` API 分组。
- 配置 CORS，允许前端开发地址访问。
- 在 `internal/response` 中定义统一成功和错误响应。
- 在 `internal/response` 或 `internal/handlers` 中统一映射错误码和 HTTP 状态码。

交付物：

- 可运行的 Gin 服务。
- `GET /healthz` 返回 `{ "status": "ok" }`。
- 数据库连接失败时服务启动应明确报错。

验收标准：

- `curl http://localhost:8080/healthz` 成功。
- 服务启动日志能看到监听地址。
- 配置缺失时错误信息明确，尤其是 `DATABASE_URL`。

### 4.3 M2：数据库模型与迁移

目标：实现 ticket、tag、ticket_tags 的 GORM 模型和自动迁移。

任务：

- 在 `internal/models` 中定义：
  - `Ticket`
  - `Tag`
  - 可选 `TicketTag` 显式 join model。
- 使用 `gorm.Model` 提供标准字段。
- 为 `Ticket.Tags` 配置 `many2many:ticket_tags`。
- 为关联配置 `OnUpdate:CASCADE,OnDelete:CASCADE`。
- 实现 `AutoMigrate`，启动时创建或更新表结构。
- 为标签名称添加唯一索引。
- 为列表查询相关字段添加索引：
  - `tickets.created_at`
  - `tickets.completed`
  - `ticket_tags.ticket_id`
  - `ticket_tags.tag_id`

交付物：

- GORM 模型。
- 数据库自动建表逻辑。
- 本地 PostgreSQL 中能看到 `tickets`、`tags`、`ticket_tags`。

验收标准：

- 启动后数据库表结构符合 `spec.md`。
- 重启服务不会重复创建异常表结构。
- ticket 与 tag 的多对多关系可正常写入。

### 4.4 M3：后端 Tag API

目标：完成标签 CRUD，为 ticket 功能提供基础数据。

任务：

- 定义请求 DTO：
  - `CreateTagRequest`
  - `UpdateTagRequest`
- 定义响应 DTO：
  - `TagResponse`
- 实现校验：
  - `name` trim 后必填，长度 1 到 40。
  - `color` 为空或符合 `#RRGGBB`。
  - `name` 唯一，冲突返回 `409 Conflict`。
- 实现 repository：
  - 创建标签。
  - 按名称检查重复。
  - 查询标签列表，按名称升序。
  - 按 ID 查询标签。
  - 更新标签。
  - 删除标签。
- 实现 service：
  - 处理 trim、唯一性、颜色校验和不存在资源。
- 实现 handler 和路由：
  - `POST /api/v1/tags`
  - `GET /api/v1/tags`
  - `PUT /api/v1/tags/{id}`
  - `DELETE /api/v1/tags/{id}`

交付物：

- 标签完整 CRUD API。
- 统一响应和错误结构。

验收标准：

- 可以创建、查询、更新、删除标签。
- 重复标签名称返回 `409`。
- 非法颜色返回 `400`。
- 删除不存在标签返回 `404`。

### 4.5 M4：后端 Ticket API

目标：完成 ticket 生命周期管理和标签关联。

任务：

- 定义请求 DTO：
  - `CreateTicketRequest`
  - `UpdateTicketRequest`
  - `ListTicketsQuery`
- 定义响应 DTO：
  - `TicketResponse`
  - `TicketListResponse`
  - `PaginationResponse`
- 实现校验：
  - `title` trim 后必填，长度 1 到 120。
  - `description` 最大 5000。
  - `tagIds` 必须为正整数数组。
  - `page` 默认 1。
  - `pageSize` 默认 20，最大 100。
- 实现 repository：
  - 创建 ticket。
  - 查询 ticket 详情并预加载 tags。
  - 更新 ticket 基础字段。
  - 删除 ticket。
  - 按 ID 查询 ticket。
  - 替换 ticket 标签集合。
  - 添加单个标签。
  - 移除单个标签。
  - 列表查询、分页、总数统计。
- 实现 service：
  - 创建 ticket 时校验 tagIds 全部存在。
  - 更新 ticket 时以 `tagIds` 替换完整标签集合。
  - 完成 ticket 时设置 `completed=true` 和 `completedAt=now`。
  - 取消完成时设置 `completed=false` 和 `completedAt=null`。
  - 添加标签保持幂等。
  - 删除 ticket 不删除标签。
- 实现 handler 和路由：
  - `POST /api/v1/tickets`
  - `GET /api/v1/tickets`
  - `GET /api/v1/tickets/{id}`
  - `PUT /api/v1/tickets/{id}`
  - `DELETE /api/v1/tickets/{id}`
  - `POST /api/v1/tickets/{id}/complete`
  - `POST /api/v1/tickets/{id}/uncomplete`
  - `POST /api/v1/tickets/{id}/tags/{tagId}`
  - `DELETE /api/v1/tickets/{id}/tags/{tagId}`

交付物：

- ticket 完整生命周期 API。
- ticket-tag 关联 API。
- 带 tags 的 ticket 响应。

验收标准：

- 创建 ticket 默认未完成。
- 完成和取消完成能正确维护 `completedAt`。
- 更新 ticket 能替换完整标签集合。
- 添加重复标签不会产生重复关联。
- 删除 ticket 后查询详情返回 `404`。
- 删除 ticket 后标签仍存在。

### 4.6 M5：搜索、筛选与分页

目标：完成 ticket 列表的核心查询能力。

任务：

- 支持 `title` 查询参数，使用 PostgreSQL `ILIKE`。
- 支持 `completed=true|false` 查询参数。
- 支持 `tagIds=1,2,3` 查询参数。
- 多标签筛选采用 AND 语义。
- 默认按 `createdAt` 倒序。
- 实现分页：
  - `page`
  - `pageSize`
  - `total`
  - `totalPages`
- 保证列表项预加载 tags。
- 对非法查询参数返回 `400`。

实现要点：

- 无标签筛选时直接查询 `tickets`。
- 有标签筛选时 join `ticket_tags`。
- 多标签 AND 可使用 `GROUP BY tickets.id HAVING COUNT(DISTINCT ticket_tags.tag_id) = ?`。
- 总数统计需要和列表查询保持相同过滤条件。

交付物：

- 稳定的 `GET /api/v1/tickets` 查询接口。

验收标准：

- title 搜索大小写不敏感。
- 单标签筛选结果正确。
- 多标签筛选只返回同时包含全部标签的 ticket。
- title、tagIds、completed 可以组合查询。
- 分页总数和当前页数据一致。

### 4.7 M6：后端测试

目标：保证后端核心业务不因后续前端接入而反复回归。

任务：

- 编写 service 测试：
  - 标签创建、重复名称、非法颜色。
  - ticket 创建、更新、删除。
  - 完成和取消完成。
  - 标签关联替换、添加、移除。
- 编写 handler 测试：
  - API 状态码。
  - 请求校验错误。
  - 统一响应结构。
- 编写 repository 或集成测试：
  - title 搜索。
  - 单标签筛选。
  - 多标签 AND 筛选。
  - 分页总数。

交付物：

- 后端测试用例。
- 可执行的测试命令记录到 README。

验收标准：

- `go test ./...` 通过。
- 关键业务规则均有测试覆盖。

### 4.8 M7：前端基础工程

目标：完成 Vite、Tailwind、shadcn/ui 和 API client 的基础搭建。

任务：

- 安装前端依赖。
- 配置 Tailwind CSS。
- 配置 TypeScript path alias：
  - `@/components`
  - `@/lib`
  - `@/types`
  - `@/hooks`
- 初始化 shadcn/ui，并设置 `rsc=false`。
- 添加首批 shadcn/ui 组件：
  - `button`
  - `input`
  - `textarea`
  - `badge`
  - `dialog`
  - `sheet`
  - `alert-dialog`
  - `select`
  - `popover`
  - `command`
  - `skeleton`
  - `sonner`
- 定义前端类型：
  - `Tag`
  - `Ticket`
  - `Pagination`
  - `ApiResponse`
- 实现 `src/lib/api.ts`：
  - 读取 `VITE_API_BASE_URL`。
  - 统一封装 GET、POST、PUT、DELETE。
  - 统一处理后端错误结构。

交付物：

- 可启动的前端工程。
- 可复用 API client。
- 与后端响应一致的 TypeScript 类型。

验收标准：

- `npm run dev` 正常启动。
- `npm run build` 通过。
- 前端能请求 `GET /healthz` 或 API 基础接口。

### 4.9 M8：前端标签功能

目标：完成标签管理 UI，并为 ticket 表单和筛选器提供数据。

任务：

- 实现 `useTags` hook：
  - 加载标签列表。
  - 创建标签。
  - 更新标签。
  - 删除标签。
  - 暴露 loading、error、refresh 状态。
- 实现标签管理组件：
  - 标签列表。
  - 新建标签表单。
  - 编辑标签表单。
  - 删除确认弹窗。
- 实现颜色输入：
  - 可使用原生 color input。
  - 同时保留十六进制文本校验。
- 标签名称重复时展示明确错误。

交付物：

- 标签管理弹窗或侧栏。
- 标签 CRUD 前端闭环。

验收标准：

- 用户可以在页面创建、编辑、删除标签。
- 删除标签前有确认。
- 操作成功后标签列表刷新。
- API 错误能以 toast 或表单错误展示。

### 4.10 M9：前端 Ticket 列表与筛选

目标：完成主页面的 ticket 列表、搜索、筛选和分页。

任务：

- 实现 `useTickets` hook：
  - 根据查询条件加载 ticket 列表。
  - 支持刷新。
  - 暴露 loading、error、pagination。
- 实现列表页面布局：
  - 顶部搜索和筛选工具栏。
  - ticket 列表。
  - 空状态。
  - 分页控制。
- 实现 title 搜索：
  - 输入框。
  - 300ms debounce。
- 实现标签多选筛选：
  - 使用 `Popover` + `Command` 或其他 shadcn/ui 组合。
  - 支持选择多个标签。
- 实现完成状态筛选：
  - 全部。
  - 未完成。
  - 已完成。
- 实现分页：
  - 上一页。
  - 下一页。
  - 当前页和总页数。

交付物：

- 可搜索、可筛选、可分页的 ticket 主列表。

验收标准：

- 页面初次加载展示 ticket 列表。
- title 输入后自动刷新列表。
- 标签筛选和完成状态筛选可组合使用。
- 加载中、空列表、错误状态都有明确展示。

### 4.11 M10：前端 Ticket 表单与操作

目标：完成 ticket 创建、编辑、删除、完成和取消完成。

任务：

- 实现 ticket 创建/编辑表单：
  - title。
  - description。
  - tagIds 多选。
- 实现表单校验：
  - title 必填且最大 120。
  - description 最大 5000。
- 实现创建 ticket。
- 实现编辑 ticket。
- 实现删除 ticket，并添加二次确认。
- 实现完成和取消完成快捷按钮。
- 操作成功后刷新列表。
- 操作失败时展示 toast 或字段错误。

交付物：

- ticket 创建/编辑弹窗或抽屉。
- ticket 列表中的完成、取消完成、删除入口。

验收标准：

- 用户可以完整管理 ticket 生命周期。
- 创建或编辑 ticket 时可以设置标签。
- 完成状态变化能立即反映到列表。
- 删除 ticket 后列表不再展示该 ticket。

### 4.12 M11：前后端联调

目标：完成真实 PostgreSQL、真实后端、真实前端之间的端到端联通。

任务：

- 使用 Docker Compose 启动 PostgreSQL。
- 启动后端服务并确认 AutoMigrate 成功。
- 启动前端服务并配置 `VITE_API_BASE_URL`。
- 完成手动验收数据流：
  - 创建标签。
  - 创建 ticket。
  - 编辑 ticket。
  - 添加和移除标签。
  - title 搜索。
  - 标签筛选。
  - 完成和取消完成。
  - 删除标签。
  - 删除 ticket。
- 检查浏览器控制台无明显错误。
- 检查后端日志无未处理异常。

交付物：

- 前后端可本地联调运行。
- README 中补齐联调步骤。

验收标准：

- `spec.md` 中的手动验收场景全部通过。
- 重启服务后数据仍存在。

### 4.13 M12：前端测试与最终验收

目标：为主要 UI 行为提供基础保障，并完成交付前检查。

任务：

- 添加前端测试工具。
- 覆盖关键交互：
  - ticket 列表渲染。
  - 搜索输入触发查询。
  - 标签筛选交互。
  - 创建/编辑表单校验。
  - 删除确认。
  - API 错误提示。
- 运行前端构建。
- 运行后端测试。
- 整理 README。
- 对照 `spec.md` 的验收标准逐项检查。

交付物：

- 前端测试用例。
- 最终 README。
- 可通过测试和构建的项目。

验收标准：

- `go test ./...` 通过。
- `npm run build` 通过。
- 前端主要操作路径可手动完成。

## 5. 任务依赖关系

推荐执行顺序：

```text
M0 -> M1 -> M2 -> M3 -> M4 -> M5 -> M6 -> M7 -> M8 -> M9 -> M10 -> M11 -> M12
```

关键依赖：

- M2 依赖 M1 的数据库连接。
- M4 依赖 M3 的标签基础能力。
- M5 依赖 M4 的 ticket 查询能力。
- M8、M9、M10 依赖 M7 的前端基础工程。
- M11 依赖后端 API 和前端页面基本完成。

可并行项：

- M6 的部分 service 测试可在 M3、M4 完成后同步补充。
- M7 可在后端 M3 到 M5 开发期间并行启动。
- README 可以随每个里程碑增量更新。

## 6. 后端实现细节清单

### 6.1 配置

- `DATABASE_URL` 必须存在。
- `HTTP_ADDR` 默认 `:8080`。
- `CORS_ALLOWED_ORIGINS` 开发默认 `http://localhost:5173`。
- 配置读取失败应在启动阶段失败，不应延迟到请求阶段。

### 6.2 响应结构

所有业务 API 使用统一结构：

```json
{
  "data": {},
  "error": null
}
```

错误结构：

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "title is required",
    "details": {}
  }
}
```

`DELETE` 成功接口返回 `204 No Content`，不需要响应体。

### 6.3 错误映射

- 参数绑定失败：`400 VALIDATION_ERROR`
- 业务校验失败：`400 VALIDATION_ERROR`
- 资源不存在：`404 NOT_FOUND`
- 标签名称重复：`409 CONFLICT`
- 数据库异常：`500 DATABASE_ERROR`
- 未预期异常：`500 INTERNAL_ERROR`

### 6.4 事务

必须使用事务的操作：

- 创建 ticket 并设置标签。
- 更新 ticket 并替换标签集合。
- 删除 ticket 并清理关联。
- 删除 tag 并清理关联。
- 批量添加或移除关联。

### 6.5 查询

ticket 列表查询必须覆盖：

- 无条件查询。
- title 模糊查询。
- completed 筛选。
- 单标签筛选。
- 多标签 AND 筛选。
- 组合筛选。
- 分页。
- tags 预加载。

## 7. 前端实现细节清单

### 7.1 页面布局

主页面直接展示操作界面，不做落地页：

- 顶部：应用名称、标签管理入口、新建 ticket 按钮。
- 工具栏：title 搜索、标签筛选、完成状态筛选。
- 主区域：ticket 列表。
- 底部：分页。

### 7.2 组件拆分

建议组件：

- `TicketPage`
- `TicketToolbar`
- `TicketList`
- `TicketListItem`
- `TicketFormDialog`
- `TicketDeleteDialog`
- `TagFilter`
- `StatusFilter`
- `TagManagerDialog`
- `TagForm`
- `TagBadge`
- `PaginationControls`

### 7.3 API 方法

建议封装：

- `listTickets(params)`
- `getTicket(id)`
- `createTicket(payload)`
- `updateTicket(id, payload)`
- `deleteTicket(id)`
- `completeTicket(id)`
- `uncompleteTicket(id)`
- `addTicketTag(ticketId, tagId)`
- `removeTicketTag(ticketId, tagId)`
- `listTags()`
- `createTag(payload)`
- `updateTag(id, payload)`
- `deleteTag(id)`

### 7.4 状态处理

每个主要请求至少处理：

- `idle`
- `loading`
- `success`
- `error`

用户可见状态：

- 加载时展示 skeleton 或按钮 loading。
- 空结果展示空状态。
- 失败时展示 toast 或 inline error。
- 删除和危险操作使用确认弹窗。

## 8. 验收检查表

### 8.1 后端

- [ ] `GET /healthz` 正常。
- [ ] `POST /api/v1/tags` 正常。
- [ ] `GET /api/v1/tags` 正常。
- [ ] `PUT /api/v1/tags/{id}` 正常。
- [ ] `DELETE /api/v1/tags/{id}` 正常。
- [ ] `POST /api/v1/tickets` 正常。
- [ ] `GET /api/v1/tickets` 正常。
- [ ] `GET /api/v1/tickets/{id}` 正常。
- [ ] `PUT /api/v1/tickets/{id}` 正常。
- [ ] `DELETE /api/v1/tickets/{id}` 正常。
- [ ] `POST /api/v1/tickets/{id}/complete` 正常。
- [ ] `POST /api/v1/tickets/{id}/uncomplete` 正常。
- [ ] `POST /api/v1/tickets/{id}/tags/{tagId}` 正常。
- [ ] `DELETE /api/v1/tickets/{id}/tags/{tagId}` 正常。
- [ ] title 搜索大小写不敏感。
- [ ] 多标签筛选为 AND 语义。
- [ ] 分页 total 和 items 正确。
- [ ] `go test ./...` 通过。

### 8.2 前端

- [ ] 页面可以加载 ticket 列表。
- [ ] 页面可以加载标签列表。
- [ ] 可以创建、编辑、删除标签。
- [ ] 可以创建、编辑、删除 ticket。
- [ ] 可以完成和取消完成 ticket。
- [ ] 可以为 ticket 设置多个标签。
- [ ] 可以按 title 搜索。
- [ ] 可以按标签筛选。
- [ ] 可以按完成状态筛选。
- [ ] 可以组合筛选。
- [ ] 删除操作有确认弹窗。
- [ ] 表单校验错误展示清晰。
- [ ] API 错误展示清晰。
- [ ] `npm run build` 通过。

### 8.3 集成

- [ ] 前端通过真实 API 完成所有操作。
- [ ] PostgreSQL 重启后数据仍存在。
- [ ] 浏览器控制台无明显错误。
- [ ] 后端日志无未处理 panic。
- [ ] README 启动步骤可复现。

## 9. 风险与应对

### 9.1 多标签 AND 查询复杂

风险：列表数据和 total 统计条件不一致，导致分页错误。

应对：

- 抽取统一查询构建函数。
- 为列表和 total 复用相同 filter 条件。
- 单独编写多标签筛选测试。

### 9.2 GORM 多对多关联重复写入

风险：重复添加标签导致 join table 冲突或重复数据。

应对：

- 使用联合主键或唯一约束。
- 添加前先查询或使用 GORM association 幂等逻辑。
- 添加重复标签的测试。

### 9.3 软删除与关联清理不一致

风险：软删除 ticket 或 tag 后，关联表仍影响筛选。

应对：

- 删除前显式清理 association。
- 查询时只基于未删除的 tickets 和 tags。
- 增加删除后筛选测试。

### 9.4 前端状态刷新遗漏

风险：创建、编辑、删除后列表或标签筛选器未刷新。

应对：

- 所有 mutation 成功后统一调用对应 `refresh`。
- 标签变更后刷新 tags，并刷新 tickets。
- 将刷新逻辑集中在 hooks 或页面容器中。

## 10. 建议提交节奏

建议按里程碑提交，便于回滚和 review：

1. 项目脚手架和本地开发环境。
2. 后端基础框架、配置、数据库连接。
3. GORM 模型和 AutoMigrate。
4. Tag API。
5. Ticket API。
6. 搜索、筛选、分页。
7. 后端测试。
8. 前端基础工程和 shadcn/ui。
9. 标签管理 UI。
10. ticket 列表、筛选和分页 UI。
11. ticket 表单和操作 UI。
12. 联调、测试和 README。

## 11. 最终交付物

- 可运行后端服务。
- 可运行前端应用。
- PostgreSQL 本地开发配置。
- API 与前端交互完整联通。
- 后端测试与前端构建通过。
- README 包含本地启动、测试和联调步骤。
- `specs/spec.md` 与 `specs/implementation-plan.md` 作为项目规划文档。
