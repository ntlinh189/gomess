CREATE TABLE message_attachments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    message_id BIGINT UNSIGNED NOT NULL,
    type ENUM('image', 'video', 'audio', 'file') NOT NULL,
    object_key VARCHAR(1024) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(128) NOT NULL,
    size_bytes BIGINT UNSIGNED NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_attachment_message
        FOREIGN KEY (message_id) REFERENCES messages(id)
        ON DELETE CASCADE,

    INDEX idx_attachment_message (message_id)
);