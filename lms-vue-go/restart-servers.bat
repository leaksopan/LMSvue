@echo off
echo Killing any existing processes on port 3001...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3001 ^| findstr LISTENING') do (
    echo Killing process with PID %%a
    taskkill /F /PID %%a 2>nul
)

echo Starting backend server...
start cmd /k "cd backend && go run main.go"
echo Waiting for backend to start...
timeout /t 5 /nobreak > nul

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
