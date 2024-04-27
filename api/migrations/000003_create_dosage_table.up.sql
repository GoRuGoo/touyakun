CREATE TABLE IF NOT EXISTS dosage(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    time_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (time_id) REFERENCES time(id)
    );