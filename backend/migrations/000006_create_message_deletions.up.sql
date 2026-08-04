CREATE TABLE message_deletions (
    message_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    deleted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (message_id, user_id),

    CONSTRAINT fk_message_deletion_message
        FOREIGN KEY (message_id) REFERENCES messages(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_message_deletion_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);