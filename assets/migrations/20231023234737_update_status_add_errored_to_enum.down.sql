ALTER TABLE
    images
MODIFY
    COLUMN `status` ENUM('processing', 'processed') NOT NULL DEFAULT 'processing';