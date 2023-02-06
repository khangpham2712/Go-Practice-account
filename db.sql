CREATE DATABASE IF NOT EXISTS `account`;

CREATE TABLE IF NOT EXISTS `account`.`accounts` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
);