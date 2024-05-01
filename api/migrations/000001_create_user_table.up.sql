CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    line_user_id VARCHAR(255) NOT NULL UNIQUE,
    morning_medication_time TIME DEFAULT '08:00:00',
    afternoon_medication_time TIME DEFAULT '12:00:00',
    evening_medication_time TIME DEFAULT '18:00:00',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
