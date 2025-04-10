-- Create database if not exists
CREATE DATABASE IF NOT EXISTS lms_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE lms_db;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    role ENUM('admin', 'teacher', 'student') NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Create students table
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    class VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Create questions table
CREATE TABLE IF NOT EXISTS questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('multiple_choice', 'essay') NOT NULL,
    question TEXT NOT NULL,
    options JSON NULL,
    answer TEXT NULL,
    image_url VARCHAR(255) NULL,
    score INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Create student_answers table
CREATE TABLE IF NOT EXISTS student_answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT NOT NULL,
    question_id INT NOT NULL,
    answer TEXT NOT NULL,
    score INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Insert default admin user (password: admin123)
INSERT INTO users (username, password, email, role) VALUES
('admin', 'admin123', 'admin@example.com', 'admin'),
('teacher', 'admin123', 'teacher@example.com', 'teacher'),
('student', 'admin123', 'student@example.com', 'student');

-- Insert sample students
INSERT INTO students (user_id, name, class) VALUES
(3, 'Budi Santoso', '10A'),
(3, 'Ani Wijaya', '10A'),
(3, 'Dian Pratama', '10B'),
(3, 'Rini Susanti', '10B'),
(3, 'Ahmad Rizki', '10C');

-- Insert sample questions
INSERT INTO questions (type, question, options, answer, image_url, score) VALUES
('multiple_choice', 'Ibukota Indonesia adalah...', '["Jakarta", "Bandung", "Surabaya", "Yogyakarta", "Medan"]', 'A', NULL, 1),
('essay', 'Jelaskan mengapa belajar pemrograman penting di era digital?', NULL, NULL, NULL, 5),
('multiple_choice', 'Bahasa pemrograman yang berjalan di lingkungan browser adalah...', '["JavaScript", "Java", "Python", "Go", "C++"]', 'A', NULL, 2),
('multiple_choice', 'Mana yang bukan termasuk framework JavaScript?', '["Django", "React", "Angular", "Vue", "Svelte"]', 'A', 'https://example.com/frameworks.jpg', 3);
