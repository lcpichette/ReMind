CREATE DATABASE dev;

USE dev;

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
CREATE TABLE Analytics
(
    id                  INT AUTO_INCREMENT UNIQUE,
    action_value        VARCHAR(1024),
    element_text        VARCHAR(2048),
    element_type        VARCHAR(64),
    event               VARCHAR(64),
    content_grouping    VARCHAR(64),
    created_at          DATETIME
);

INSERT INTO Users (name,username,password,created_at,last_seen) VALUES ( 'Tod', 'tod@gmail.com', 'Todster1987!', now(), now() );

INSERT INTO Messages (raw,created_at,user_id) VALUES ("Hello World!",now(),1);

SELECT * FROM Users;

SELECT * FROM Analytics;

SELECT raw, created_at FROM Messages WHERE user_id = 1;