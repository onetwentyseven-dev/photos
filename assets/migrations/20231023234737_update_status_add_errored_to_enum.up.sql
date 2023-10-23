ALTER TABLE
    images
MODIFY
    COLUMN STATUS ENUM('processing', 'processed', 'errored') NOT NULL DEFAULT 'processing';