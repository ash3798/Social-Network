CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id serial PRIMARY KEY,
    comment_text VARCHAR(255) NOT NULL,
    sender_username VARCHAR(50),
    receiver_username VARCHAR(50) ,
    parent_comment_id INT ,
    comment_time BIGINT NOT NULL ,
    CONSTRAINT fk_sender_username
        FOREIGN KEY (sender_username)
            REFERENCES users(username) ,
    CONSTRAINT fk_receiver_username
        FOREIGN KEY (receiver_username)
            REFERENCES users(username)        
);

CREATE TABLE IF NOT EXISTS reactions (
    id serial PRIMARY KEY,
    reaction VARCHAR(10) CHECK (reaction = 'like' or reaction = 'dislike' or reaction = '+1'),
    comment_id INT NOT NULL,
    CONSTRAINT fk_comment_id
        FOREIGN KEY (comment_id)
            REFERENCES comments(id)
            ON DELETE CASCADE
);