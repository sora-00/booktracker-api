package request

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

// newRequestWithChiRoute は chi.URLParam が効く *http.Request を返す（テスト用）。
func newRequestWithChiRoute(method, path string, routeParams map[string]string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	rctx := chi.NewRouteContext()
	for k, v := range routeParams {
		rctx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}
