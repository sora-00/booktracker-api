# BookTracker 開発用 Makefile

.DEFAULT_GOAL := help

ifneq (,$(wildcard .env))
include .env
export
endif

GCP_PROJECT_ID ?= booktracker
DATASTORE_PORT ?= 8059
DSADMIN_PORT ?= 8060
APP_INTERNAL_PORT ?= 8081
API_PORT ?= 8085

# ------------------------------------------
# コマンド一覧
# ------------------------------------------

## help: この一覧を表示
help:
	@grep -E '^## ' Makefile | sed 's/^## //'

## ⚙️ .env を .env.example から作成
env:
	@test -f .env || cp .env.example .env
	@echo ".env is ready"

## 🚀 local起動（datastore + dsadmin + api）
dev:
	@echo "🚀 Starting local stack (datastore + dsadmin + api)..."
	DATASTORE_EMULATOR_HOST=datastore:$(DATASTORE_PORT) docker compose --profile local up --build -d
	@echo "Datastore emulator: localhost:$(DATASTORE_PORT)  dsadmin: http://localhost:$(DSADMIN_PORT)  API: http://localhost:$(API_PORT)"
	@echo "Following logs (Ctrl-C to detach)..."
	docker compose --profile local logs -f

## 🚀 prod想定起動（apiのみ・本番 Datastore）
prod:
	@echo "🚀 Starting api (prod profile, real Cloud Datastore)..."
	DATASTORE_EMULATOR_HOST= docker compose --profile prod up --build -d api
	@echo "API: http://localhost:$(API_PORT)"

## 🚀 本番 Cloud Datastore 直結（Docker: ADC は ~/.config/gcloud）
dev-cloud:
	@echo "🚀 Starting api (cloud profile, real Cloud Datastore)..."
	@echo "    Ensure: gcloud auth application-default login"
	DATASTORE_EMULATOR_HOST= docker compose --profile cloud up --build -d api
	@echo "API: http://localhost:$(API_PORT)"
	DATASTORE_EMULATOR_HOST= docker compose --profile cloud logs -f api

## 🧠 本番 Cloud Datastore 直結（ホストで go run。ADC のみ。DATASTORE_EMULATOR_HOST は付けない）
dev-cloud-local:
	@echo "🚀 Starting Go on :$(API_PORT) → real Cloud Datastore (unset DATASTORE_EMULATOR_HOST)..."
	@echo "    Ensure: gcloud auth application-default login"
	PORT=$(API_PORT) GCP_PROJECT_ID=$(GCP_PROJECT_ID) env -u DATASTORE_EMULATOR_HOST go run main.go

## 🧠 ローカルのみ Go サーバー起動（エミュレータは Docker の datastore サービスを先に起動）
dev-local: ensure-datastore-docker
	@echo "🚀 Starting local Go server on :$(API_PORT) (DATASTORE_EMULATOR_HOST=localhost:$(DATASTORE_PORT))..."
	PORT=$(API_PORT) DATASTORE_EMULATOR_HOST=localhost:$(DATASTORE_PORT) GCP_PROJECT_ID=$(GCP_PROJECT_ID) go run main.go

## 🗄 Datastore エミュレータ（Docker）起動保証
ensure-datastore-docker:
	@echo "🚀 Ensuring datastore container is up..."
	docker compose --profile local up -d datastore
	@sleep 2

## 🐳 local全サービス起動（ビルドなし）
up:
	DATASTORE_EMULATOR_HOST=datastore:$(DATASTORE_PORT) docker compose --profile local up -d

## 🧹 Docker 停止（ボリュームは残す）
down:
	docker compose --profile local --profile prod --profile cloud down

## 🧪 テスト（domain / usecase・request のみ、race 付き）
test:
	go test -race ./app/domain/... ./app/usecase/...

## 🧠 Go サーバー起動（環境変数は自分で設定）
run:
	PORT=$(API_PORT) go run main.go

## 🔍 ログ確認
logs:
	docker compose --profile local logs -f
