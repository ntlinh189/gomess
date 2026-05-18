CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,

    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    avatar TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_provider_provider_id (
        provider,
        provider_id
    ),

    UNIQUE KEY uk_email (
        email
    )
);