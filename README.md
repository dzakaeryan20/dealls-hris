# Go Payroll System Boilerplate

A modular monolithic API for a simple payroll system built with Go and PostgreSQL.

## Features

-   User Authentication (JWT)
-   Admin and Employee Roles
-   Employee Submissions (Attendance, Overtime, Reimbursement)
-   Admin Payroll Processing
-   Payslip Generation
-   Payroll Summary Reports
-   Context-aware for cancellation and timeouts.

## 🚀 How to Run

1.  **Clone the repository.**

2.  **Create a `.env` file:**
    ```bash
    cp .env.example .env
    ```
    Edit the `.env` file with your desired settings. **Important:** For the very first run, ensure `RUN_SEEDER=true` to populate the database with fake data. You can set it to `false` afterward.

3.  **Run with Docker Compose:**
    ```bash
    docker-compose up --build
    ```

The API will be available at `http://localhost:8080`.

## 👨‍💻 API Usage

### Authentication

-   `POST /api/v1/auth/login`
    -   Login with user credentials to get a JWT.
    -   Admin: `{"username": "admin", "password": "password123"}`
    -   Employee: `{"username": "employee1", "password": "password123"}`

**Note:** All other endpoints require an `Authorization: Bearer <your_jwt_token>` header.

### Employee Endpoints

-   `POST /api/v1/attendance` (Submits attendance for today)
-   `POST /api/v1/overtime`
    -   Body: `{"date": "2025-09-05", "hours": 2}`
-   `POST /api/v1/reimbursement`
    -   Body: `{"date": "2025-09-05", "description": "Client lunch", "amount": 50.50}`
-   `GET /api/v1/payslip/{period_id}` (Get your payslip for a specific period)

### Admin Endpoints

-   `POST /api/v1/admin/payroll-period`
    -   Body: `{"start_date": "2025-09-01", "end_date": "2025-09-30"}`
-   `POST /api/v1/admin/payroll/{period_id}/run` (Processes payroll for the period)
-   `GET /api/v1/admin/payroll/{period_id}/summary` (Gets a summary of all payslips for the period)