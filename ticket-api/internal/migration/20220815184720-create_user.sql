-- +migrate Up
CREATE TABLE IF NOT EXISTS `users` (
  `id` BINARY(16) NOT NULL,
  `firebase_id` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100),
  `nick_name` VARCHAR(30),
  `username` VARCHAR(30) NOT NULL,
  `bio` VARCHAR(500),
  `genre` JSON,
  `gender` VARCHAR(10),
  `age` VARCHAR(30),
  `role` INT,
  `location` TEXT,
  `temp_email` TINYINT DEFAULT 0,
  `stripe_id` VARCHAR(100),
  `profile_picture` VARCHAR(100),
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME,
  PRIMARY KEY (`id`),
  CONSTRAINT User_Email_Unique UNIQUE(`email`),
  CONSTRAINT User_Firebase_Unique UNIQUE (`firebase_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
-- +migrate Down
DROP TABLE IF EXISTS `users`;