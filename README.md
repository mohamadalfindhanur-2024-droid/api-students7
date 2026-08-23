# Tugas 2 Praktikum Pemrograman Backend Lanjut - REST API Students

Repositori ini berisi implementasi REST API untuk manajemen data mahasiswa menggunakan Go dan framework Fiber. Data disimpan secara in-memory (slice).

## Dokumentasi Endpoint (Kontrak API)

### 1. GET /api/v1/students
Mengambil daftar mahasiswa dengan fitur pagination, search, sort, dan filter.
- Query Params:
  - page (default 1)
  - limit (default 10, max 100)
  - search (mencari keyword pada nama/NIM)
  - sort (id, nim, name, grade, created_at)
  - order (asc / desc)
  - is_active (true / false)
- Response: 200 OK
```json
{
  "success": true,
  "message": "daftar mahasiswa berhasil diambil",
  "data": [
    {
      "id": 1,
      "nim": "5025211001",
      "name": "Alfin",
      "grade": 90,
      "is_active": true,
      "created_at": "2026-08-23T17:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}