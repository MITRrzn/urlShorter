CREATE TABLE link_stats_daily (
                                  id BIGSERIAL PRIMARY KEY,
                                  link_id BIGINT NOT NULL,
                                  stat_date DATE NOT NULL,
                                  clicks_count BIGINT NOT NULL DEFAULT 0,

                                  CONSTRAINT fk_link_stats_daily_link
                                      FOREIGN KEY (link_id)
                                          REFERENCES links(id)
                                          ON DELETE CASCADE,

                                  CONSTRAINT uq_link_stats_daily_link_date
                                      UNIQUE (link_id, stat_date)
);

CREATE INDEX idx_link_stats_daily_stat_date
    ON link_stats_daily(stat_date);