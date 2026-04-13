-- +migrate Up
CREATE TABLE events (
                       id VARCHAR(255) PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description VARCHAR(255) NOT NULL,
                       date DATE NOT NULL,
                       time VARCHAR(255) NOT NULL DEFAULT '',
                       location VARCHAR(255),
                       category VARCHAR(255),
                       messages INTEGER,
                       max_messages INTEGER,
                       image VARCHAR(255),
                       status VARCHAR(255),
                       organizer VARCHAR(255),
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_event ON events(id);
CREATE INDEX idx_title ON events(title);

-- +migrate Down
DROP TABLE events;

