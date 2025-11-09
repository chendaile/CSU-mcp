## CSU MCP 本地服务

一个封装好的中南大学校园信息聚合服务，内置 REST API 与 Model Context Protocol (MCP) 代理，只需启动一个 Docker 容器就能让 AI 助手帮你查询成绩、课表、校车、招聘等数据。

---

### 1. 这个 MCP 可以做什么？
- **API 能力**：提供 `/api/v1/jwc/:id/:pwd/grade`、`/rank`、`/class/:term/:week`、`/bus/search/:start/:end/:day`、`/job/...` 等端点，统一 token 验证，方便自建前端或小程序调用。
- **MCP 工具**：开放 `csu.grade`、`csu.rank`、`csu.classes`、`csu.bus_search`、`csu.jobs` 五个工具，Claude、Kilocode、Continue、Open Interpreter 等支持 MCP 的助手都能直接调用。
- **单机部署**：无需公网或复杂配置，适合本地自用；作者信息和 landing 页面已内置，可直接展示在 `http://localhost`。

---

### 2. 用 Docker 一键启动
> 先安装 Docker（macOS/Windows 用 Docker Desktop，Linux 用 `apt install docker-ce`）。

#### 拉取镜像
```bash
docker pull ghcr.io/chendaile/csu-mcp:latest
```

#### 快速运行（默认 API 端口 12000，MCP 端口 13000）
```bash
docker run -d \
  --name csu-mcp \
  -p 12000:12000 \
  -p 13000:13000 \
  ghcr.io/chendaile/csu-mcp:latest
```

完成后：
- 浏览器访问 `http://localhost:12000/` 查看 API 首页及路由说明；
- 访问 `http://localhost:13000/` 查看 MCP 入口；
- 测试接口：
  ```bash
  curl "http://localhost:12000/api/v1/jwc/学号/密码/grade?token=csugo-token"
  ```

#### 自定义端口或 token（可选）
```bash
docker run -d --name csu-mcp \
  -e CSUGO_HTTP_PORT=8080 \
  -e CSUGO_BASE_URL=http://127.0.0.1:8080 \
  -e CSUGO_TOKEN=my-token \
  -e MCP_HTTP_ADDR=:18080 \
  -p 8080:8080 \
  -p 18080:18080 \
  ghcr.io/chendaile/csu-mcp:latest
```
常用环境变量：

| 变量 | 默认值 | 说明 |
| ---- | ------ | ---- |
| `CSUGO_HTTP_PORT` | 12000 | 容器内 API 监听端口 |
| `CSUGO_BASE_URL` | `http://127.0.0.1:12000` | MCP 访问 API 的地址 |
| `CSUGO_TOKEN` | `csugo-token` | API 调用需附带的 token |
| `MCP_HTTP_ADDR` | `:13000` | MCP 监听地址 |

> 若宿主机需要科学上网/校园网代理，请在 `docker run` 时加 `-e http_proxy -e https_proxy`，确保容器也能访问教务系统。

---

### 3. 常见 AI 助手的 MCP 接入示例

#### Claude Desktop
```bash
claude mcp add csu-mcp http://localhost:13000
```
重启 Claude，工具面板会出现 `csu.grade` 等按钮，直接让 Claude 查询成绩/课表即可。

#### Kilocode / Continue / 其他 MCP 客户端
大多数 MCP 客户端使用同样格式：
```bash
mcp-client add csu-mcp http://localhost:13000
```
或在图形界面中填写：
- 名称：`csu-mcp`
- Endpoint：`http://localhost:13000`

#### Open Interpreter 0.3+
```bash
oi mcp add csu-mcp http://localhost:13000
```
之后在对话里输入“调用 csu.grade 查询我的成绩”即可触发工具调用。

---

如有疑问或建议，欢迎访问 [项目仓库](https://github.com/chendaile/CSU-mcp) 提 Issue。祝使用愉快！
