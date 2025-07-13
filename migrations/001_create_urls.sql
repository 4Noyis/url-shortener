CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    clicks INTEGER DEFAULT 0,
    user_id VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_short_url ON urls(short_url);
CREATE INDEX idx_long_url_hash ON urls(MD5(long_url));
CREATE INDEX idx_created_at ON urls(created_at);