@echo off
echo Starting backend server...
start cmd /k "cd backend && go run main.go"
echo Starting frontend server...
start cmd /k "cd frontend/lms-app && npm run serve"
echo Servers are starting. Please wait a moment...
echo.
echo Backend URL: http://localhost:3001
echo Frontend URL: http://localhost:8080
echo Test API page: http://localhost:8080/test-api
echo.
echo Press any key to exit this window (servers will continue running)
pause > nul
