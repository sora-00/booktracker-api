package httpx

import (
	"errors"
	"net/http"

	domainerr "github.com/sora-00/booktracker-api/app/domain/errors"
)

// IsNotFound は repository 層の ErrNotFound を判定する。
func IsNotFound(err error) bool {
	return errors.Is(err, domainerr.ErrNotFound)
}

// StatusAndMessage は err を HTTP ステータスとレスポンスメッセージに変換する。
func StatusAndMessage(err error, notFoundMessage string) (int, string) {
	if IsNotFound(err) {
		return http.StatusNotFound, notFoundMessage
	}
	return http.StatusInternalServerError, err.Error()
}
