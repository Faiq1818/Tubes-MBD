
CREATE TABLE IF NOT EXISTS sensor_readings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sensor_id UUID NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    
    acc_x DOUBLE PRECISION NOT NULL,
    acc_y DOUBLE PRECISION NOT NULL,
    acc_z DOUBLE PRECISION NOT NULL,
    
    pga DOUBLE PRECISION,
    sta_lta DOUBLE PRECISION,
    is_trigger BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

