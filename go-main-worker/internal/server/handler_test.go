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

func TestNPlusOneSensorReadingsQueryDoesNotJoinSensors(t *testing.T) {
	query := nPlusOneSensorReadingsQuery()

	if strings.Contains(query, "JOIN sensors") {
		t.Fatalf("expected N+1 readings query to avoid join, got: %s", query)
	}

	if !strings.Contains(query, "FROM sensor_readings") {
		t.Fatalf("expected N+1 readings query to select from sensor_readings, got: %s", query)
	}
}

func TestNPlusOneSensorQueryLoadsOneSensorByID(t *testing.T) {
	query := nPlusOneSensorQuery()

	if !strings.Contains(query, "FROM sensors") {
		t.Fatalf("expected N+1 sensor query to select from sensors, got: %s", query)
	}

	if !strings.Contains(query, "WHERE id = $1") {
		t.Fatalf("expected N+1 sensor query to lookup one sensor by id, got: %s", query)
	}
}

func TestNoNPlusOneSensorReadingsQueryUsesSingleJoin(t *testing.T) {
	query := noNPlusOneSensorReadingsQuery()

	if !strings.Contains(query, "JOIN sensors s ON s.id = sr.sensor_id") {
		t.Fatalf("expected no-N+1 query to join sensors in one query, got: %s", query)
	}
}
