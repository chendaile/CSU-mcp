```
┏━╸┏━┓╻ ╻   ┏┳┓┏━╸┏━┓
┃  ┗━┓┃ ┃   ┃┃┃┃  ┣━┛
┗━╸┗━┛┗━┛   ╹ ╹┗━╸╹                
   by oftheloneliness · github.com/chendaile
```

## CSU MCP 本地服务

中南大学校园信息聚合 + MCP 代理，一条 Docker 命令即可在本地启动，AI 助手马上能查成绩、看课表、问校车、看招聘。

---

### 1. 能力速览
- **REST API**：统一 token 校验，提供 `/api/v1/jwc/:id/:pwd/grade`、`/rank`、`/class/:term/:week`、`/bus/search/:start/:end/:day`、`/job/...` 等接口。
- **MCP 工具集**：开放 `csu.grade`、`csu.rank`、`csu.classes`、`csu.bus_search`、`csu.jobs`，Claude/Kilocode/Continue/Open Interpreter 等 MCP 客户端可直接调用。
- **本地自用**：镜像内置 landing 页面与日志，适合个人在笔记本、台式机或实验服务器上自建，数据只在本机流转。

---

### 2. Docker 快速启动

> 请先安装 Docker（macOS/Win 用 Docker Desktop，Linux 可 `apt install docker-ce`）。

#### 2.1 直接 `docker run`
```bash
docker pull ghcr.io/chendaile/csu-mcp:latest

docker run -d \
  --name csu-mcp \
  -p 12000:12000 \
  -p 13000:13000 \
  ghcr.io/chendaile/csu-mcp:latest
```
访问地址：
- API 首页：<http://localhost:12000/>
- MCP 入口：<http://localhost:13000/>
- API 调用示例：
  ```bash
  curl "http://localhost:12000/api/v1/jwc/学号/密码/grade?token=csugo-token"
  ```

**自定义端口或 token(一般默认即可)**
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

#### 2.2 使用 `docker compose`
项目已经自带 `docker-compose.yml`，在仓库根目录执行：
```bash
git clone https://github.com/chendaile/CSU-mcp.git
cd CSU-mcp
docker compose up -d         # 自动 build & run
```
可通过 `.env` 文件覆盖以下变量：

| 变量 | 默认值 | 说明 |
| ---- | ------ | ---- |
| `CSU_API_PORT` | 12000 | 映射到宿主机的 API 端口 |
| `CSU_MCP_PORT` | 13000 | 映射到宿主机的 MCP 端口 |
| `CSUGO_HTTP_PORT` | 12000 | 容器内 API 监听端口 |
| `CSUGO_BASE_URL` | `http://127.0.0.1:12000` | MCP 访问 API 的地址 |
| `CSUGO_TOKEN` | `csugo-token` | API token |
| `MCP_HTTP_ADDR` | `:13000` | 容器内 MCP 监听地址 |


---

### 3. AI 助手接入（MCP 案例）

#### Claude Desktop
```bash
claude mcp add --transport http csu-mcp http://localhost:13000
```
重启 Claude，工具栏会出现 `csu.*` 系列按钮。让 Claude 执行“查询本学期课表”即可自动调用 MCP 工具。

#### Kilocode、Continue、Cursor 等 IDE
```bash
mcp-client add csu-mcp http://localhost:13000
```
或在图形界面里填写：
- Name: `csu-mcp`
- Endpoint: `http://localhost:13000`

#### Open Interpreter ≥ 0.3
```bash
oi mcp add csu-mcp http://localhost:13000
```
之后在对话中直接请求“帮我查成绩/校车/招聘”即可。

---

如需更多玩法或反馈建议，欢迎访问 [github.com/chendaile/CSU-mcp](https://github.com/chendaile/CSU-mcp) 提 Issue。祝使用愉快！
