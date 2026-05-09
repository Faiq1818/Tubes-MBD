CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS stations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    altitude FLOAT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sensors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id UUID NOT NULL,
    model_name VARCHAR(50) NOT NULL,
    sampling_rate INT NOT NULL,
    sensitivity_threshold FLOAT NOT NULL,
    last_calibration DATE,
    
    CONSTRAINT fk_station
        FOREIGN KEY(station_id) 
        REFERENCES stations(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS earthquake_event (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    detected_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    epicenter GEOGRAPHY(POINT, 4326) NOT NULL,
    magnitude_richter FLOAT NOT NULL,
    depth_km FLOAT NOT NULL,
    validation_status VARCHAR(20) DEFAULT 'pending',
    
    CONSTRAINT check_validation_status 
        CHECK (validation_status IN ('pending', 'confirmed', 'false alarm'))
);

CREATE INDEX IF NOT EXISTS idx_stations_location ON stations USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_earthquake_epicenter ON earthquake_event USING GIST(epicenter);
CREATE INDEX IF NOT EXISTS idx_sensors_station_id ON sensors(station_id);
