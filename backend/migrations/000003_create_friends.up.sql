CREATE TABLE friends (
    user1_id BIGINT UNSIGNED NOT NULL,
    user2_id BIGINT UNSIGNED NOT NULL,

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(
        user1_id,
        user2_id
    ),

    CONSTRAINT fk_friend_user1
        FOREIGN KEY(user1_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_friend_user2
        FOREIGN KEY(user2_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_friend_users
        CHECK(user1_id < user2_id)
);