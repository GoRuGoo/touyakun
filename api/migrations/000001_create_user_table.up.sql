CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    line_user_id VARCHAR(255) NOT NULL UNIQUE,
    morning_medication_time TIME,
    afternoon_medication_time TIME,
    evening_medication_time TIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
