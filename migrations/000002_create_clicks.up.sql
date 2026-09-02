CREATE TABLE clicks (
                        id BIGSERIAL PRIMARY KEY,
                        link_id BIGINT NOT NULL,
                        clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        referer TEXT NULL,
                        user_agent TEXT NULL,
                        ip_hash VARCHAR(64) NULL,

                        CONSTRAINT fk_clicks_link
                            FOREIGN KEY (link_id)
                                REFERENCES links(id)
                                ON DELETE CASCADE
);

CREATE INDEX idx_clicks_link_id
    ON clicks(link_id);

CREATE INDEX idx_clicks_clicked_at
    ON clicks(clicked_at);