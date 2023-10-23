ALTER TABLE
    images
ADD
    COLUMN processing_errors JSON
AFTER
    `status`;