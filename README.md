# Sistem Absensi & Payroll Karyawan

Aplikasi ini dibangun dengan **Go + Gin + GORM** untuk mengelola sistem absensi harian dan payroll karyawan, lengkap dengan autentikasi berbasis **JWT**.

---

## ⚙️ Teknologi

- Go (Golang)
- Gin (Web Framework)
- GORM (ORM untuk database)
- JWT (JSON Web Token)

---

## 📂 Relasi Tabel

- Departments **(1) -- (n)** Employees
- Employees **(1) -- (n)** Attendances
- Employees **(1) -- (n)** Payrolls

---

## 🚀 Cara Menjalankan

1. Clone repo ini

   ```bash
   git clone https://github.com/username/sistem-absensi-payroll.git
   cd sistem-absensi-payroll
   ```

2. Atur database di file `config/config.go`

3. Jalankan migrasi otomatis

   ```go
   config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})
   ```

4. (Opsional) Jalankan seeding data:

   ```go
   migrations.Seed()
   ```

5. Jalankan aplikasi
   ```bash
   go run main.go
   ```

Aplikasi akan jalan di `http://localhost:8080`

---

## 🔑 Endpoint API

### 1. Authentication

#### Register

```
POST /api/auth/register
```

Request:

```json
{
  "name": "Fathir achmad",
  "username": "fathir",
  "password": "123456",
  "department_id": 1
}
```

Response:

```json
{ "message": "User registered successfully" }
```

#### Login

```
POST /api/auth/login
```

Request:

```json
{
  "username": "fathir",
  "password": "123456"
}
```

Response:

```json
{ "token": "JWT_TOKEN_HERE" }
```

---

### 2. Attendance

> Semua endpoint attendance butuh **Authorization: Bearer <token>**

#### Check-in

```
POST /api/attendances/checkin
```

Response:

```json
{
  "message": "Checked in successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "date": "2025-08-27",
    "check_in_at": "2025-08-27T08:01:00Z"
  }
}
```

#### Check-out

```
POST /api/attendances/checkout
```

Response:

```json
{
  "message": "Checked out successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "date": "2025-08-27",
    "check_in_at": "2025-08-27T08:01:00Z",
    "check_out_at": "2025-08-27T17:03:00Z"
  }
}
```

---

### 3. Payroll

> Butuh **Authorization: Bearer <token>**

#### Get Payroll (Employee aktif)

```
GET /api/attendances/payroll
```

Response:

```json
{
  "employee_id": 1,
  "payrolls": [
    {
      "amount": 3000000,
      "period_start": "2025-08-01",
      "period_end": "2025-08-31"
    }
  ]
}
```

---

### 4. Departments

> Butuh **Authorization: Bearer <token>**

#### Create Department

```
POST /api/departments
```

Request:

```json
{ "name": "Finance" }
```

Response:

```json
{
  "message": "Department created successfully",
  "data": { "id": 1, "name": "Finance" }
}
```

#### Get All Departments

```
GET /api/departments
```

Response:

```json
{
  "data": [
    { "id": 1, "name": "Finance" },
    { "id": 2, "name": "HRD" }
  ]
}
```

#### Update Department

```
PATCH /api/departments/1
```

Request:

```json
{ "name": "Updated Department" }
```

Response:

```json
{ "message": "Department updated successfully" }
```

#### Delete Department

```
DELETE /api/departments/1
```

Response:

```json
{ "message": "Department deleted successfully" }
```

---

### 5. Employee (CRUD Profile)

> Semua endpoint employee butuh **Authorization: Bearer <token>**

#### Get Profile

```
GET /api/employee/me
```

Response:

```json
{
  "id": 1,
  "name": "Fathir achmad",
  "username": "fathir",
  "department": { "id": 1, "name": "HRD" }
}
```

#### Update Profile

```
PUT /api/employee/me
```

Request:

```json
{
  "name": "Updated Name",
  "password": "newpassword123",
  "department_id": 2
}
```

Response:

```json
{ "message": "Profile updated successfully" }
```

#### Delete Account

```
DELETE /api/employee/me
```

Response:

```json
{ "message": "Account deleted successfully" }
```

---

## 🗄️ Seed Data (Opsional)

Tambahkan file `migrations/seed.go`:

```go
package migrations

import (
    "attendance-payroll/models"
    "attendance-payroll/config"
)

func Seed() {
    config.DB.FirstOrCreate(&models.Department{}, models.Department{Name: "HRD"})
}
```

Panggil `migrations.Seed()` di `main.go` setelah `AutoMigrate`.

---

## ✅ Fitur

- [x] Register & Login dengan JWT
- [x] Absensi Check-in & Check-out (1 record per hari)
- [x] Perhitungan Payroll berdasarkan jam kerja (>= 8 jam = Rp100.000/hari)
- [x] CRUD Department
- [x] CRUD Employee (Profile: Get, Update, Delete)
