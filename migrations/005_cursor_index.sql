-- Index untuk Keyset / Cursor Pagination pada tabel users (Modul 7)
CREATE INDEX IF NOT EXISTS users_created_at_id_desc_idx
    ON users (created_at DESC, id DESC);
