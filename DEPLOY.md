# 公网部署指南

把 project-alpha 一键部署到公网，三层独立、域名固定、零费用。

## 架构

```
+-------------------+       HTTPS        +-------------------+
|  Cloudflare Pages |  ----------------> |   Fly.io Backend  |
|  (React 前端)     |                    |   (Go + Gin)      |
|  ...pages.dev     |                    |   ...fly.dev      |
+-------------------+                    +---------+---------+
         |                                          |
         | 静态资源 CDN 加速                          | TLS + SSL
         |                                          v
         |                                 +-------------------+
         |                                 |   Neon Postgres   |
         |                                 |   (免费 0.5GB)    |
         |                                 +-------------------+
```

## 一次性账号注册（约 3 分钟）

| 服务 | 用途 | 注册链接 |
|------|------|----------|
| Fly.io | 跑 Go 后端 | https://fly.io/app/sign-up （推荐 GitHub 一键登录） |
| Cloudflare | 跑前端 + CDN | https://dash.cloudflare.com/sign-up |
| Neon | 免费 PostgreSQL | https://console.neon.tech/sign-up |

## 数据库准备（30 秒）

1. 登录 Neon 控制台 → **New Project**
2. Region 选 `Asia Pacific (Singapore)` 或靠近你的区域
3. 复制 **Connection string**，形如：
   ```
   postgres://user:pass@ep-xxx.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
   ```

## 一条命令部署

在项目根目录执行：

```bash
DATABASE_URL='postgres://user:pass@ep-xxx.ap-southeast-1.aws.neon.tech/neondb?sslmode=require' \
    ./scripts/deploy.sh
```

脚本会自动：

1. 验证 `fly` 和 `wrangler` 是否登录
2. 在 Fly.io 创建 app `project-alpha-shilianjie`（如果不存在）
3. 推送 Docker 镜像并部署后端
4. 等后端 `/healthz` 绿
5. 把 `seed.sql` 写入生产数据库（10 个标签 + 10 个 ticket）
6. 构建前端（注入生产 API 地址）
7. 部署到 Cloudflare Pages
8. 打印两个公网 URL

## 部署后

固定公网域名：

- 前端：https://project-alpha-shilianjie.pages.dev
- 后端：https://project-alpha-shilianjie.fly.dev
- 健康检查：https://project-alpha-shilianjie.fly.dev/healthz

后续运维：

```bash
# 查看后端日志
fly logs --app project-alpha-shilianjie

# 重新部署后端（代码改动后）
./scripts/deploy.sh

# 验证线上状态
./scripts/healthcheck.sh

# 重新灌种子（会清空现有数据）
DATABASE_URL=... go -C backend run ./cmd/seed
```

## 自定义域名（可选，永久更稳）

如果希望用 `yourname.com` 而不是 `*.pages.dev`：

1. 买个域名（`.top` / `.xyz` 一年几块钱）
2. Cloudflare Pages → Custom domain → 添加 `yourname.com`
3. 后端同理：`fly certs add api.yourname.com`

## 重新部署

任意代码改动后：

```bash
./scripts/deploy.sh   # 会跳过已存在的 app/project，只更新镜像
```

## 故障排查

| 现象 | 解决 |
|------|------|
| `fly: command not found` | `brew install flyctl` 或 `curl -L https://fly.io/install.sh \| sh` |
| `wrangler: command not found` | `npm i -g wrangler` |
| `Not logged in to Fly.io` | `fly auth login` |
| `Not logged in to Cloudflare` | `wrangler login` |
| Backend 部署成功但前端 404 | Cloudflare Pages 第一次需要 1-2 分钟传播 |
| 种子报错 "tables do not exist" | 后端没启动成功，等 `/healthz` 绿了再跑种子 |
| 前端请求后端 CORS 错误 | 检查 `fly secrets list --app project-alpha-shilianjie`，确保 `DATABASE_URL` 设了；CORS 默认是 `*`，不应该报错 |
