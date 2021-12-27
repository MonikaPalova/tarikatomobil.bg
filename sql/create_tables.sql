CREATE TABLE IF NOT EXISTS `photos` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `bytes` MEDIUMBLOB NOT NULL,
    `extension` VARCHAR(4) NOT NULL
);

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `name` VARCHAR(64) NOT NULL,
    `email` VARCHAR(64) NOT NULL,
    `phone_number` VARCHAR(15) NOT NULL,
    `photo_id` VARCHAR(36),
    `times_passenger` INTEGER NOT NULL,
    `times_driver` INTEGER NOT NULL
);
