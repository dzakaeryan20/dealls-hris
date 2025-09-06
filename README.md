# Dokumentasi Proyek HRIS Payroll

Selamat datang di dokumentasi resmi HRIS Payroll. Sistem ini adalah sebuah API backend yang dirancang untuk mengelola penggajian karyawan berdasarkan kehadiran, lembur, dan *reimbursement*.

## 1. Arsitektur Perangkat Lunak

Aplikasi ini dibangun menggunakan arsitektur **Modular Monolith**. Ini berarti seluruh logika aplikasi berada dalam satu unit *deployment* (monolith), tetapi kode diorganisir ke dalam modul-modul independen berdasarkan domain bisnisnya. Pendekatan ini menawarkan kemudahan pengembangan seperti monolith, sambil mempertahankan keteraturan dan skalabilitas seperti *microservices*.



### Struktur Direktori Utama
-   **`/cmd`**: Titik masuk (entrypoint) aplikasi. Di sinilah fungsi `main` berada, yang bertugas menginisialisasi dan menyatukan semua komponen.
-   **`/internal`**: Direktori ini berisi seluruh kode inti aplikasi yang tidak boleh diimpor oleh proyek lain.
    -   **`/api`**: Menangani semua hal yang berkaitan dengan lapisan HTTP.
        -   `/handler`: Menerjemahkan *request* HTTP menjadi panggilan ke *service*.
        -   `/middleware`: Berisi "pos pemeriksaan" untuk setiap *request*, seperti otentikasi (JWT), pelacakan IP, dan *request ID*.
        -   `/router`: Mendefinisikan semua *endpoint* API dan menghubungkannya ke *handler* yang sesuai.
    -   **`/domain`**: Jantung dari aplikasi. Berisi logika bisnis murni.
        -   `/auth`, `/payroll`, `/attendance`, `/overtime`, `/reimbursement` , `/user`: Setiap folder adalah modul domain yang memiliki `model`, `service` (logika bisnis), dan `repository` (kontrak ke database).
    -   **`/platform`**: Berisi kode yang berinteraksi dengan dunia luar.
        -   `/database`: Konfigurasi dan koneksi ke database PostgreSQL.
        -   `/seeder`: Logika untuk mengisi data awal (dummy data) ke database.
-   **`/pkg`**: (Opsional) Digunakan untuk kode yang aman untuk dibagikan dan diimpor oleh proyek lain.

### Alur Data
Sebuah *request* dari klien akan mengikuti alur berikut:
**Router** → **Middleware** (Auth, IP Tracker, dll.) → **Handler** (Validasi Input) → **Service** (Eksekusi Logika Bisnis) → **Repository** (Akses Database) → **Database**

---
## 2. Panduan Penggunaan (How-To Guides)

### Menjalankan Proyek Secara Lokal
Proyek ini dirancang untuk berjalan dengan mudah menggunakan Docker.

1.  **Prasyarat**: Pastikan Anda memiliki **Docker** dan **Docker Compose** ter-install di sistem Anda.
2.  **Konfigurasi**: Salin file `.env.example` menjadi `.env`.
    ```bash
    cp .env.example .env
    ```
3.  **Sesuaikan `.env`**: Buka file `.env` dan sesuaikan konfigurasinya jika perlu. Untuk menjalankan pertama kali, pastikan `RUN_SEEDER=true` untuk mengisi database dengan data admin dan 100 karyawan.
4.  **Jalankan Aplikasi**: Buka terminal di direktori utama proyek dan jalankan:
    ```bash
    docker-compose up --build
    ```
5.  **Akses API**: Server API akan berjalan dan dapat diakses di `http://localhost:8080`.

### Menjalankan Pengujian (Testing)
Aplikasi ini dilengkapi dengan dua jenis tes: *unit test* dan *integration test*.

-   **Menjalankan Unit Test (Cepat)**: Tes ini tidak memerlukan database dan hanya memverifikasi logika di *service layer*.
    ```bash
    go test -v ./...
    ```
-   **Menjalankan Integration Test (Lengkap)**: Tes ini akan secara otomatis menjalankan kontainer database PostgreSQL temporer menggunakan Docker untuk menguji seluruh alur API.
    ```bash
    go test -v ./... -tags=integration
    ```

---
## 3. Referensi API

Semua *endpoint* yang membutuhkan otentikasi harus menyertakan *header* berikut:
`Authorization: Bearer <your_jwt_token>`

### 🏛️ Otentikasi

#### `POST /api/v1/auth/login`
-   **Deskripsi**: Mengotentikasi pengguna dan mengembalikan JWT.
-   **Request Body**:
    ```json
    {
        "username": "admin",
        "password": "password123"
    }
    ```
-   **Response Sukses (200 OK)**:
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

---
### 👨‍💼 Endpoint Karyawan

#### `POST /api/v1/attendance`
-   **Deskripsi**: Mengajukan absensi untuk hari ini.
-   **Otentikasi**: Perlu token **Karyawan**.
-   **Request Body**: Kosong.
-   **Response Sukses (201 Created)**:
    ```json
    {
        "message": "Attendance submitted successfully for today"
    }
    ```

#### `POST /api/v1/overtime`
-   **Deskripsi**: Mengajukan jam lembur.
-   **Otentikasi**: Perlu token **Karyawan**.
-   **Request Body**:
    ```json
    {
        "date": "2025-09-10",
        "hours": 2
    }
    ```
-   **Response Sukses (201 Created)**:
    ```json
    {
        "message": "Overtime submitted successfully"
    }
    ```

#### `POST /api/v1/reimbursement`
-   **Deskripsi**: Mengajukan *reimbursement*.
-   **Otentikasi**: Perlu token **Karyawan**.
-   **Request Body**:
    ```json
    {
        "date": "2025-09-11",
        "description": "Biaya makan siang dengan klien",
        "amount": 150000
    }
    ```
-   **Response Sukses (201 Created)**:
    ```json
    {
        "message": "Reimbursement submitted successfully"
    }
    ```

#### `GET /api/v1/payslip/{period_id}`
-   **Deskripsi**: Melihat slip gaji pribadi untuk periode tertentu.
-   **Otentikasi**: Perlu token **Karyawan**.
-   **Response Sukses (200 OK)**:
    ```json
    {
        "id": "payslip-uuid",
        "user_id": "employee-uuid",
        "payroll_period_id": "period-uuid",
        "base_salary": 10000000,
        "prorated_salary": 9500000,
        "overtime_pay": 500000,
        "reimbursement_total": 150000,
        "total_pay": 10150000,
        "created_at": "...",
        "updated_at": "...",
        "created_by": "admin-uuid",
        "updated_by": "admin-uuid"
    }
    ```

---
### ⚙️ Endpoint Admin

#### `POST /api/v1/admin/payroll-period`
-   **Deskripsi**: Membuat periode penggajian baru.
-   **Otentikasi**: Perlu token **Admin**.
-   **Request Body**:
    ```json
    {
        "start_date": "2025-09-01",
        "end_date": "2025-09-30"
    }
    ```
-   **Response Sukses (201 Created)**:
    ```json
    {
        "id": "period-uuid-baru",
        "start_date": "2025-09-01T00:00:00Z",
        "end_date": "2025-09-30T00:00:00Z",
        "status": "pending",
        // ...
    }
    ```

#### `POST /api/v1/admin/payroll/{period_id}/run`
-   **Deskripsi**: Menjalankan dan memproses kalkulasi gaji untuk semua karyawan dalam satu periode.
-   **Otentikasi**: Perlu token **Admin**.
-   **Request Body**: Kosong.
-   **Response Sukses (200 OK)**:
    ```json
    {
        "message": "Payroll run successfully"
    }
    ```

#### `GET /api/v1/admin/payroll/{period_id}/summary`
-   **Deskripsi**: Mendapatkan ringkasan total pengeluaran gaji untuk satu periode.
-   **Otentikasi**: Perlu token **Admin**.
-   **Response Sukses (200 OK)**:
    ```json
    {
        "payroll_period_id": "period-uuid",
        "employee_pays": [
            {
                "user_id": "employee-uuid-1",
                "take_home_pay": 10150000
            },
            {
                "user_id": "employee-uuid-2",
                "take_home_pay": 9800000
            }
        ],
        "total_payout": 19950000
    }
    ```