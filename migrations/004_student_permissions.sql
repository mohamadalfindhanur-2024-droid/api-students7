-- ---------------------------------------------------------------
-- permissions untuk entity students (Modul 6 Tugas Mandiri C.1)
-- ---------------------------------------------------------------
INSERT INTO permissions (name, description) VALUES
    ('student:list', 'Melihat daftar seluruh mahasiswa'),
    ('student:read:any', 'Melihat data mahasiswa milik siapa pun'),
    ('student:create', 'Mendaftarkan data mahasiswa baru'),
    ('student:update:any', 'Mengubah data mahasiswa milik siapa pun'),
    ('student:delete', 'Menghapus data mahasiswa')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create')
ON CONFLICT DO NOTHING;

-- Tambahkan kolom owner_id pada tabel students sebagai dasar pemeriksaan kepemilikan
ALTER TABLE students
    ADD COLUMN IF NOT EXISTS owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Bila ada data lama di tabel students, hubungkan ke user pertama yang ada
UPDATE students SET owner_id = (SELECT id FROM users ORDER BY id LIMIT 1) WHERE owner_id IS NULL;

CREATE INDEX IF NOT EXISTS students_owner_id_idx ON students (owner_id);
