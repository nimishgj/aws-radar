package collector

import (
	"testing"
	"time"
)

func TestCostCollectorShouldRefreshDaily(t *testing.T) {
	c := NewCostCollector("daily")
	now := time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC)

	if !c.shouldRefresh(now) {
		t.Fatalf("expected refresh when no previous fetch exists")
	}

	c.lastFetch = now.Add(-2 * time.Hour)
	if c.shouldRefresh(now) {
		t.Fatalf("did not expect refresh on same UTC date")
	}

	c.lastFetch = now.Add(-26 * time.Hour)
	if !c.shouldRefresh(now) {
		t.Fatalf("expected refresh on next UTC date")
	}
}

func TestCostCollectorInvalidFrequencyFallsBackToDaily(t *testing.T) {
	c := NewCostCollector("weird")
	if c.frequency != "daily" {
		t.Fatalf("expected invalid frequency to fall back to daily, got %s", c.frequency)
	}
	now := time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC)
	c.lastFetch = now.Add(-2 * time.Hour)

	if c.shouldRefresh(now) {
		t.Fatalf("expected daily fallback to avoid refresh on same date")
	}
}

func TestCostCollectorShouldRefreshHourlyWeeklyMonthly(t *testing.T) {
	now := time.Date(2026, 2, 28, 12, 45, 0, 0, time.UTC)

	hourly := NewCostCollector("hourly")
	hourly.lastFetch = now.Add(-30 * time.Minute)
	if hourly.shouldRefresh(now) {
		t.Fatalf("expected hourly collector not to refresh within same hour")
	}
	hourly.lastFetch = now.Add(-90 * time.Minute)
	if !hourly.shouldRefresh(now) {
		t.Fatalf("expected hourly collector to refresh on next hour")
	}

	weekly := NewCostCollector("weekly")
	weekly.lastFetch = now.AddDate(0, 0, -3)
	if weekly.shouldRefresh(now) {
		t.Fatalf("expected weekly collector not to refresh in same ISO week")
	}
	weekly.lastFetch = now.AddDate(0, 0, -8)
	if !weekly.shouldRefresh(now) {
		t.Fatalf("expected weekly collector to refresh in new ISO week")
	}

	monthly := NewCostCollector("monthly")
	monthly.lastFetch = now.AddDate(0, 0, -10)
	if monthly.shouldRefresh(now) {
		t.Fatalf("expected monthly collector not to refresh in same month")
	}
	monthly.lastFetch = now.AddDate(0, -1, 0)
	if !monthly.shouldRefresh(now) {
		t.Fatalf("expected monthly collector to refresh in new month")
	}
}
