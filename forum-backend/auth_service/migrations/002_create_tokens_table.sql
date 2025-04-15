-- +migrate Up
CREATE TABLE IF NOT EXISTS tokens (
                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
                                      user_id INTEGER NOT NULL,
                                      token VARCHAR(255) NOT NULL UNIQUE,
    expiry_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
    );

-- +migrate Down
DROP TABLE IF EXISTS tokens;