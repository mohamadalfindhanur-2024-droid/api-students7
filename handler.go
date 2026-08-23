package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var students []Student
var nextStudentID = 1

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

func matchesSearch(s Student, keyword string) bool {
	keyword = strings.ToLower(keyword)
	return strings.Contains(strings.ToLower(s.Name), keyword) ||
		strings.Contains(strings.ToLower(s.NIM), keyword)
}

func parseParamID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// 1. GET /api/v1/students (List dengan paginasi, search, sort, filter)
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	// Filter
	var filtered []Student
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.Search != "" && !matchesSearch(s, q.Search) {
			continue
		}
		filtered = append(filtered, s)
	}

	// Sort
	sort.SliceStable(filtered, func(i, j int) bool {
		var isLess bool
		switch q.Sort {
		case "nim":
			isLess = filtered[i].NIM < filtered[j].NIM
		case "name":
			isLess = filtered[i].Name < filtered[j].Name
		case "grade":
			isLess = filtered[i].Grade < filtered[j].Grade
		case "created_at":
			isLess = filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		default:
			isLess = filtered[i].ID < filtered[j].ID
		}

		if q.Order == "desc" {
			return !isLess
		}
		return isLess
	})

	// Pagination
	total := len(filtered)
	totalPages := 0
	if total > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	start := (q.Page - 1) * q.Limit
	if start > total {
		start = total
	}

	end := start + q.Limit
	if end > total {
		end = total
	}

	var paginated []Student
	if total > 0 && start < total {
		paginated = filtered[start:end]
	} else {
		paginated = []Student{}
	}

	return okList(c, "daftar mahasiswa berhasil diambil", paginated, &Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// 2. GET /api/v1/students/:id (Ambil 1 data)
func getStudent(c *fiber.Ctx) error {
	id, valid := parseParamID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	return ok(c, "mahasiswa ditemukan", students[idx])
}

// 3. POST /api/v1/students (Tambah data baru)
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "nilai harus berada di antara 0 dan 100"
	}

	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	// Cek apakah NIM sudah ada (mengembalikan 409 Conflict)
	for _, s := range students {
		if strings.EqualFold(s.NIM, req.NIM) {
			return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
		}
	}

	newStudent := Student{
		ID:        nextStudentID,
		NIM:       req.NIM,
		Name:      req.Name,
		Grade:     req.Grade,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	students = append(students, newStudent)
	nextStudentID++

	location := "/api/v1/students/" + strconv.Itoa(newStudent.ID)
	return created(c, "mahasiswa berhasil ditambahkan", newStudent, location)
}

// 4. PUT /api/v1/students/:id (Mengganti seluruh field)
func replaceStudent(c *fiber.Ctx) error {
	id, valid := parseParamID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "nilai harus berada di antara 0 dan 100"
	}

	// Validasi jika NIM diubah ke NIM orang lain yang sudah terpakai
	for _, s := range students {
		if s.ID != id && strings.EqualFold(s.NIM, req.NIM) {
			errs["nim"] = "NIM sudah terdaftar pada mahasiswa lain"
		}
	}

	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	students[idx].NIM = req.NIM
	students[idx].Name = req.Name
	students[idx].Grade = req.Grade
	students[idx].IsActive = req.IsActive

	return ok(c, "data mahasiswa berhasil diganti seluruhnya", students[idx])
}

// 5. PATCH /api/v1/students/:id (Hanya mengubah field yang dikirim)
func patchStudent(c *fiber.Ctx) error {
	id, valid := parseParamID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	errs := map[string]string{}

	if req.NIM != nil {
		trimmedNIM := strings.TrimSpace(*req.NIM)
		if trimmedNIM == "" {
			errs["nim"] = "NIM tidak boleh kosong"
		} else {
			for _, s := range students {
				if s.ID != id && strings.EqualFold(s.NIM, trimmedNIM) {
					errs["nim"] = "NIM sudah terdaftar pada mahasiswa lain"
				}
			}
		}
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errs["name"] = "nama tidak boleh kosong"
	}

	if req.Grade != nil && (*req.Grade < 0 || *req.Grade > 100) {
		errs["grade"] = "nilai harus berada di antara 0 dan 100"
	}

	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	if req.NIM != nil {
		students[idx].NIM = strings.TrimSpace(*req.NIM)
	}
	if req.Name != nil {
		students[idx].Name = strings.TrimSpace(*req.Name)
	}
	if req.Grade != nil {
		students[idx].Grade = *req.Grade
	}
	if req.IsActive != nil {
		students[idx].IsActive = *req.IsActive
	}

	return ok(c, "data mahasiswa berhasil diperbarui sebagian", students[idx])
}

// 6. DELETE /api/v1/students/:id (Hapus data)
func deleteStudent(c *fiber.Ctx) error {
	id, valid := parseParamID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	students = append(students[:idx], students[idx+1:]...)
	return noContent(c)
}