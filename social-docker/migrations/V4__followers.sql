-- CREATE TABLE IF NOT EXISTS followers (
--     id BIGINT PRIMARY KEY AUTO_INCREMENT,
--     follower_id BIGINT NOT NULL,
--     followee_id BIGINT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
--     FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
-- );


CREATE TABLE IF NOT EXISTS followers (
    user_id UUID NOT NULL,
    follower_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
);