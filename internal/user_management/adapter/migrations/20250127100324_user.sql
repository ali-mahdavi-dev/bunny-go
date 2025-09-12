-- Active: 1757671557877@@127.0.0.1@5432@bunny_go
-- migrate:up
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    avatar_identifier VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Optional: Index for soft delete queries
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- migrate:down
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP TABLE IF EXISTS users;