# Barang API

REST API untuk mengelola data barang menggunakan **Go, Gin, GORM, MariaDB, dan JWT Authentication**.

Project ini menyediakan autentikasi pengguna dan CRUD barang dengan validasi input, pagination, pencarian, pengurutan aman, serta soft delete.

## Fitur

- REST API menggunakan Go dan Gin.
- Database MariaDB dengan GORM ORM.
- Auto migration untuk tabel pengguna dan barang.
- Autentikasi JWT menggunakan algoritma HS256.
- Password disimpan menggunakan hash bcrypt.
- Endpoint profil pengguna yang dilindungi JWT.
- CRUD barang lengkap.
- Validasi SKU, nama, deskripsi, harga, dan stok.
- Pencegahan SKU duplikat.
- Pagination dengan metadata total data dan total halaman.
- Pencarian berdasarkan SKU atau nama barang.
- Pengurutan berdasarkan kolom yang diizinkan.
- Soft delete menggunakan GORM.
- Health check API dan koneksi database.
- Konfigurasi melalui file `.env`.

## Teknologi

| Teknologi | Fungsi                     |
| --------- | -------------------------- |
| Go        | Bahasa pemrograman backend |
| Gin       | HTTP framework             |
| GORM      | ORM untuk database         |
| MariaDB   | Database relasional        |
| JWT       | Autentikasi API            |
| bcrypt    | Hash password              |
| godotenv  | Membaca konfigurasi `.env` |

## Persyaratan

Pastikan perangkat sudah memiliki:

- Go yang kompatibel dengan versi pada `go.mod`.
- MariaDB yang sedang berjalan.
- Git (opsional).
- PowerShell untuk menjalankan contoh pengujian pada Windows.

## Struktur Project

```text
barang-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   └── barang_handler.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   ├── user.go
│   │   └── barang.go
│   ├── repository/
│   ├── routes/
│   └── utils/
│       └── response.go
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Instalasi

### 1. Clone atau masuk ke project

Jika project sudah dibuat secara lokal:

```powershell
cd barang-api
```

### 2. Unduh dependency

```powershell
go mod tidy
```

### 3. Buat database MariaDB

Masuk ke MariaDB menggunakan akun administrator:

```sql
CREATE DATABASE db_barang
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE USER 'barang_app'@'localhost'
    IDENTIFIED BY 'GANTI_PASSWORD_DATABASE';

GRANT ALL PRIVILEGES ON db_barang.*
    TO 'barang_app'@'localhost';

FLUSH PRIVILEGES;
```

Ganti password database dengan password yang kuat.

Jika akun `barang_app` sudah dibuat pada tahap sebelumnya, tidak perlu membuatnya kembali. Pastikan akun tersebut mempunyai izin yang diperlukan pada database `db_barang`.

### 4. Konfigurasi environment

Buat file `.env` di root project berdasarkan `.env.example`.

Contoh:

```dotenv
APP_ENV=development
APP_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=barang_app
DB_PASSWORD=GANTI_PASSWORD_DATABASE
DB_NAME=db_barang

JWT_SECRET=GANTI_DENGAN_SECRET_RANDOM_MINIMAL_32_BYTE
JWT_EXPIRE_MINUTES=60

ADMIN_USERNAME=admin
ADMIN_PASSWORD=GANTI_DENGAN_PASSWORD_ADMIN_KUAT
```

Sesuaikan nilai database dengan konfigurasi MariaDB lokal.

**Penting:**

- Jangan menggunakan contoh password di atas untuk production.
- `JWT_SECRET` harus berupa secret acak yang kuat, minimal 32 karakter sesuai validasi konfigurasi aplikasi.
- Password admin minimal 12 byte dan maksimal 72 byte karena menggunakan bcrypt.
- Jangan mengunggah `.env` ke repository Git.

Contoh membuat JWT secret acak menggunakan PowerShell:

```powershell
$bytes = New-Object byte[] 32
$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
$rng.GetBytes($bytes)

$secret = [Convert]::ToBase64String($bytes)

$rng.Dispose()
Write-Output $secret
```

Salin hasilnya ke `JWT_SECRET` di `.env`.

### 5. Jalankan pemeriksaan kode

```powershell
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/api
```

### 6. Jalankan server

```powershell
go run ./cmd/api
```

Secara default API berjalan pada:

```text
http://localhost:8080
```

Pada startup, aplikasi menghubungkan database, menjalankan auto migration, dan membuat akun admin awal jika belum ada.

**Catatan:** jika username admin sudah ada, proses startup tidak menimpa password yang tersimpan. Mengubah `ADMIN_PASSWORD` di `.env` tidak otomatis mereset password admin yang sudah dibuat.

## Endpoint API

Base URL:

```text
http://localhost:8080
```

### Health Check

| Method | Endpoint  | Autentikasi | Keterangan                         |
| ------ | --------- | ----------- | ---------------------------------- |
| GET    | `/health` | Tidak       | Memeriksa API dan koneksi database |

### Authentication

| Method | Endpoint          | Autentikasi | Keterangan                  |
| ------ | ----------------- | ----------- | --------------------------- |
| POST   | `/api/auth/login` | Tidak       | Login dan mendapatkan JWT   |
| GET    | `/api/auth/me`    | Ya          | Mendapatkan profil pengguna |

### Barang

| Method | Endpoint          | Autentikasi | Keterangan         |
| ------ | ----------------- | ----------- | ------------------ |
| GET    | `/api/barang`     | Ya          | Daftar barang      |
| GET    | `/api/barang/:id` | Ya          | Detail barang      |
| POST   | `/api/barang`     | Ya          | Menambahkan barang |
| PUT    | `/api/barang/:id` | Ya          | Memperbarui barang |
| DELETE | `/api/barang/:id` | Ya          | Soft delete barang |

Semua endpoint barang memerlukan token JWT yang valid.

## Autentikasi

### Login

Request:

```http
POST /api/auth/login
Content-Type: application/json
```

Body:

```json
{
  "username": "admin",
  "password": "PASSWORD_ADMIN_KAMU"
}
```

Contoh respons:

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "access_token": "TOKEN_JWT",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "admin"
    }
  }
}
```

Nilai `expires_in` adalah durasi token dalam detik. Nilai default adalah 3600 detik atau 60 menit.

### Menggunakan token

Sertakan token pada header request:

```http
Authorization: Bearer TOKEN_JWT
```

Jika token tidak dikirim, tidak valid, atau kedaluwarsa, API mengembalikan HTTP `401 Unauthorized`.

## CRUD Barang

### 1. Menambahkan barang

```http
POST /api/barang
Authorization: Bearer TOKEN_JWT
Content-Type: application/json
```

Body:

```json
{
  "sku": "BRG001",
  "nama": "Indomie Goreng",
  "deskripsi": "Mi instan goreng",
  "harga": 3500,
  "stok": 20
}
```

Respons sukses: HTTP `201 Created`.

### 2. Mendapatkan daftar barang

```http
GET /api/barang
Authorization: Bearer TOKEN_JWT
```

Parameter query:

| Parameter    | Default      | Keterangan                         |
| ------------ | ------------ | ---------------------------------- |
| `page`       | `1`          | Nomor halaman, minimal 1           |
| `limit`      | `10`         | Jumlah data per halaman, 1–100     |
| `search`     | kosong       | Pencarian SKU atau nama            |
| `sort_by`    | `created_at` | Kolom pengurutan yang diizinkan    |
| `sort_order` | `desc`       | Arah pengurutan: `asc` atau `desc` |

Contoh:

```http
GET /api/barang?page=1&limit=10&search=indomie&sort_by=nama&sort_order=asc
```

Kolom yang dapat digunakan untuk `sort_by`:

- `id`
- `sku`
- `nama`
- `harga`
- `stok`
- `created_at`
- `updated_at`

Nama kolom lain akan ditolak dengan HTTP `400 Bad Request`.

Contoh respons:

```json
{
  "success": true,
  "message": "Data barang berhasil diambil",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total_items": 0,
      "total_pages": 0
    },
    "filters": {
      "search": "indomie",
      "sort_by": "nama",
      "sort_order": "asc"
    }
  }
}
```

### 3. Mendapatkan detail barang

```http
GET /api/barang/1
Authorization: Bearer TOKEN_JWT
```

ID harus berupa bilangan bulat positif. Jika barang tidak ditemukan atau sudah di-soft-delete, API mengembalikan HTTP `404 Not Found`.

### 4. Memperbarui barang

```http
PUT /api/barang/1
Authorization: Bearer TOKEN_JWT
Content-Type: application/json
```

Body:

```json
{
  "sku": "BRG001",
  "nama": "Indomie Goreng Jumbo",
  "deskripsi": "Kemasan jumbo",
  "harga": 4500,
  "stok": 15
}
```

Endpoint PUT memerlukan semua field pada body, termasuk `harga` dan `stok`. Gunakan angka bulat untuk kedua field tersebut.

Respons sukses: HTTP `200 OK`.

### 5. Menghapus barang

```http
DELETE /api/barang/1
Authorization: Bearer TOKEN_JWT
```

API menggunakan soft delete. Data tidak langsung dihapus secara fisik, melainkan diberi timestamp penghapusan melalui kolom `deleted_at`.

Barang yang sudah dihapus tidak muncul pada daftar normal dan tidak bisa diakses melalui endpoint detail.

## Validasi dan Penanganan Error

| Kondisi                                   |                 HTTP status |
| ----------------------------------------- | --------------------------: |
| Request berhasil                          |                    `200 OK` |
| Barang berhasil dibuat                    |               `201 Created` |
| Input tidak valid                         |           `400 Bad Request` |
| Token tidak valid atau tidak dikirim      |          `401 Unauthorized` |
| Barang tidak ditemukan                    |             `404 Not Found` |
| SKU duplikat                              |              `409 Conflict` |
| Database tidak tersedia pada health check |   `503 Service Unavailable` |
| Kesalahan internal server                 | `500 Internal Server Error` |

Validasi barang:

- SKU wajib diisi, maksimal 50 karakter.
- Nama wajib diisi, maksimal 150 karakter.
- Deskripsi opsional, maksimal 5.000 karakter.
- Harga wajib berupa bilangan bulat dan tidak boleh negatif.
- Stok wajib berupa bilangan bulat dan tidak boleh negatif.
- SKU harus unik, termasuk terhadap barang yang sudah di-soft-delete.

Harga menggunakan integer Rupiah, misalnya `3500` berarti Rp3.500.

## Pengujian dengan PowerShell

Jalankan server terlebih dahulu:

```powershell
go run ./cmd/api
```

Buka jendela PowerShell baru untuk pengujian.

### 1. Login

```powershell
$baseUrl = "http://localhost:8080"

$loginBody = @{
    username = "admin"
    password = "PASSWORD_ADMIN_KAMU"
} | ConvertTo-Json

$login = Invoke-RestMethod `
    -Uri "$baseUrl/api/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $loginBody

$token = $login.data.access_token

$headers = @{
    Authorization = "Bearer $token"
}
```

### 2. Tambah barang

```powershell
$body = @{
    sku = "BRG001"
    nama = "Indomie Goreng"
    deskripsi = "Mi instan goreng"
    harga = 3500
    stok = 20
} | ConvertTo-Json

$result = Invoke-RestMethod `
    -Uri "$baseUrl/api/barang" `
    -Method Post `
    -Headers $headers `
    -ContentType "application/json" `
    -Body $body

$result | ConvertTo-Json -Depth 10

$barangId = $result.data.id
```

### 3. Lihat daftar dan cari barang

```powershell
Invoke-RestMethod `
    -Uri "$baseUrl/api/barang?page=1&limit=10" `
    -Method Get `
    -Headers $headers

Invoke-RestMethod `
    -Uri "$baseUrl/api/barang?search=indomie" `
    -Method Get `
    -Headers $headers
```

### 4. Update barang

```powershell
$updateBody = @{
    sku = "BRG001"
    nama = "Indomie Goreng Jumbo"
    deskripsi = "Kemasan jumbo"
    harga = 4500
    stok = 15
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "$baseUrl/api/barang/$barangId" `
    -Method Put `
    -Headers $headers `
    -ContentType "application/json" `
    -Body $updateBody
```

### 5. Hapus barang

```powershell
Invoke-RestMethod `
    -Uri "$baseUrl/api/barang/$barangId" `
    -Method Delete `
    -Headers $headers
```

Gunakan SKU yang berbeda saat mengulangi tes jika SKU sebelumnya sudah digunakan.

## Keamanan

- Simpan kredensial database dan JWT secret di `.env`.
- Jangan commit file `.env` ke Git.
- Gunakan HTTPS saat API dipublikasikan.
- Gunakan password admin dan database yang kuat.
- Jangan mencatat token JWT atau password ke log.
- Terapkan rate limiting pada login sebelum deployment publik.
- Pertimbangkan pembatasan ukuran request body.
- Gunakan akun MariaDB khusus aplikasi dengan hak akses minimum yang diperlukan.
- Simpan secret production di secret manager atau mekanisme pengelolaan secret yang sesuai dengan lingkungan deployment.

## Catatan Implementasi

- Migrasi otomatis dijalankan pada saat startup melalui GORM `AutoMigrate`.
- Akun admin awal dibuat jika username yang dikonfigurasi belum tersedia.
- Perubahan `ADMIN_PASSWORD` tidak otomatis mengubah password admin yang sudah ada.
- JWT pada versi ini digunakan untuk autentikasi, belum untuk pembagian role atau hak akses.
- Soft delete tidak membebaskan SKU untuk dipakai kembali.
- Endpoint CRUD belum memiliki rate limiting khusus.
- Pengujian pada bagian PowerShell adalah pengujian manual; automated unit dan integration test dapat ditambahkan berikutnya.

## Lisensi

Tentukan lisensi project sesuai kebutuhan sebelum memublikasikannya sebagai repository terbuka.
