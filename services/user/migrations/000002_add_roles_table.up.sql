CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name) VALUES ('admin'), ('moderator'), ('user');

ALTER TABLE users ADD COLUMN role_id INT REFERENCES roles(id);

UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user');

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;