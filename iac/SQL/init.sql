-- CREATE DATABASE dev;

USE myrdsinstance; -- db id

CREATE TABLE Users
(
    id          INT AUTO_INCREMENT UNIQUE,
    name        VARCHAR(64),
    username    VARCHAR(255),
    -- TODO: Hashed passwords are far more secure
    password    VARCHAR(255),
    created_at  DATETIME,
    last_seen   DATETIME,
    PRIMARY KEY (id)
);
CREATE TABLE Messages
(
    id          INT AUTO_INCREMENT UNIQUE,
    raw         VARCHAR(1024),
    rich        VARCHAR(1024),
    created_at  DATETIME,
    user_id     INT,
    CONSTRAINT uid
        FOREIGN KEY (user_id)
        REFERENCES Users(id),
    PRIMARY KEY (id)
);
