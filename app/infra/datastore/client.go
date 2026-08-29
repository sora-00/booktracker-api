package datastore

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/datastore"
)

// ErrNoClient は WithContext でクライアントが載っていない（または nil）ときに FromContext が返す。
var ErrNoClient = errors.New("datastore client not found in context")

type contextKey struct{}

// WithContext は context に Datastore クライアントを入れる。middleware などでリクエストごとに呼ぶ。
func WithContext(ctx context.Context, client *datastore.Client) context.Context {
	return context.WithValue(ctx, contextKey{}, client)
}

// FromContext は context から Datastore クライアントを取得する。
// Paircare の FromContext（単一返り値 + error）に近い形。ミドルウェア未設定時は ErrNoClient。
//
// 注: cloud.google.com/go/datastore では boom のような遅延 New は行わない。
// 接続は main で1回 NewClient し、WithContext でリクエストに載せる想定（テストは WithContext で注入）。
func FromContext(ctx context.Context) (*datastore.Client, error) {
	client, ok := ctx.Value(contextKey{}).(*datastore.Client)
	if !ok || client == nil {
		return nil, ErrNoClient
	}
	return client, nil
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
