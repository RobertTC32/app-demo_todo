-- +goose Up
CREATE TABLE IF NOT EXISTS todos (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    -- auto_increment starts from 1
    title VARCHAR(30) CHARACTER SET utf8 NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME DEFAULT NULL,
    -- nullable, only set when completed is true
    version INT UNSIGNED NOT NULL DEFAULT 1,
    -- used for optimistic locking
    PRIMARY KEY (id),
    UNIQUE KEY unique_title (title)
);
-- +goose Down
DROP TABLE IF EXISTS todos;