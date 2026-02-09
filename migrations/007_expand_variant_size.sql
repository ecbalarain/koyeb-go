-- Migration: Expand variant size column for longer labels
ALTER TABLE variants
    MODIFY COLUMN size VARCHAR(50) NOT NULL;
