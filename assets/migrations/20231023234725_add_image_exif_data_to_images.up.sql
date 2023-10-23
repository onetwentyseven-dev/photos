ALTER TABLE
    images
ADD
    COLUMN image_exif_data JSON
AFTER
    processing_errors;