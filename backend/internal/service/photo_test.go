package service

import (
	"testing"
	"time"
)

// dateMonthsAgo returns a YYYY-MM-DD string for ~monthsAgo months before now.
func dateMonthsAgo(monthsAgo int) string {
	t := time.Now().AddDate(0, -monthsAgo, 0)
	return t.Format("2006-01-02")
}

func TestCalculateAgeGroup_InvalidDate(t *testing.T) {
	if got := CalculateAgeGroup("not-a-date"); got != "1y" {
		t.Fatalf("expected 1y for invalid date, got %q", got)
	}
	if got := CalculateAgeGroup(""); got != "1y" {
		t.Fatalf("expected 1y for empty, got %q", got)
	}
}

func TestCalculateAgeGroup_Newborn(t *testing.T) {
	// Future date → months <= 0 → "1m"
	got := CalculateAgeGroup("2999-01-01")
	if got != "1m" {
		t.Fatalf("expected 1m for newborn, got %q", got)
	}
}

func TestCalculateAgeGroup_Months(t *testing.T) {
	for _, m := range []int{1, 6, 11} {
		got := CalculateAgeGroup(dateMonthsAgo(m))
		// Format is "<n>m"
		if len(got) < 2 || got[len(got)-1] != 'm' {
			t.Fatalf("expected *m for %d months ago, got %q", m, got)
		}
	}
}

func TestCalculateAgeGroup_YearCappedAt20(t *testing.T) {
	got := CalculateAgeGroup("1990-01-01")
	if got != "20y" {
		t.Fatalf("expected 20y for very old pet, got %q", got)
	}
}

func TestCalculateAgeGroup_Years(t *testing.T) {
	// ~30 months ago → 2y
	got := CalculateAgeGroup(dateMonthsAgo(30))
	if got != "2y" {
		t.Fatalf("expected 2y for ~30 months ago, got %q", got)
	}
}
