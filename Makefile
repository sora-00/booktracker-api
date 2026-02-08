# BookTracker 開発用 Makefile

# デフォルトターゲット
.DEFAULT_GOAL := help

API_PORT := 8085

# ------------------------------------------
# コマンド一覧
# ------------------------------------------

## 🚀 すべて起動（Docker + エミュレータ）
## 事前に別ターミナルで: export PATH="/opt/homebrew/opt/openjdk@21/bin:$$PATH" && gcloud beta emulators datastore start --project=booktracker
dev:
	@echo "🚀 Starting API in Docker (connects to Datastore emulator on host)..."
	@echo "   Ensure emulator is running in another terminal: gcloud beta emulators datastore start --project=booktracker"
	docker compose up --build -d
	@echo "Docker is up. Following logs (Ctrl-C to detach)..."
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
