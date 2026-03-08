package service

import "time"

// RemainingDays は目標日までの残り日数（0以上）。過ぎている場合は0。
func RemainingDays(target time.Time) int {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	targetDay := target.UTC().Truncate(24 * time.Hour)
	diff := targetDay.Sub(now)
	days := int(diff / (24 * time.Hour))
	if days < 0 {
		return 0
	}
	return days
}
