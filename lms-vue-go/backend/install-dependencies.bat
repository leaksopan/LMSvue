@echo off
echo Installing Go dependencies...
go get github.com/go-sql-driver/mysql
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go mod tidy
echo Dependencies installed successfully!
echo.
echo Next steps:
echo 1. Set up the database by following the instructions in backend/database/README.md
echo 2. Run the server with 'go run main.go'
pause
