CREATE TABLE IF NOT EXISTS `photos` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `bytes` MEDIUMBLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS `users` (
    `name` VARCHAR(36) NOT NULL PRIMARY KEY,
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
    FOREIGN KEY (from_user) REFERENCES users(name),
    FOREIGN KEY (for_user) REFERENCES users(name)
);

CREATE TABLE IF NOT EXISTS `automobiles` (
    `reg_num` VARCHAR(16) NOT NULL PRIMARY KEY,
    `photo_id` VARCHAR(36) NOT NULL,
    `comment` VARCHAR(256),
    `owner_name` VARCHAR(36) NOT NULL UNIQUE,
    FOREIGN KEY (photo_id) REFERENCES photos(id),
    FOREIGN KEY (owner_name) REFERENCES users(name)
);

CREATE TABLE IF NOT EXISTS `trips` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `location_from` VARCHAR(64) NOT NULL,
    `location_to` VARCHAR(64) NOT NULL,
    `departure_time` DATETIME NOT NULL,
    `driver_name` VARCHAR(36) NOT NULL,
    `price` FLOAT NOT NULL,
    `max_passengers` INTEGER NOT NULL,
    `air_conditioning` BOOLEAN,
    `smoking` BOOLEAN,
    `pets` BOOLEAN,
    `comment` VARCHAR(256),
    FOREIGN KEY (driver_name) REFERENCES users(name)
);

CREATE TABLE IF NOT EXISTS `trip_participations` (
    `trip_id` VARCHAR(36) NOT NULL,
    `passenger_name` VARCHAR(36) NOT NULL,
    FOREIGN KEY (trip_id) REFERENCES trips(id),
    FOREIGN KEY (passenger_name) REFERENCES users(name),
    UNIQUE (trip_id, passenger_name)
);

CREATE TABLE IF NOT EXISTS `sessions` (
    `id` VARCHAR(128) NOT NULL PRIMARY KEY,
    `owner` VARCHAR(36) NOT NULL,
    `expiration` DATETIME NOT NULL,
    FOREIGN KEY (owner) REFERENCES users(name)
);