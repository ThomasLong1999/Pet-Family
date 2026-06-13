# Cat Manager

**[中文](README.md)**

A family-friendly pet management web app for tracking profiles, weight trends, and health records (deworming, vaccination, etc.) for cats, dogs, hamsters, rabbits, and more.

## Tech Stack

- **Backend**: Go + chi + SQLite (modernc.org/sqlite)
- **Frontend**: Vue 3 + TypeScript + Vite
- **Charts**: ECharts
- **Deployment**: Single binary (Go embed.FS)

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.22+ | Compile backend |
| [Node.js](https://nodejs.org/) | 18+ | Build frontend |

## Build

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/cat-manager.git
cd cat-manager
```

### 2. Build the frontend

```bash
cd frontend
npm install
npm run build
cd ..
```

### 3. Copy frontend assets to backend

```bash
# Linux / macOS
cp -r frontend/dist backend/cmd/server/static

# Windows (PowerShell)
Copy-Item -Recurse -Force frontend\dist backend\cmd\server\static

# Windows (CMD)
xcopy /E /I /Y frontend\dist backend\cmd\server\static
```

### 4. Compile the backend

```bash
cd backend
go build -o ../cat-manager ./cmd/server
cd ..
```

The output binary `cat-manager` (Linux/macOS) or `cat-manager.exe` (Windows) will be in the project root.

### Quick build (requires make)

```bash
make build
```

> Note: The Makefile uses Unix commands (`cp -r`). On Windows, run it inside Git Bash or MSYS2.

## Run

```bash
# Linux / macOS
./cat-manager

# Windows
.\cat-manager.exe
```

Visit `http://localhost:8080` after startup.

## Configure Listen Address

The listen IP and port can be set via command-line flags or environment variables.

### Command-line flags

```bash
# Custom port
./cat-manager -port 3000

# Custom IP and port
./cat-manager -host 192.168.1.100 -port 3000

# Localhost only
./cat-manager -host 127.0.0.1 -port 8080

# All interfaces (default)
./cat-manager -port 8080
```

### Environment variables

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

> Command-line flags take precedence over environment variables. If neither is set, the server listens on all interfaces (`0.0.0.0`) at port `8080`.

## Development

Start the frontend and backend separately:

```bash
# Terminal 1 — backend
cd backend && go run ./cmd/server

# Terminal 2 — frontend
cd frontend && npm run dev
```

Or use make to start both:

```bash
make dev
```

The frontend dev server runs at `http://localhost:5173` and proxies API requests to the backend automatically.

## Data Storage

- **Database**: `cat-manager.db` (SQLite, auto-created on first run)
- **Uploads**: `uploads/` directory (avatars, photos, etc.)

## Project Structure

```
cat-manager/
├── backend/
│   ├── cmd/server/          # Entry point & embedded frontend
│   ├── internal/
│   │   ├── handler/         # HTTP handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Data access layer
│   │   ├── model/           # Data models
│   │   └── middleware/      # Middleware
│   └── migrations/          # Database migrations
├── frontend/
│   ├── src/                 # Vue 3 source
│   └── public/              # Static assets
├── Makefile
└── README.md
```

## License

MIT
