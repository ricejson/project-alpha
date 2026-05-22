# project-alpha 需求与设计文档

## 1. 项目概述

project-alpha 是一个轻量级项目管理工具，面向无需账号体系的小型团队或个人使用场景。系统核心围绕 ticket 管理与标签分类展开，用户可以创建、编辑、删除、完成、取消完成 ticket，并通过标签或标题快速检索 ticket。

本项目采用前后端分离架构：

- 后端：Go + Gin，提供 RESTful JSON API。
- ORM：GORM，连接 PostgreSQL。
- 前端：TypeScript + Vite + React + Tailwind CSS + shadcn/ui。
- 数据库：PostgreSQL。

当前版本不包含用户系统、权限控制、多项目空间、评论、附件、通知、审计日志等功能。

## 2. 目标与非目标

### 2.1 目标

- 支持 ticket 的完整生命周期管理。
- 支持标签的创建、编辑、删除与分类使用。
- 支持为 ticket 添加和移除多个标签。
- 支持按照标签筛选 ticket。
- 支持按照 title 搜索 ticket。
- 提供清晰、稳定、易扩展的 REST API。
- 提供简洁、可用、响应式的 Web 操作界面。
- 使用 PostgreSQL 持久化数据。

### 2.2 非目标

- 不实现登录、注册、用户权限或多租户隔离。
- 不实现复杂工作流、状态机、看板列或冲刺管理。
- 不实现文件上传、富文本编辑、评论和通知。
- 不实现实时协作或 WebSocket。
- 不实现全文搜索引擎，标题搜索使用数据库模糊匹配即可。

## 3. 用户角色与使用场景

本系统只有一种隐含角色：工具使用者。

典型使用场景：

- 用户创建一个待处理 ticket，填写标题、描述并选择标签。
- 用户在 ticket 列表中按标签查看某类事项。
- 用户通过标题关键字查找已有 ticket。
- 用户完成某个 ticket 后将其标记为完成。
- 用户误标完成后取消完成。
- 用户维护标签名称和颜色，使 ticket 分类更易识别。

## 4. 功能需求

### 4.1 Ticket 管理

系统必须支持：

- 创建 ticket。
- 查看 ticket 列表。
- 查看 ticket 详情。
- 编辑 ticket。
- 删除 ticket。
- 将 ticket 标记为完成。
- 将已完成 ticket 取消完成。

ticket 字段：

- `id`：唯一标识。
- `title`：标题，必填。
- `description`：描述，可选。
- `completed`：是否完成。
- `completedAt`：完成时间，可选。
- `tags`：关联标签列表。
- `createdAt`：创建时间。
- `updatedAt`：更新时间。

业务规则：

- `title` 不能为空，去除首尾空格后长度必须在 1 到 120 个字符之间。
- `description` 可为空，最大长度 5000 个字符。
- 新建 ticket 默认 `completed=false`，`completedAt=null`。
- 标记完成时设置 `completed=true`，并写入当前时间到 `completedAt`。
- 取消完成时设置 `completed=false`，并清空 `completedAt`。
- 删除 ticket 时必须同时移除 ticket 与标签之间的关联关系。
- 删除 ticket 不应删除标签本身。

### 4.2 标签管理

系统必须支持：

- 创建标签。
- 查看标签列表。
- 编辑标签。
- 删除标签。

标签字段：

- `id`：唯一标识。
- `name`：标签名称，必填且唯一。
- `color`：标签颜色，可选，使用十六进制色值。
- `createdAt`：创建时间。
- `updatedAt`：更新时间。

业务规则：

- `name` 不能为空，去除首尾空格后长度必须在 1 到 40 个字符之间。
- `name` 在系统内唯一，建议大小写不敏感唯一。
- `color` 如提供，必须符合 `#RRGGBB` 格式。
- 删除标签时必须移除它与所有 ticket 的关联关系。
- 删除标签不应删除任何 ticket。

### 4.3 Ticket 标签关联

系统必须支持：

- 给 ticket 添加标签。
- 从 ticket 移除标签。
- 创建或编辑 ticket 时一次性设置标签列表。

业务规则：

- 一个 ticket 可以关联 0 到多个标签。
- 一个标签可以关联 0 到多个 ticket。
- 重复添加同一个标签应保持幂等，不创建重复关系。
- 给不存在的 ticket 或标签建立关联时返回 404。
- 支持通过替换标签 ID 列表的方式更新 ticket 的标签集合。

### 4.4 搜索与筛选

系统必须支持：

- 按 title 搜索 ticket。
- 按一个或多个标签筛选 ticket。
- 支持按照完成状态筛选 ticket。
- 支持组合查询：title + tags + completed。

查询规则：

- title 搜索使用不区分大小写的模糊匹配。
- 标签筛选支持传入多个标签 ID。
- 当传入多个标签时，默认语义为 AND，即返回同时拥有所有指定标签的 ticket。
- 列表默认按 `createdAt` 倒序排序。
- 支持分页，避免列表无限增长。

## 5. 页面与交互设计

### 5.1 页面结构

前端采用单页应用。第一阶段建议包含以下视图：

- Ticket 列表主页面。
- Ticket 创建/编辑弹窗或抽屉。
- Ticket 详情面板。
- 标签管理弹窗或侧栏。

无需登录页和用户设置页。

### 5.2 Ticket 列表页

主要区域：

- 顶部工具栏：title 搜索框、标签筛选器、完成状态筛选器、新建 ticket 按钮。
- Ticket 列表：展示 title、description 摘要、标签、完成状态、更新时间。
- 空状态：无 ticket 或搜索无结果时显示简洁提示与创建入口。
- 分页控件：上一页、下一页、当前页信息。

交互要求：

- 输入 title 关键字后触发搜索，建议使用 300ms debounce。
- 标签筛选器支持多选。
- 完成状态筛选支持全部、未完成、已完成。
- 点击 ticket 行打开详情或编辑入口。
- 完成/取消完成操作应可从列表快速触发。
- 删除 ticket 前必须二次确认。

### 5.3 Ticket 创建/编辑

表单字段：

- title 输入框。
- description 多行文本框。
- 标签多选控件。

交互要求：

- 创建成功后关闭表单并刷新列表。
- 编辑成功后关闭表单并刷新当前数据。
- 表单提交时展示 loading 状态。
- 校验错误应显示在对应字段附近。
- 后端返回错误时应显示可读提示。

### 5.4 标签管理

标签管理界面应支持：

- 查看所有标签。
- 新建标签。
- 修改标签名称和颜色。
- 删除标签。

交互要求：

- 标签颜色使用颜色选择器或预设色板。
- 删除标签前必须二次确认。
- 标签名称重复时显示明确错误。

### 5.5 UI 组件建议

基于 shadcn/ui 的组件组合：

- `Button`：主要操作、次要操作、危险操作。
- `Input`：title 搜索和表单输入。
- `Textarea`：ticket description。
- `Badge`：标签展示。
- `Dialog` 或 `Sheet`：创建/编辑 ticket、标签管理。
- `Select` 或 `Tabs`：完成状态筛选。
- `Command` 或 `Popover`：标签多选筛选。
- `AlertDialog`：删除确认。
- `Toast` 或 `Sonner`：操作结果反馈。
- `Skeleton`：加载态。

前端应避免营销式页面结构，默认展示可直接操作的项目管理界面。

## 6. API 设计

### 6.1 API 基础约定

后端 API 统一挂载在：

```text
/api/v1
```

响应格式建议：

```json
{
  "data": {},
  "error": null
}
```

错误响应建议：

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

HTTP 状态码：

- `200 OK`：查询、更新、状态变更成功。
- `201 Created`：创建成功。
- `204 No Content`：删除成功。
- `400 Bad Request`：请求格式或参数错误。
- `404 Not Found`：资源不存在。
- `409 Conflict`：唯一约束冲突，例如标签名称重复。
- `500 Internal Server Error`：未预期服务端错误。

### 6.2 Ticket API

#### 创建 ticket

```http
POST /api/v1/tickets
```

请求体：

```json
{
  "title": "Implement ticket list",
  "description": "Build ticket list with filters",
  "tagIds": [1, 2]
}
```

响应：`201 Created`

```json
{
  "data": {
    "id": 1,
    "title": "Implement ticket list",
    "description": "Build ticket list with filters",
    "completed": false,
    "completedAt": null,
    "tags": [
      {
        "id": 1,
        "name": "frontend",
        "color": "#2563EB"
      }
    ],
    "createdAt": "2026-05-22T12:00:00Z",
    "updatedAt": "2026-05-22T12:00:00Z"
  },
  "error": null
}
```

#### 查询 ticket 列表

```http
GET /api/v1/tickets?title=list&tagIds=1,2&completed=false&page=1&pageSize=20
```

查询参数：

- `title`：可选，标题关键字。
- `tagIds`：可选，逗号分隔的标签 ID 列表。
- `completed`：可选，`true` 或 `false`。
- `page`：可选，默认 `1`。
- `pageSize`：可选，默认 `20`，最大 `100`。

响应：`200 OK`

```json
{
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 0,
      "totalPages": 0
    }
  },
  "error": null
}
```

#### 查询 ticket 详情

```http
GET /api/v1/tickets/{id}
```

#### 更新 ticket

```http
PUT /api/v1/tickets/{id}
```

请求体：

```json
{
  "title": "Implement ticket search",
  "description": "Support title and tag search",
  "tagIds": [1, 3]
}
```

说明：

- `tagIds` 表示替换该 ticket 的完整标签集合。
- 不建议通过该接口修改完成状态，完成状态使用专门接口。

#### 删除 ticket

```http
DELETE /api/v1/tickets/{id}
```

响应：`204 No Content`

#### 标记 ticket 完成

```http
POST /api/v1/tickets/{id}/complete
```

#### 取消 ticket 完成

```http
POST /api/v1/tickets/{id}/uncomplete
```

#### 给 ticket 添加标签

```http
POST /api/v1/tickets/{id}/tags/{tagId}
```

#### 从 ticket 移除标签

```http
DELETE /api/v1/tickets/{id}/tags/{tagId}
```

### 6.3 Tag API

#### 创建标签

```http
POST /api/v1/tags
```

请求体：

```json
{
  "name": "frontend",
  "color": "#2563EB"
}
```

#### 查询标签列表

```http
GET /api/v1/tags
```

建议默认按名称升序排序。

#### 更新标签

```http
PUT /api/v1/tags/{id}
```

请求体：

```json
{
  "name": "backend",
  "color": "#16A34A"
}
```

#### 删除标签

```http
DELETE /api/v1/tags/{id}
```

响应：`204 No Content`

## 7. 数据库设计

### 7.1 表结构

#### tickets

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | primary key | 主键 |
| title | varchar(120) | not null | 标题 |
| description | text | not null default '' | 描述 |
| completed | boolean | not null default false | 是否完成 |
| completed_at | timestamptz | null | 完成时间 |
| created_at | timestamptz | not null | 创建时间 |
| updated_at | timestamptz | not null | 更新时间 |
| deleted_at | timestamptz | null, indexed | 软删除时间 |

#### tags

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | primary key | 主键 |
| name | varchar(40) | not null, unique | 标签名称 |
| color | varchar(7) | null | 十六进制颜色 |
| created_at | timestamptz | not null | 创建时间 |
| updated_at | timestamptz | not null | 更新时间 |
| deleted_at | timestamptz | null, indexed | 软删除时间 |

#### ticket_tags

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| ticket_id | bigint | primary key, foreign key | ticket ID |
| tag_id | bigint | primary key, foreign key | tag ID |
| created_at | timestamptz | not null | 关联创建时间 |

约束：

- `ticket_tags.ticket_id` 引用 `tickets.id`，删除 ticket 时级联删除关联。
- `ticket_tags.tag_id` 引用 `tags.id`，删除 tag 时级联删除关联。
- `(ticket_id, tag_id)` 作为联合主键，避免重复关联。

### 7.2 GORM 模型建议

GORM 文档提供 `gorm.Model`，包含 `ID`、`CreatedAt`、`UpdatedAt` 和带索引的 `DeletedAt` 字段。当前项目可以嵌入该结构以获得标准字段和软删除能力。

```go
type Ticket struct {
    gorm.Model
    Title       string     `gorm:"size:120;not null"`
    Description string     `gorm:"type:text;not null;default:''"`
    Completed   bool       `gorm:"not null;default:false"`
    CompletedAt *time.Time
    Tags        []Tag      `gorm:"many2many:ticket_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Tag struct {
    gorm.Model
    Name  string   `gorm:"size:40;not null;uniqueIndex"`
    Color *string  `gorm:"size:7"`
    Tickets []Ticket `gorm:"many2many:ticket_tags;"`
}
```

说明：

- `many2many:ticket_tags` 用于声明 ticket 与 tag 的多对多关系。
- 查询 ticket 列表和详情时应预加载 tags，避免前端获得不完整数据。
- 使用 PostgreSQL 时，建议通过环境变量配置 DSN。

### 7.3 索引建议

- `tickets.deleted_at`：支持软删除查询。
- `tickets.created_at`：支持默认排序。
- `tickets.completed`：支持完成状态筛选。
- `tags.name`：唯一索引。
- `ticket_tags.ticket_id`：支持按 ticket 查询标签。
- `ticket_tags.tag_id`：支持按标签筛选 ticket。

如果 title 搜索性能成为问题，可后续增加 `pg_trgm` 和 GIN 索引；第一阶段无需实现。

## 8. 后端设计

### 8.1 目录结构建议

```text
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
```

职责划分：

- `config`：读取环境变量和配置默认值。
- `database`：初始化 PostgreSQL 连接和迁移。
- `models`：GORM 模型定义。
- `repositories`：数据库查询封装。
- `services`：业务规则和事务。
- `handlers`：Gin HTTP 请求处理。
- `router`：路由注册和中间件。
- `response`：统一响应和错误结构。
- `validation`：请求校验辅助函数。

### 8.2 Gin 约定

根据 Gin 当前文档：

- 使用 `gin.Default()` 创建路由器，默认包含 logger 和 recovery 中间件。
- 使用 `r.Group("/api/v1")` 注册版本化 API。
- 使用 `ShouldBindJSON` 绑定和验证 JSON 请求体。
- 使用 `c.JSON(status, body)` 返回 JSON。
- 使用 `gin-contrib/cors` 配置前端跨域访问。

后端应提供健康检查：

```http
GET /healthz
```

响应：

```json
{
  "status": "ok"
}
```

### 8.3 配置项

环境变量：

- `APP_ENV`：运行环境，默认 `development`。
- `HTTP_ADDR`：监听地址，默认 `:8080`。
- `DATABASE_URL`：PostgreSQL DSN，必填。
- `CORS_ALLOWED_ORIGINS`：允许的前端 Origin，开发默认 `http://localhost:5173`。

### 8.4 事务边界

以下操作建议使用事务：

- 创建 ticket 并关联标签。
- 更新 ticket 并替换标签集合。
- 删除 ticket 并清理关联。
- 删除 tag 并清理关联。
- 批量添加或移除 ticket 标签。

### 8.5 查询实现建议

列表查询基础逻辑：

- 从 `tickets` 表开始查询。
- 如传入 `title`，使用 `ILIKE '%keyword%'`。
- 如传入 `completed`，增加布尔过滤条件。
- 如传入 `tagIds`，通过 `ticket_tags` 关联表筛选。
- 多标签 AND 语义可使用 `GROUP BY tickets.id HAVING COUNT(DISTINCT ticket_tags.tag_id) = ?`。
- 使用 `Preload("Tags")` 返回标签数据。
- 使用 `Limit`、`Offset` 和总数查询实现分页。

## 9. 前端设计

### 9.1 目录结构建议

```text
frontend/
  src/
    app/
      App.tsx
    components/
      tickets/
      tags/
      layout/
      ui/
    lib/
      api.ts
      utils.ts
    types/
      ticket.ts
      tag.ts
    hooks/
      useTickets.ts
      useTags.ts
    main.tsx
    index.css
```

### 9.2 shadcn/ui 与 Tailwind 约定

根据 shadcn/ui 当前文档：

- 使用 `npx shadcn@latest init` 初始化组件配置。
- Vite 项目应配置路径别名，例如 `@/components`、`@/lib/utils`。
- `components.json` 应设置 `rsc=false`，因为当前项目是 Vite SPA。
- 组件通过 CLI 添加到本地代码库，作为可维护的本地组件使用。

建议初始安装组件：

```text
button input textarea badge dialog sheet alert-dialog select popover command skeleton sonner
```

### 9.3 前端状态管理

第一阶段无需引入复杂全局状态库。建议：

- 使用 React hooks 管理页面局部状态。
- 使用自定义 hooks 封装 API 请求。
- 可选使用 TanStack Query 管理服务端状态、缓存和自动刷新；如不引入，应保证加载态、错误态和刷新逻辑清晰。

### 9.4 API Client

前端应集中封装 HTTP 请求：

- 基础地址从 `VITE_API_BASE_URL` 读取，默认 `http://localhost:8080/api/v1`。
- 所有请求统一处理 JSON 解析。
- 统一处理后端错误结构。
- 暴露类型化方法，例如 `listTickets`、`createTicket`、`updateTicket`、`completeTicket`、`listTags`。

### 9.5 TypeScript 类型

```ts
export interface Tag {
  id: number
  name: string
  color: string | null
  createdAt: string
  updatedAt: string
}

export interface Ticket {
  id: number
  title: string
  description: string
  completed: boolean
  completedAt: string | null
  tags: Tag[]
  createdAt: string
  updatedAt: string
}
```

## 10. 校验与错误处理

### 10.1 后端校验

请求体 DTO 应使用 Gin binding tag 和自定义校验：

- `title`：required，trim 后非空，最大 120。
- `description`：最大 5000。
- `tagIds`：必须是正整数数组。
- `name`：required，trim 后非空，最大 40。
- `color`：为空或匹配 `^#[0-9A-Fa-f]{6}$`。

### 10.2 前端校验

前端应在提交前做基础校验：

- 必填字段非空。
- 最大长度限制。
- 颜色格式限制。

后端仍作为最终可信校验来源。

### 10.3 错误码建议

- `VALIDATION_ERROR`
- `NOT_FOUND`
- `CONFLICT`
- `DATABASE_ERROR`
- `INTERNAL_ERROR`

## 11. 开发与运行设计

### 11.1 本地开发服务

建议端口：

- 后端：`http://localhost:8080`
- 前端：`http://localhost:5173`
- PostgreSQL：`localhost:5432`

### 11.2 Docker Compose 建议

本地开发可使用 Docker Compose 启动 PostgreSQL：

```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: project_alpha
      POSTGRES_PASSWORD: project_alpha
      POSTGRES_DB: project_alpha
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### 11.3 数据迁移

第一阶段可以使用 GORM `AutoMigrate` 创建表结构，便于快速开发。

后续如项目扩大，建议迁移到显式 migration 工具，例如 golang-migrate，以保证生产环境 schema 变更可追踪、可回滚。

## 12. 测试策略

### 12.1 后端测试

建议覆盖：

- Ticket 创建、更新、删除。
- 完成和取消完成。
- 标签创建、更新、删除。
- 添加和移除 ticket 标签。
- title 搜索。
- 单标签和多标签筛选。
- 标签重复名称冲突。
- 不存在资源返回 404。

测试层级：

- Service 单元测试：验证业务规则。
- Handler 集成测试：验证 HTTP 状态码、请求体和响应体。
- Repository 集成测试：可使用测试 PostgreSQL 或容器化数据库。

### 12.2 前端测试

建议覆盖：

- Ticket 列表渲染。
- 搜索和筛选交互。
- 创建/编辑表单校验。
- 完成/取消完成按钮行为。
- 删除确认弹窗。
- API 错误提示。

### 12.3 手动验收

验收场景：

1. 创建两个标签：frontend、backend。
2. 创建三个 ticket，分别关联不同标签。
3. 按 frontend 筛选，只显示带 frontend 的 ticket。
4. 按 title 关键字搜索，只显示匹配标题的 ticket。
5. 组合 title 和标签筛选，结果正确。
6. 将一个 ticket 标记完成，列表状态立即更新。
7. 取消完成，状态和完成时间正确恢复。
8. 删除标签后，相关 ticket 仍存在，但不再显示该标签。
9. 删除 ticket 后，列表和详情接口均无法再查询该 ticket。

## 13. 验收标准

项目达到以下条件即可认为第一阶段完成：

- 可以通过前端完成 ticket 的创建、编辑、删除、完成和取消完成。
- 可以通过前端创建、编辑和删除标签。
- 可以通过前端给 ticket 添加和移除标签。
- 可以按 title 搜索 ticket。
- 可以按一个或多个标签筛选 ticket。
- 后端 API 返回结构稳定，错误响应可读。
- PostgreSQL 中数据持久化正确，重启服务后数据仍存在。
- 前端在桌面和移动视口下均可正常操作。
- README 或启动说明中包含本地运行步骤。

## 14. 后续扩展方向

可在第一阶段稳定后考虑：

- 多项目空间。
- 用户系统和权限。
- ticket 优先级、截止日期、负责人。
- 评论和活动记录。
- 看板视图。
- 批量操作。
- 全文搜索和高级筛选。
- 导入导出。
