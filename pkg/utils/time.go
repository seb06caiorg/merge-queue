package utils

import (
	"fmt"
	"time"
)

// TimeUtils provides utility functions for time operations.
type TimeUtils struct{}

// NewTimeUtils creates a new TimeUtils instance.
func NewTimeUtils() *TimeUtils {
	return &TimeUtils{}
}

// FormatDuration returns a human-readable duration string.
func (tu *TimeUtils) FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

// IsToday checks if a time is today.
func (tu *TimeUtils) IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsThisWeek checks if a time is within the current week.
func (tu *TimeUtils) IsThisWeek(t time.Time) bool {
	now := time.Now()
	year, week := now.ISOWeek()
	tYear, tWeek := t.ISOWeek()
	return year == tYear && week == tWeek
}

// StartOfDay returns the start of the day for the given time.
func (tu *TimeUtils) StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for the given time.
func (tu *TimeUtils) EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// DaysBetween calculates the number of days between two times.
func (tu *TimeUtils) DaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	startDay := tu.StartOfDay(start)
	endDay := tu.StartOfDay(end)

	return int(endDay.Sub(startDay).Hours() / 24)
}

// FormatRelativeTime returns a human-readable relative time string.
func (tu *TimeUtils) FormatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 0 {
		// Future time.
		diff = -diff
		if diff < time.Minute {
			return "in a few seconds"
		}
		if diff < time.Hour {
			minutes := int(diff.Minutes())
			if minutes == 1 {
				return "in 1 minute"
			}
			return fmt.Sprintf("in %d minutes", minutes)
		}
		if diff < 24*time.Hour {
			hours := int(diff.Hours())
			if hours == 1 {
				return "in 1 hour"
			}
			return fmt.Sprintf("in %d hours", hours)
		}
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "tomorrow"
		}
		return fmt.Sprintf("in %d days", days)
	}

	// Past time.
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(diff.Hours() / 24)
	if days == 1 {
		return "yesterday"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}
	weeks := days / 7
	if weeks == 1 {
		return "1 week ago"
	}
	if weeks < 4 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	return t.Format("Jan 2, 2006")
}
