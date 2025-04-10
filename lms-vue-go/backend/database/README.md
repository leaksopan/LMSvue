# Database Setup Instructions

This directory contains the SQL scripts needed to set up the database for the LMS application.

## Prerequisites

- MySQL Server 5.7+ or MariaDB 10.2+
- MySQL client or phpMyAdmin

## Setup Instructions

### Option 1: Using MySQL Command Line

1. Open a terminal or command prompt
2. Navigate to this directory
3. Connect to MySQL server:
   ```
   mysql -u root -p
   ```
4. Enter your MySQL root password when prompted
5. Import the SQL script:
   ```
   source lms_db.sql
   ```

### Option 2: Using phpMyAdmin

1. Open phpMyAdmin in your web browser
2. Click on the "Import" tab
3. Click "Browse" and select the `lms_db.sql` file
4. Click "Go" to import the database

## Database Configuration

The application is configured to connect to the database with the following default settings:

- Host: localhost
- Port: 3306
- Username: root
- Password: (empty)
- Database: lms_db

If you need to change these settings, update the `DefaultConfig()` function in `backend/config/database.go`.

## Database Schema

The database consists of the following tables:

1. `users` - Stores user account information
2. `students` - Stores student information
3. `questions` - Stores questions for quizzes and tests
4. `student_answers` - Stores student answers to questions

## Default Users

The script creates the following default users:

1. Admin User:
   - Username: admin
   - Password: admin123
   - Email: admin@example.com

2. Teacher User:
   - Username: teacher
   - Password: admin123
   - Email: teacher@example.com

3. Student User:
   - Username: student
   - Password: admin123
   - Email: student@example.com

**Note:** For security reasons, you should change these passwords in a production environment.
