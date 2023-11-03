ALTER TABLE
    images
MODIFY
    COLUMN STATUS ENUM('processing', 'processed', 'errored', 'queued') NOT NULL DEFAULT 'queued';