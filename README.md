# CSU MCP

一个聚合中南大学统一认证数据的 REST + MCP 服务，可查成绩、课表、校车、招聘，并供 Claude / Kilocode / VS Code 等客户端直接调用 `csu.*` 工具。

---

## 1. 用处
- **REST API**：提供 `/api/v1/jwc/:id/:pwd/grade`、`/rank`、`/class/:term/:week`、`/bus/search/:start/:end/:day`、`/job/...` 等统一接口。
- **MCP 工具**：自动开放 `csu.grade`、`csu.rank`、`csu.classes`、`csu.bus_search`、`csu.jobs`，AI 助手无需额外集成。
- **单镜像部署**：API、MCP、Landing 全打包，适合自托管或直接复用公共服务。

---

## 2. 直接连接公共服务器

公共 MCP 地址：`http://8.162.3.110:13000`

| 客户端 | 配置示例 |
| ------ | -------- |
| Kilocode | `"csu": { "type": "streamable-http", "url": "http://8.162.3.110:13000" }` |
| Claude Desktop | `claude mcp add --transport http csu http://8.162.3.110:13000` |
| VS Code MCP 扩展 | `"csu": { "type": "http", "url": "http://8.162.3.110:13000" }` |

保存后即可在对应客户端直接调用 `csu.*` 工具，无需本地部署。

---

## 3. 自己拉镜像运行

镜像在 GHCR 与阿里云同步，可任选其一：

```bash
docker pull ghcr.io/chendaile/csu-mcp:latest
# 或
docker pull crpi-qpej8ufiacto12s1.cn-hangzhou.personal.cr.aliyuncs.com/csu-mcp/csu-mcp:<版本号>
```

运行：

```bash
docker run -d \
  --name csu-mcp \
  -p 12000:12000 \
  -p 13000:13000 \
  <TARGET_IMAGE>
```

启动后：
- `http://localhost:12000/` —— REST API 首页，可直接带 `token=csugo-token` 调接口。
- `http://localhost:13000/` —— MCP 入口，本地 MCP 客户端指向此地址即可。

若需自定义端口或 token，复制 `configs/api/conf/app.conf` 修改后，通过 `-v /path/app.conf:/app/configs/api/conf/app.conf:ro` 挂载覆盖。欢迎 Issue / PR。***
