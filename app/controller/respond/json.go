package respond

import (
	"encoding/json"
	"net/http"
)

// Data は正常系レスポンスとして v をそのまま JSON で返す。
func Data(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Error は仕様に従いエラー系レスポンス { "error": message } を返す。
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
