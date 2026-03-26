.PHONY: build run clean test

# 构建后端
build:
	cd apps/server && go build -o bin/server ./cmd/server

# 运行后端
run:
	cd apps/server && go run cmd/server/main.go

# 运行数据库迁移
migrate-up:
	golang-migrate -path apps/server/db/migrations -database "postgres://postgres:Admin_2026@localhost:5432/cloudim?sslmode=disable" up

migrate-down:
	golang-migrate -path apps/server/db/migrations -database "postgres://postgres:Admin_2026@localhost:5432/cloudim?sslmode=disable" down

# 测试
test:
	cd apps/server && go test ./...

# 清理
clean:
	rm -rf apps/server/bin
	rm -rf apps/client/dist
	rm -rf apps/client/dist-electron

# 前端开发
frontend-dev:
	cd apps/client && npm run dev

# 前端构建
frontend-build:
	cd apps/client && npm run build

# Electron 开发
electron-dev:
	cd apps/client && npm run electron:dev

# Electron 打包
electron-build:
	cd apps/client && npm run electron:build

# 帮助
help:
	@echo "CloudIM Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build          - 构建后端"
	@echo "  run            - 运行后端"
	@echo "  migrate-up     - 执行数据库迁移"
	@echo "  migrate-down   - 回滚数据库迁移"
	@echo "  test           - 运行测试"
	@echo "  clean          - 清理构建产物"
	@echo "  frontend-dev   - 启动前端开发服务器"
	@echo "  frontend-build - 构建前端"
	@echo "  electron-dev   - 启动 Electron 开发模式"
	@echo "  electron-build - 打包 Electron 应用"
