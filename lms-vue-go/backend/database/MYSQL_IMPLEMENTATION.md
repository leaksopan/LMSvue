# MySQL Implementation in LMS System

This document explains how MySQL is integrated into the LMS (Learning Management System) application.

## Architecture Overview

The application follows a layered architecture:

1. **Handlers Layer** - Handles HTTP requests and responses
2. **Repository Layer** - Manages data access and persistence
3. **Models Layer** - Defines data structures
4. **Database Layer** - Manages database connections

## Database Connection

The database connection is managed in `config/database.go`. The application uses the standard Go `database/sql` package with the MySQL driver.

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)
```

The connection is initialized in the `main.go` file:

```go
dbConfig := config.DefaultConfig()
err := config.InitDB(dbConfig)
if err != nil {
    log.Fatal("Error connecting to database:", err)
}
defer config.CloseDB()
```

## Repository Pattern

The application uses the repository pattern to abstract database operations:

1. **UserRepository** - Manages user data
2. **StudentRepository** - Manages student data
3. **QuestionRepository** - Manages question data
4. **StudentAnswerRepository** - Manages student answers

Each repository provides methods for CRUD operations:

- `FindAll()` - Retrieves all records
- `FindByID()` - Retrieves a record by ID
- `Create()` - Creates a new record
- `Update()` - Updates an existing record
- `Delete()` - Deletes a record

## SQL Queries

The application uses raw SQL queries instead of an ORM. This provides more control over the database operations and makes it easier to optimize queries.

Example query from `StudentRepository.FindAll()`:

```go
query := `
    SELECT id, name, class, email
    FROM students s
    JOIN users u ON s.user_id = u.id
`
```

## Data Models

The data models are defined in the `models` package:

1. **User** - Represents a user account
2. **Student** - Represents a student
3. **Question** - Represents a question
4. **StudentAnswer** - Represents a student's answer to a question

Each model corresponds to a table in the database.

## Transaction Management

For operations that require multiple database operations, transactions are used to ensure data consistency.

Example:

```go
tx, err := db.Begin()
if err != nil {
    return err
}

// Perform multiple operations

if err != nil {
    tx.Rollback()
    return err
}

return tx.Commit()
```

## Error Handling

Database errors are handled at the repository layer and converted to appropriate HTTP responses in the handlers.

Example:

```go
if errors.Is(err, sql.ErrNoRows) {
    return nil, nil // Not found
}
return nil, err // Other error
```

## Security Considerations

1. **SQL Injection Prevention** - All queries use parameterized statements to prevent SQL injection
2. **Connection Pooling** - The database connection pool is configured with appropriate limits
3. **Error Handling** - Database errors are properly handled and don't expose sensitive information

## Database Schema

The database schema is defined in `database/lms_db.sql` and includes:

1. **users** - Stores user accounts
2. **students** - Stores student information
3. **questions** - Stores questions
4. **student_answers** - Stores student answers

## Future Improvements

1. **Migration System** - Implement a database migration system for schema changes
2. **Query Builder** - Consider using a query builder for complex queries
3. **Caching** - Implement caching for frequently accessed data
4. **Read/Write Splitting** - Separate read and write operations for better performance
