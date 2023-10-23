CREATE TABLE IF NOT EXISTS images (
    `id` VARCHAR(256) PRIMARY KEY,
    `user_id` VARCHAR(256) NOT NULL,
    `name` VARCHAR(256) NOT NULL,
    `description` TEXT,
    `status` ENUM('processing', 'processed') NOT NULL DEFAULT 'processing',
    `ts_created` DATETIME NOT NULL,
    `ts_updated` DATETIME NOT NULL,
    KEY `fk_images_user_id` (`user_id`)
);