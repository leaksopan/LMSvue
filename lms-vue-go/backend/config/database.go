package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// DefaultConfig returns the default database configuration
func DefaultConfig() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     3306, // Update this if your MySQL is running on a different port
		User:     "root",
		Password: "", // Update this if you have a password set
		DBName:   "lms_db",
	}
}

// FormatDSN formats the DSN (Data Source Name) for MySQL connection
func (c *DBConfig) FormatDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

var (
	// DB is the global database connection
	DB *sql.DB
)

// InitDB initializes the database connection
func InitDB(config DBConfig) error {
	var err error

	// Log connection attempt
	dsn := config.FormatDSN()
	log.Printf("Connecting to MySQL database at %s:%d/%s", config.Host, config.Port, config.DBName)

	// Open database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
