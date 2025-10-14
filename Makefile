# BookTracker 開発用 Makefile

# デフォルトターゲット
.DEFAULT_GOAL := help

# docker-compose.yml に合わせて設定
DB_SERVICE := db
API_PORT := 8085

# ------------------------------------------
# コマンド一覧
# ------------------------------------------

## 🚀 すべて起動（Docker + Go）
dev:
	@echo "🚀 Starting Docker containers and Go server..."
	docker compose up -d
	@echo "Docker containers are up. Following logs (Ctrl-C to detach)..."
	docker compose logs -f

## 🧠 ローカルのみGoサーバー起動（Dockerは使わない）
dev-local:
	@echo "🚀 Starting local Go server on :$(API_PORT) ..."
	PORT=$(API_PORT) go run main.go

## 🐳 Docker だけ起動
up:
	docker compose up -d

## 🧹 Docker 停止
down:
	docker compose down

## 🧠 Go サーバー起動
run:
	PORT=$(API_PORT) go run main.go

## 🔍 ログ確認
logs:
	docker compose logs -f
