# LMS Vue + Go

Simple Learning Management System built with Vue.js frontend and Go (Gin) backend.

## Project Structure

- `/backend`: Contains the Go (Gin) backend application.
- `/frontend/lms-app`: Contains the Vue.js frontend application.

## Setup

### Prerequisites

- Node.js (v16 or later recommended)
- Go (v1.18 or later recommended)

### Backend

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install Go dependencies:
    ```bash
    go mod tidy
    ```

### Frontend

1.  Navigate to the frontend directory:
    ```bash
    cd frontend/lms-app
    ```
2.  Install Node.js dependencies:
    ```bash
    npm install
    ```

## Running the Application

### Backend

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Run the Go server:
    ```bash
    go run main.go
    ```
    The backend will run on `http://localhost:3001`.

### Frontend

1.  Navigate to the frontend directory:
    ```bash
    cd frontend/lms-app
    ```
2.  Run the Vue development server:
    ```bash
    npm run serve
    ```
    The frontend will be accessible, usually on `http://localhost:8080` or another port if 8080 is busy.
