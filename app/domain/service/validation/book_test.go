package validation

import "testing"

func TestValidateBookStatus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want bool
	}{
		{"unread", true},
		{"reading", true},
		{"read", true},
		{"", false},
		{"done", false},
		{"UNREAD", false},
	}
	for _, tc := range cases {
		if ok := ValidateBookStatus(tc.in) == nil; ok != tc.want {
			t.Errorf("ValidateBookStatus(%q) ok=%v, want %v", tc.in, ok, tc.want)
		}
	}
}

func TestValidatePagesAll(t *testing.T) {
	t.Parallel()
	if err := ValidatePagesAll(100); err != nil {
		t.Error(err)
	}
	if err := ValidatePagesAll(0); err == nil {
		t.Fatal("expected error for 0")
	}
	if err := ValidatePagesAll(10000); err == nil {
		t.Fatal("expected error for > MaxPagesDigits")
	}
}

func TestValidateTargetReadPagesPerDay(t *testing.T) {
	t.Parallel()
	if err := ValidateTargetReadPagesPerDay(10, 100); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTargetReadPagesPerDay(0, 100); err == nil {
		t.Fatal("expected error")
	}
	if err := ValidateTargetReadPagesPerDay(100, 100); err == nil {
		t.Fatal("expected error when >= pagesAll")
	}
}
