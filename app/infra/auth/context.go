package auth

import "context"

type userIDKey struct{}

// WithUserID は context に認証ユーザーID（Firebase UID）を格納する。middleware で使用。
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFromContext は context からユーザーIDを取得する。未設定の場合は空文字。
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey{}).(string)
	return v
}
