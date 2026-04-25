-- 20260426003_create_stats.up.sql
CREATE TABLE IF NOT EXISTS stats (
    id SERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- User-Agent данные
    browser VARCHAR(50),
    device_type VARCHAR(20),
    
    -- Гео данные (по IP)
    country VARCHAR(100),
    
    -- Источник перехода
    referer TEXT
);

CREATE INDEX idx_stats_link_id ON stats(link_id);
CREATE INDEX idx_stats_clicked_at ON stats(clicked_at);
CREATE INDEX idx_stats_browser ON stats(browser);
CREATE INDEX idx_stats_device_type ON stats(device_type);
CREATE INDEX idx_stats_country ON stats(country);