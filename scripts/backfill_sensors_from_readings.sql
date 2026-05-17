BEGIN;

-- sensor_readings only stores sensor_id, so create a placeholder station first.
INSERT INTO stations (id, name, location, altitude, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Unknown Station',
    ST_SetSRID(ST_MakePoint(0, 0), 4326)::geography,
    0,
    true
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sensors (
    id,
    station_id,
    model_name,
    sampling_rate,
    sensitivity_threshold,
    last_calibration
)
SELECT DISTINCT
    sr.sensor_id,
    '00000000-0000-0000-0000-000000000001'::uuid,
    'unknown',
    0,
    0,
    NULL::date
FROM sensor_readings sr
WHERE NOT EXISTS (
    SELECT 1
    FROM sensors s
    WHERE s.id = sr.sensor_id
);

COMMIT;
