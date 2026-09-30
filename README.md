# BUMDes Backend API

REST API menggunakan Go, Gin, GORM, dan PostgreSQL.

## Struktur Project

```
cmd/
  api/                entry point aplikasi (HTTP server)
  migrate/            entry point migrasi database
  seed/                entry point seeding data awal
docs/                 hasil generate swagger (swag init), jangan diedit manual
internal/
  config/            load konfigurasi dari env
  database/          koneksi DB, auto-migrate, seeder
  entity/            model database (User, Role, Permission)
  dto/               request/response payload
  repository/        akses data (GORM)
  service/           business logic
  handler/           HTTP handler (Gin) + anotasi swagger
  middleware/         JWT auth middleware
  router/            routing
migrations/          (opsional, jika ingin pakai golang-migrate)
```

## Fitur

- **Auth**: `POST /api/v1/auth/login` — login dengan email & password, mengembalikan JWT.
- **Me**: `GET /api/v1/me` — profil user yang sedang login beserta role-nya.
- **Role Management**: `GET /api/v1/roles`, `GET /api/v1/roles/:id` — daftar role beserta permission yang dimiliki.
- **Account Management**: `GET /api/v1/accounts` (list, dengan filter `?search=` untuk nama/email dan `?name_role=` untuk nama role — isi `all` atau kosongkan untuk semua role), `POST /api/v1/accounts` (create), `PUT /api/v1/accounts/:id` (update), `DELETE /api/v1/accounts/:id` (delete).
- **User Logs**: `GET /api/v1/logs` — riwayat aktivitas user (lihat bagian di bawah).
- **Master Data** (`/api/v1/master-data/*`):
  - **Parking Prices**: `GET /master-data/parking-prices` (list, filter `?search=` untuk nama), `GET /master-data/parking-prices/:id` (detail), `POST /master-data/parking-prices` (create), `PUT /master-data/parking-prices/:id` (update), `DELETE /master-data/parking-prices/:id` (delete). Kolom: `name`, `price`, `status`.
  - **Payment Methods**: `GET /master-data/payment-methods` (list, filter `?search=` untuk method), `GET /master-data/payment-methods/:id` (detail), `POST /master-data/payment-methods` (create), `PUT /master-data/payment-methods/:id` (update), `DELETE /master-data/payment-methods/:id` (delete). Kolom: `method`, `status`.
  - **Revenue Categories**: `GET /master-data/revenue-categories` (list, filter `?search=` mencari di `category_name` maupun `source`), `GET /master-data/revenue-categories/:id` (detail), `POST /master-data/revenue-categories` (create), `PUT /master-data/revenue-categories/:id` (update), `DELETE /master-data/revenue-categories/:id` (delete). Kolom: `category_name`, `source`, `status`.
  - **Ticket Prices**: `GET /master-data/ticket-prices` (list, filter `?search=` untuk type_name), `GET /master-data/ticket-prices/:id` (detail), `POST /master-data/ticket-prices` (create), `PUT /master-data/ticket-prices/:id` (update), `DELETE /master-data/ticket-prices/:id` (delete). Kolom: `id` (UUID), `type_name`, `price`, `status`.
  - **Tenant Types**: `GET /master-data/tenant-types` (list, filter `?search=` untuk type), `GET /master-data/tenant-types/:id` (detail), `POST /master-data/tenant-types` (create), `PUT /master-data/tenant-types/:id` (update), `DELETE /master-data/tenant-types/:id` (delete). Kolom: `id` (UUID), `type`, `description`, `status`.
  - **Attraction Prices**: `GET /master-data/attraction-prices` (list, filter `?search=` untuk nama), `GET /master-data/attraction-prices/:id` (detail), `POST /master-data/attraction-prices` (create), `PUT /master-data/attraction-prices/:id` (update), `DELETE /master-data/attraction-prices/:id` (delete). Kolom: `id` (UUID), `name`, `price`, `status`.
  - **Kios Locations**: `GET /master-data/kios_locations` (list, filter `?search=` untuk kios), `GET /master-data/kios_locations/:id` (detail), `POST /master-data/kios_locations` (create), `PUT /master-data/kios_locations/:id` (update), `DELETE /master-data/kios_locations/:id` (delete). Kolom: `id` (UUID), `kios`, `size`, `price`, `status`, `tenant_id` (nullable, UUID tenant yang sedang menyewa — otomatis disinkron oleh backend saat tenant dibuat/dipindah kios/dihapus, bukan diisi manual lewat form kios).
- **Ticketing** (`/api/v1/ticketing`, transaksi — bukan master data): `GET /ticketing` (list — default filter tanggal hari ini, bisa dikombinasi `?date=YYYY-MM-DD`, `?user_id=`, `?jenis=` untuk type_name, `?payment=` untuk payment_id), `GET /ticketing/:id` (detail), `POST /ticketing` (create), `PUT /ticketing/:id` (update), `DELETE /ticketing/:id` (delete). Kolom: `id` (UUID), `ticket_price_id`, `type_name`, `price`, `payment_id`, `payment_method_name`, `count`, `user_id`, `date`. Saat create/update, kirim `ticket_price_id` & `payment_id` — backend otomatis mengisi `type_name`/`price` dari tarif tiket dan `payment_method_name` dari metode pembayaran; `user_id` dari user yang login, `date` otomatis tanggal hari ini saat create.
- **Parking** (`/api/v1/parkings`, transaksi — bukan master data): `GET /parkings` (list — default filter tanggal hari ini, bisa dikombinasi `?date=YYYY-MM-DD` dan `?lokasi=` untuk location), `GET /parkings/:id` (detail), `POST /parkings` (create), `DELETE /parkings/:id` (delete). Tidak ada update. Kolom: `id` (UUID), `type`, `location`, `count`, `payment_method_id`, `price`, `user_id`, `status`. Saat create, kirim `parking_price_id` (bukan `type`/`price` langsung) — backend otomatis mengisi `type` = nama tarif parkir tersebut dan `price` = tarif × `count`; `user_id` otomatis diisi dari user yang sedang login.
- **Attraction / Wahana** (`/api/v1/attractions`, transaksi — bukan master data): `GET /attractions` (list — default filter tanggal hari ini, bisa dikombinasi `?date=YYYY-MM-DD`, `?attraction_name=`, `?payment_method=` untuk payment_method_name), `GET /attractions/:id` (detail), `POST /attractions` (create), `PUT /attractions/:id` (update), `DELETE /attractions/:id` (delete). Kolom: `id` (UUID), `attraction_price_id`, `attraction_name`, `price`, `count`, `payment_id`, `payment_method_name`, `user_id`, `date`. Saat create/update, kirim `attraction_price_id` & `payment_id` — backend otomatis mengisi `attraction_name`/`price` dari tarif wahana dan `payment_method_name` dari metode pembayaran; `user_id` dari user yang login, `date` otomatis tanggal hari ini saat create.
- **Tenants** (`/api/v1/tenants`, bukan master data): `GET /tenants` (list dengan **pagination** `?page=&limit=`, filter `?search=` untuk nama), `GET /tenants/:id` (detail), `POST /tenants` (create), `PUT /tenants/:id` (update), `DELETE /tenants/:id` (delete). Kolom: `id` (UUID), `name`, `tenant_type_id`, `type`, `phone`, `email` (opsional), `kios_location_id`, `kios_name`, `price`, `status`, `rental_type`, `start_payment` (tanggal patokan mulai pembayaran sewa, format `YYYY-MM-DD`, wajib diisi). Saat create/update, kirim `tenant_type_id` & `kios_location_id` — backend otomatis mengisi `type` dari tipe tenant dan `kios_name` dari lokasi kios yang dipilih.
  - **Search helper untuk form Tenant**: `GET /tenants/tenant-types?q=&limit=` (cari tenant type berdasarkan kolom `type`, default `limit=10`), `GET /tenants/kios_locations?q=&limit=` (cari lokasi kios berdasarkan kolom `kios`, default `limit=10`) — dipakai untuk isi dropdown/autocomplete `tenant_type_id` dan `kios_location_id` saat mengisi form tenant.
- **Kios Available** (`/api/v1/kios_availables`): `GET /kios_availables` — menampilkan seluruh kios dari tabel `kios_locations` beserta status ketersediaan: `"Terisi"` (ada `tenant_id`) atau `"Tersedia"` (kosong). Status ini otomatis mengikuti `kios_locations.tenant_id`, yang disinkron backend setiap kali tenant dibuat, pindah kios, atau dihapus. Field `price`: kalau `"Terisi"` diambil dari `price` tenant yang menyewa, kalau `"Tersedia"` diambil dari `price` kios itu sendiri.

Semua endpoint selain login membutuhkan header `Authorization: Bearer <token>`.

> **Catatan**: primary key tabel `users` menggunakan UUID (bukan auto-increment), di-generate di aplikasi (bukan default DB). Semua field `id`/`user_id` yang merujuk ke user (JWT `user_id`, `user.id` di response, `user_logs.user_id`) berupa string UUID. Primary key tabel lain (roles, permissions, master data, dll) tetap auto-increment biasa.

## User Logs

Setiap aktivitas user pada endpoint yang butuh auth otomatis dicatat ke tabel `user_logs`:

- **view** — otomatis tercatat untuk setiap request `GET` yang berhasil (mis. buka `/me`, `/roles`, `/accounts`).
- **create** — tercatat saat `POST /accounts`, kolom `data` berisi `{"input": {...}}` (field password diredact jadi `***`).
- **update** — tercatat saat `PUT /accounts/:id`, kolom `data` berisi `{"old": {...}, "new": {...}}`.
- **delete** — tercatat saat `DELETE /accounts/:id`, kolom `data` bernilai `null`.

Setiap baris log juga menyimpan `user_id` pelaku, `module`, `method`, `path`, `ip_address`, dan `user_agent`. Lihat riwayatnya lewat:

```
GET /api/v1/logs?user_id=&action=&module=&page=1&limit=20
```

Semua filter (`user_id`, `action`, `module`) opsional.

## API Documentation (Swagger)

Setiap endpoint didokumentasikan dengan Swagger. Setelah server jalan, buka:

```
http://localhost:8080/swagger/index.html
```

Untuk memanggil endpoint yang butuh auth langsung dari Swagger UI, klik tombol **Authorize** lalu isi `Bearer <token>` (token didapat dari response `POST /auth/login`).

Jika ada perubahan pada anotasi swagger di handler, regenerate docs dengan:

```bash
make swag
```

(butuh `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@latest`)

## Menjalankan

1. Copy `.env.example` menjadi `.env` dan sesuaikan konfigurasi database.
2. Jalankan PostgreSQL (bisa pakai `docker compose up -d`).
3. Jalankan migrasi database:

```bash
make migrate
```

4. Jalankan seeder untuk data awal:

```bash
make seed
```

5. Jalankan aplikasi:

```bash
make run
```

Migrasi dan seeding dipisah dari proses start API — jalankan manual setiap ada perubahan schema (`make migrate`) atau saat butuh reset data awal (`make seed`, aman dijalankan berulang kali karena idempotent).

Data yang dibuat oleh seeder:

- Role: `superadmin` (semua permission), `admin` (semua permission), `member` (view-only)
- Permission: `user.view`, `user.manage`, `role.view`, `role.manage`
- User default:
  - `superadmin@bumdes.local` / `password123` (role `superadmin`)
  - `admin@bumdes.local` / `password123` (role `admin`)

## Build

```bash
make build
```
