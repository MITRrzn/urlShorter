CREATE TABLE links (
                       id BIGSERIAL PRIMARY KEY,
                       short_code VARCHAR(7) NOT NULL UNIQUE,
                       original_url TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       expires_at TIMESTAMPTZ NULL
);