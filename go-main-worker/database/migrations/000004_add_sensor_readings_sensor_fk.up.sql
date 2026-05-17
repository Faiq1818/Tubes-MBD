ALTER TABLE sensor_readings
ADD CONSTRAINT fk_sensor_readings_sensor
FOREIGN KEY (sensor_id)
REFERENCES sensors(id)
ON DELETE RESTRICT;
