package request

import (
	"strings"
	"testing"
)

func TestNewBookGet(t *testing.T) {
	req := newRequestWithChiRoute("GET", "/books", nil)
	got, err := NewBookGet(req)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected non-nil")
	}
}

func TestNewBookGetByID(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		params  map[string]string
		wantID  int
		wantErr string
	}{
		{"ok", map[string]string{"id": "42"}, 42, ""},
		{"empty id", map[string]string{"id": ""}, 0, "book id is required"},
		{"missing param", nil, 0, "book id is required"},
		{"not a number", map[string]string{"id": "x"}, 0, "invalid book id"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := newRequestWithChiRoute("GET", "/books/42", tc.params)
			got, err := NewBookGetByID(req)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("err = %v, want %q", err, tc.wantErr)
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

func TestNewBookGetByStatus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		params  map[string]string
		want    string
		wantErr string
	}{
		{"reading", map[string]string{"status": "reading"}, "reading", ""},
		{"unread", map[string]string{"status": "unread"}, "unread", ""},
		{"read", map[string]string{"status": "read"}, "read", ""},
		{"empty", map[string]string{"status": ""}, "", "status is required"},
		{"invalid", map[string]string{"status": "done"}, "", "status must be unread, reading, or read"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := newRequestWithChiRoute("GET", "/books/status/reading", tc.params)
			got, err := NewBookGetByStatus(req)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.Status != tc.want {
				t.Fatalf("Status = %q, want %q", got.Status, tc.want)
			}
		})
	}
}

func TestNewBookDelete(t *testing.T) {
	t.Parallel()
	req := newRequestWithChiRoute("DELETE", "/books/1", map[string]string{"id": "7"})
	got, err := NewBookDelete(req)
	if err != nil {
		t.Fatal(err)
	}
	if got.BookID != 7 {
		t.Fatalf("BookID = %d, want 7", got.BookID)
	}
}
