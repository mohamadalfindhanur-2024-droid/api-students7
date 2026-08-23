package main

import "time"

// Student entitas utama mahasiswa
type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateStudentRequest untuk POST (semua field wajib ada)
type CreateStudentRequest struct {
	NIM   string  `json:"nim"`
	Name  string  `json:"name"`
	Grade float64 `json:"grade"`
}

// ReplaceStudentRequest untuk PUT (ganti seluruh data, semua wajib)
type ReplaceStudentRequest struct {
	NIM      string  `json:"nim"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// PatchStudentRequest untuk PATCH (ubah sebagian, pakai pointer agar bisa tahu field yang tidak dikirim)
type PatchStudentRequest struct {
	NIM      *string  `json:"nim,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Grade    *float64 `json:"grade,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// WebResponse amplop baku untuk semua respons JSON
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// Meta informasi paginasi
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListQuery menampung parameter query string
type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}