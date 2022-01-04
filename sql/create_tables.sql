CREATE TABLE IF NOT EXISTS `photos` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `bytes` MEDIUMBLOB NOT NULL,
    `extension` VARCHAR(4) NOT NULL
);

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `name` VARCHAR(64) NOT NULL UNIQUE,
    `password` VARCHAR(36) NOT NULL,
    `email` VARCHAR(64) NOT NULL,
    `phone_number` VARCHAR(15) NOT NULL,
    `photo_id` VARCHAR(36),
    `times_passenger` INTEGER NOT NULL,
    `times_driver` INTEGER NOT NULL,
    FOREIGN KEY (photo_id) REFERENCES photos(id)
);

CREATE TABLE IF NOT EXISTS `reviews` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `from_user` VARCHAR(36) NOT NULL,
    `for_user` VARCHAR(36) NOT NULL,
    `rating` INTEGER NOT NULL,
    `comment` VARCHAR(256),
    FOREIGN KEY (from_user) REFERENCES users(id),
    FOREIGN KEY (for_user) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS `automobiles` (
    `regNum` VARCHAR(16) NOT NULL PRIMARY KEY,
    `photo_id` VARCHAR(36),
    `comment` VARCHAR(256),
    `owner_id` VARCHAR(36) NOT NULL,
    FOREIGN KEY (photo_id) REFERENCES photos(id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS `trips` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `from` VARCHAR(64) NOT NULL,
    `to` VARCHAR(64) NOT NULL,
    `when` DATETIME NOT NULL,
    `automobile` VARCHAR(16),
    `price` FLOAT NOT NULL,
    `max_passengers` INTEGER NOT NULL,
    `air_conditioning` BOOLEAN,
    `smoking` BOOLEAN,
    `pets` BOOLEAN,
    `comment` VARCHAR(256),
    FOREIGN KEY (automobile) REFERENCES automobiles(regNum)
);

CREATE TABLE IF NOT EXISTS `trip_participations` (
    `trip_id` VARCHAR(36) NOT NULL,
    `driver_id` VARCHAR(36) NOT NULL,
    `passenger_id` VARCHAR(36) NOT NULL,
    FOREIGN KEY (trip_id) REFERENCES trips(id),
    FOREIGN KEY (driver_id) REFERENCES users(id),
    FOREIGN KEY (passenger_id) REFERENCES users(id)
);