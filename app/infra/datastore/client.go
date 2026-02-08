package datastore

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
)

type contextKey struct{}

// WithContext は context に Datastore クライアントを入れる。middleware などでリクエストごとに呼ぶ。
func WithContext(ctx context.Context, client *datastore.Client) context.Context {
	return context.WithValue(ctx, contextKey{}, client)
}

// FromContext は context から Datastore クライアントを取得する。repository 層で利用。
func FromContext(ctx context.Context) (*datastore.Client, bool) {
	client, ok := ctx.Value(contextKey{}).(*datastore.Client)
	return client, ok
}

// NewClient は GCP Cloud Datastore のクライアントを返す。
// プロジェクトIDは環境変数 GCP_PROJECT_ID または GOOGLE_CLOUD_PROJECT で指定。
// ローカルでは DATASTORE_EMULATOR_HOST=localhost:8081 でエミュレータに接続できる。
func NewClient(ctx context.Context) (*datastore.Client, error) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if projectID == "" {
		projectID = "booktracker" // エミュレータ用のダミー
	}
	return datastore.NewClient(ctx, projectID)
}
