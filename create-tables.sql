DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS room;
DROP TABLE IF EXISTS chat;

CREATE TABLE user (
  id           INT AUTO_INCREMENT NOT NULL,
  name         VARCHAR(128) NOT NULL,
  email        VARCHAR(128) NOT NULL,
  password     VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY unique_name (`email`)
);

CREATE TABLE room (
  id          INT AUTO_INCREMENT NOT NULL,
  user_id     VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE chat (
  id          INT AUTO_INCREMENT NOT NULL,
  user_id     VARCHAR(128) NOT NULL,
  room_id	  VARCHAR(128) NOT NULL,
  content     VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);
