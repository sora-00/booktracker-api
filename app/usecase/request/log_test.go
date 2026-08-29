package request

import (
	"strings"
	"testing"
)

func TestNewLogGetByID(t *testing.T) {
	t.Parallel()
	req := newRequestWithChiRoute("GET", "/logs/10", map[string]string{"id": "10"})
	got, err := NewLogGetByID(req)
	if err != nil {
		t.Fatal(err)
	}
	if got.LogID != 10 {
		t.Fatalf("LogID = %d", got.LogID)
	}
}

func TestNewLogGetByBookID(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		params  map[string]string
		wantID  int
		wantErr string
	}{
		{"ok", map[string]string{"bookId": "99"}, 99, ""},
		{"empty", map[string]string{"bookId": ""}, 0, "book id is required"},
		{"bad", map[string]string{"bookId": "nope"}, 0, "invalid book id"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := newRequestWithChiRoute("GET", "/books/1/logs", tc.params)
			got, err := NewLogGetByBookID(req)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.BookID != tc.wantID {
				t.Fatalf("BookID = %d, want %d", got.BookID, tc.wantID)
			}
		})
	}
}

func TestNewLogDeleteByBookID(t *testing.T) {
	t.Parallel()
	req := newRequestWithChiRoute("DELETE", "/books/3/logs", map[string]string{"bookId": "3"})
	got, err := NewLogDeleteByBookID(req)
	if err != nil {
		t.Fatal(err)
	}
	if got.BookID != 3 {
		t.Fatalf("BookID = %d", got.BookID)
	}
}
