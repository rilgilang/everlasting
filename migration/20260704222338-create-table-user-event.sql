-- +migrate Up
CREATE TABLE user_events (
                        id VARCHAR(255) PRIMARY KEY,
                        user_id VARCHAR(255) NOT NULL,
                        event_id jsonb,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_event_user_id ON user_events(user_id);

-- +migrate Down
DROP TABLE user_events;
