-- Migration script to add score column to questions table

-- Check if score column exists, if not add it
SET @exist := (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
               WHERE TABLE_SCHEMA = 'lms_db' 
               AND TABLE_NAME = 'questions' 
               AND COLUMN_NAME = 'score');

SET @query = IF(@exist = 0, 
                'ALTER TABLE questions ADD COLUMN score INT NOT NULL DEFAULT 1',
                'SELECT "Score column already exists"');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Update existing questions with default scores if the column was just added
UPDATE questions SET score = 1 WHERE score IS NULL;

-- Update specific questions with different scores
UPDATE questions SET score = 5 WHERE id = 2; -- Essay question
UPDATE questions SET score = 2 WHERE id = 3; -- JavaScript question
UPDATE questions SET score = 3 WHERE id = 4; -- Framework question
