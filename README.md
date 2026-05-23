# project-alpha

轻量级 project/ticket 管理工具，支持标签分类、ticket 管理、标题搜索、标签筛选、完成状态筛选和分页。

## 技术栈

- 后端：Go、Gin、GORM、PostgreSQL
- 前端：TypeScript、Vite、React、Tailwind CSS、shadcn/ui 风格组件
- 测试：Go test、Vitest、React Testing Library

## 本地运行

使用 Makefile 启动完整本地环境：

```bash
make install
make dev
```

`make dev` 会启动 PostgreSQL、后端和前端。也可以单独启动：

```bash
make db
make backend
make frontend
```

启动 PostgreSQL：

```bash
docker compose up -d postgres
```

启动后端：

```bash
cd backend
export DATABASE_URL="host=localhost user=project_alpha password=project_alpha dbname=project_alpha port=5432 sslmode=disable TimeZone=Asia/Shanghai"
export CORS_ALLOWED_ORIGINS="http://localhost:5173"
go run ./cmd/server
```

启动前端：

```bash
cd frontend
npm install
cp .env.example .env
npm run dev
```

默认地址：

- 前端：`http://localhost:5173`
- 后端健康检查：`http://localhost:8080/healthz`
- API：`http://localhost:8080/api/v1`

## 前后端联调场景

1. 在页面打开标签管理，创建标签。
2. 创建 ticket，并给 ticket 选择标签。
3. 编辑 ticket 的标题、描述和标签。
4. 使用标题搜索、标签筛选、完成状态筛选组合查询。
5. 完成 ticket，再取消完成。
6. 删除标签，确认 ticket 列表刷新后不再显示该标签。
7. 删除 ticket，确认列表中不再展示。

## 后端测试

```bash
cd backend
go test ./...
```

如果本地 Go 模块缓存缺少依赖，请先恢复网络后执行：

```bash
cd backend
go mod tidy
go test ./...
```

## 前端验证

```bash
cd frontend
npm install
npm run build
npm test
```

## 代码检查

首次使用 pre-commit 前先安装并初始化 Git hook：

```bash
python3 -m pip install --user pre-commit
pre-commit install
```

手动执行全量检查：

```bash
pre-commit run --all-files
```

也可以使用 Makefile：

```bash
make fmt
make lint
make test
make precommit
```

GitHub Actions 会在 push 到 `main` 和 pull request 时运行 pre-commit、后端 Go 检查和前端 TypeScript/Vite 测试构建。
