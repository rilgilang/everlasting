-- +migrate Up
CREATE TABLE guests (
    id VARCHAR(255) PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'regular',
    is_invitation_sent BOOLEAN NOT NULL DEFAULT FALSE,
    last_invitation_sent TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_guests_event_id ON guests(event_id);
CREATE INDEX idx_guests_status ON guests(status);
CREATE INDEX idx_guests_deleted_at ON guests(deleted_at);

-- +migrate Down
DROP TABLE guests;
