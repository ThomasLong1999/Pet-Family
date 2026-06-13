.PHONY: dev dev-backend dev-frontend build build-frontend clean

# Development mode: run frontend and backend concurrently
dev:
	@make -j2 dev-backend dev-frontend

dev-backend:
	@cd backend && go run ./cmd/server

dev-frontend:
	@cd frontend && npx vite --port 5173

# Build: compile frontend, copy to backend static, build single binary
build: build-frontend
	@echo "Building Go binary..."
	@rm -rf backend/cmd/server/static
	@cp -r frontend/dist backend/cmd/server/static
	@cd backend && go build -o ../cat-manager ./cmd/server
	@echo "✅ Built cat-manager"

build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo "✅ Frontend built"

clean:
	@rm -rf frontend/dist frontend/node_modules backend/cmd/server/static cat-manager cat-manager.exe
	@echo "✅ Cleaned"
