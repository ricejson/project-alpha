# project-alpha

轻量级 project/ticket 管理工具。当前已实现后端基础框架、数据模型、标签 API、ticket API、ticket 搜索筛选与后端测试。

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
