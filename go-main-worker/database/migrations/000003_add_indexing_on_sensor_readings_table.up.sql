CREATE INDEX IF NOT EXISTS idx_sensor_readings_sensor_id
ON sensor_readings(sensor_id);

CREATE INDEX IF NOT EXISTS idx_sensor_readings_recorded_at
ON sensor_readings(recorded_at);

CREATE INDEX IF NOT EXISTS idx_sensor_readings_sensor_time
ON sensor_readings(sensor_id, recorded_at DESC);

CREATE INDEX IF NOT EXISTS idx_sensor_readings_trigger_true
ON sensor_readings(recorded_at)
WHERE is_trigger = true;
