CREATE TABLE IF NOT EXISTS urls (
    id uuid PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL constraint uidx_urls_short_code unique,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    clicks BIGINT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_original_url ON urls(original_url);
CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at);
