-- +migrate Up
CREATE TABLE users (
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
	ciphertext VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(255) DEFAULT 'active',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user ON users(status);
CREATE INDEX idx_user_email_status ON users(email, status);
CREATE INDEX idx_user_role ON users(role);
CREATE INDEX idx_user_status ON users(status);

INSERT INTO users (
    id,
    name,
    email,
    role
) VALUES(
    '21ae4024-395b-472f-98fa-945382183418', 
    'admin', 
    'harun@digitalsekuriti.id', 
    'admin'
);

-- +migrate Down
DROP TABLE users;

