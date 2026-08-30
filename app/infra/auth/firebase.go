package auth

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// FirebaseVerifier は Firebase ID トークン検証を行う。
type FirebaseVerifier struct {
	client *auth.Client
}

// NewFirebaseVerifier は Firebase Admin SDK を初期化する。
// credentialsJSONPath が空の場合は Application Default Credentials を使う。
func NewFirebaseVerifier(ctx context.Context, credentialsJSONPath string) (*FirebaseVerifier, error) {
	var app *firebase.App
	var err error
	if strings.TrimSpace(credentialsJSONPath) != "" {
		app, err = firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsJSONPath))
	} else {
		app, err = firebase.NewApp(ctx, nil)
	}
	if err != nil {
		return nil, err
	}
	c, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseVerifier{client: c}, nil
}

// VerifyBearerToken は Authorization: Bearer <token> を検証して UID を返す。
func (v *FirebaseVerifier) VerifyBearerToken(ctx context.Context, authorizationHeader string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorizationHeader, prefix) {
		return "", errors.New("missing bearer token")
	}
	raw := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, prefix))
	if raw == "" {
		return "", errors.New("missing bearer token")
	}
	tok, err := v.client.VerifyIDToken(ctx, raw)
	if err != nil {
		return "", err
	}
	return tok.UID, nil
}

