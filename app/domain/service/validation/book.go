package validation

import "errors"

const (
	MaxLenTitle      = 30
	MaxLenAuthor     = 30
	MaxLenPublish    = 20
	MaxLenBackground = 200
	MaxPagesDigits   = 9999
)

// ValidateBookStatus は status が unread / reading / read のいずれかであることを検証する。
func ValidateBookStatus(s string) error {
	if s != "unread" && s != "reading" && s != "read" {
		return errors.New("status must be unread, reading, or read")
	}
	return nil
}

// ValidatePagesAll は pagesAll の範囲（1 < n <= 9999）を検証する。
func ValidatePagesAll(n int) error {
	if n < 1 {
		return errors.New("pagesAll must be greater than 1")
	}
	if n > MaxPagesDigits {
		return errors.New("pagesAll must be 4 digits or less")
	}
	return nil
}

// ValidateTargetReadPagesPerDay は targetReadPagesPerDay の範囲と pagesAll との関係を検証する。
func ValidateTargetReadPagesPerDay(pagesPerDay, pagesAll int) error {
	if pagesPerDay <= 0 {
		return errors.New("targetReadPagesPerDay must be greater than 0")
	}
	if pagesPerDay >= pagesAll {
		return errors.New("targetReadPagesPerDay must be less than pagesAll")
	}
	if pagesPerDay > MaxPagesDigits {
		return errors.New("targetReadPagesPerDay must be 4 digits or less")
	}
	return nil
}

