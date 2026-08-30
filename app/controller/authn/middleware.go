package authn

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sora-00/booktracker-api/app/controller/respond"
	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/infra/auth"
	httpx "github.com/sora-00/booktracker-api/app/infra/http"
)

var errUnauthorized = errors.New("unauthorized")

// RequireAuth は /api 配下で認証を必須にする middleware。
// - Authorization: Bearer <Firebase ID Token> があればそれを検証して UID を採用
// - ない場合は X-User-ID を開発用として許可
// また、初回ログイン（Me が存在しない）なら Me を作成する。
func RequireAuth(verifier *auth.FirebaseVerifier, meRepo repository.MeRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := resolveUserID(r, verifier)
			if err != nil {
				respond.Error(w, errToStatus(err), err.Error())
				return
			}

			ctx := auth.WithUserID(r.Context(), userID)
			if err := ensureMe(ctx, meRepo, userID); err != nil {
				respond.Error(w, http.StatusInternalServerError, err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func errToStatus(err error) int {
	if errors.Is(err, errUnauthorized) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func resolveUserID(r *http.Request, verifier *auth.FirebaseVerifier) (string, error) {
	if uid, ok := userIDFromBearerToken(r, verifier); ok {
		return uid, nil
	}
	if uid, ok := userIDFromDevHeader(r); ok {
		return uid, nil
	}
	return "", errUnauthorized
}

func userIDFromBearerToken(r *http.Request, verifier *auth.FirebaseVerifier) (string, bool) {
	if verifier == nil {
		return "", false
	}
	authz := strings.TrimSpace(r.Header.Get("Authorization"))
	if authz == "" {
		return "", false
	}
	uid, err := verifier.VerifyBearerToken(r.Context(), authz)
	if err != nil {
		return "", false
	}
	return uid, true
}

func userIDFromDevHeader(r *http.Request) (string, bool) {
	uid := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if uid == "" {
		return "", false
	}
	return uid, true
}

func ensureMe(ctx context.Context, meRepo repository.MeRepo, userID string) error {
	if meRepo == nil {
		return nil
	}
	_, err := meRepo.Get(ctx, userID)
	if err == nil {
		return nil
	}
	if !httpx.IsNotFound(err) {
		return err
	}
	m, err := entity.NewMe(userID)
	if err != nil {
		return err
	}
	return meRepo.Put(ctx, m)
}
