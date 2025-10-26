CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    click_count INT DEFAULT 0
);

CREATE TABLE domain_counts (
    id SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL,
    count INT DEFAULT 1
);
