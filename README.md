# Cat Manager - 猫咪管理工具

**[English](README_EN.md)**

家庭共享的宠物管理 Web 应用，用于记录猫、狗、仓鼠、兔子等宠物的基本档案、体重变化和驱虫/疫苗等健康记录。

## 技术栈

- **后端**: Go + chi + SQLite (modernc.org/sqlite)
- **前端**: Vue 3 + TypeScript + Vite
- **图表**: ECharts
- **部署**: 单一二进制文件 (Go embed.FS)

## 前置要求

| 工具 | 版本 | 用途 |
|------|------|------|
| [Go](https://go.dev/dl/) | 1.22+ | 编译后端 |
| [Node.js](https://nodejs.org/) | 18+ | 构建前端 |

## 构建

### 1. 克隆项目

```bash
git clone https://github.com/<your-username>/cat-manager.git
cd cat-manager
```

### 2. 构建前端

```bash
cd frontend
npm install
npm run build
cd ..
```

### 3. 将前端产物复制到后端

```bash
# Linux / macOS
cp -r frontend/dist backend/cmd/server/static

# Windows (PowerShell)
Copy-Item -Recurse -Force frontend\dist backend\cmd\server\static

# Windows (CMD)
xcopy /E /I /Y frontend\dist backend\cmd\server\static
```

### 4. 编译后端

```bash
cd backend
go build -o ../cat-manager ./cmd/server
cd ..
```

编译产物 `cat-manager`（Linux/macOS）或 `cat-manager.exe`（Windows）位于项目根目录。

### 一键构建（需要 make）

```bash
make build
```

> 注意：Makefile 中的 `cp -r` 在 Windows 下需要 Git Bash 或 MSYS2 环境。

## 运行

```bash
# Linux / macOS
./cat-manager

# Windows
.\cat-manager.exe
```

启动后访问 `http://localhost:8080`。

## 配置监听地址

支持通过命令行参数或环境变量指定监听 IP 和端口。

### 命令行参数

```bash
# 监听指定端口
./cat-manager -port 3000

# 监听指定 IP 和端口
./cat-manager -host 192.168.1.100 -port 3000

# 仅监听本机
./cat-manager -host 127.0.0.1 -port 8080

# 监听所有接口（默认行为）
./cat-manager -port 8080
```

### 环境变量

```bash
# Linux / macOS
export HOST=192.168.1.100
export PORT=3000
./cat-manager

# Windows (PowerShell)
$env:HOST="192.168.1.100"
$env:PORT="3000"
.\cat-manager.exe

# Windows (CMD)
set HOST=192.168.1.100
set PORT=3000
cat-manager.exe
```

> 命令行参数优先级高于环境变量；两者均未设置时默认监听所有接口的 `8080` 端口。

## 开发模式

需要分别启动前端和后端：

```bash
# 终端 1 - 后端
cd backend && go run ./cmd/server

# 终端 2 - 前端
cd frontend && npm run dev
```

或使用 make 一键启动：

```bash
make dev
```

前端开发服务器运行在 `http://localhost:5173`，会自动代理 API 请求到后端。

## 数据存储

- **数据库**: `cat-manager.db`（SQLite，首次运行自动创建）
- **上传文件**: `uploads/` 目录（头像、照片等）

## 项目结构

```
cat-manager/
├── backend/
│   ├── cmd/server/          # 入口 & 前端静态资源
│   ├── internal/
│   │   ├── handler/         # HTTP 处理层
│   │   ├── service/         # 业务逻辑层
│   │   ├── repository/      # 数据访问层
│   │   ├── model/           # 数据模型
│   │   └── middleware/      # 中间件
│   └── migrations/          # 数据库迁移
├── frontend/
│   ├── src/                 # Vue 3 源码
│   └── public/              # 静态资源
├── Makefile
└── README.md
```

## License

MIT
