package server

import (
	"strings"
	"testing"
)

func TestOptimizedSensorReadingJoinQueryUsesUUIDJoin(t *testing.T) {
	query := optimizedSensorReadingJoinQuery()

	if !strings.Contains(query, "JOIN sensors s ON s.id = sr.sensor_id") {
		t.Fatalf("expected optimized query to join UUID columns directly, got: %s", query)
	}

	if strings.Contains(query, "::text") {
		t.Fatalf("expected optimized query to avoid UUID text casts, got: %s", query)
	}
}

func TestUnoptimizedSensorReadingJoinQueryCastsUUIDsToText(t *testing.T) {
	query := unoptimizedSensorReadingJoinQuery()

	if !strings.Contains(query, "JOIN sensors s ON s.id::text = sr.sensor_id::text") {
		t.Fatalf("expected unoptimized query to cast UUID columns to text, got: %s", query)
	}
}
