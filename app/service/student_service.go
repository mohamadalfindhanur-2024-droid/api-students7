package service

import (
	"errors"
	"strconv"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	repo  repository.StudentRepository
	perms *helper.PermissionSet
}

func NewStudentService(repo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{repo: repo, perms: perms}
}

func (s *StudentService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	format, err := helper.Negotiate(c, helper.FormatJSON, helper.FormatCSV)
	if err != nil {
		return err
	}

	q, err := helper.ParseCursorQuery(c)
	if err != nil {
		return err
	}

	rows, err := s.repo.FindAfterCursor(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}

	if format == helper.FormatCSV {
		return helper.WriteStudentsCSV(c, rows)
	}

	hasMore := len(rows) > q.Limit
	if hasMore {
		rows = rows[:q.Limit]
	}

	meta := &model.CursorMeta{Limit: q.Limit, HasMore: hasMore}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		meta.NextCursor = helper.EncodeCursor(last.CreatedAt, last.ID)
	}

	return helper.SuccessCursor(c, "daftar student berhasil diambil", rows, meta)
}

func (s *StudentService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("id harus berupa bilangan bulat positif")
	}

	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return helper.Internal(err)
	}

	if !CanAccessStudent(current, student.OwnerID, s.perms, "student:read:any") {
		return helper.Forbidden("tidak berhak mengakses data mahasiswa milik user lain")
	}

	return helper.Success(c, fiber.StatusOK, "data mahasiswa berhasil diambil", student)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	student := model.Student{
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    req.Grade,
		IsActive: true,
		OwnerID:  current.UserID,
	}

	created, err := s.repo.Create(ctx, student)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM sudah terdaftar")
		}
		return helper.Internal(err)
	}

	location := "/api/v1/students/" + strconv.Itoa(created.ID)
	return helper.Created(c, "data mahasiswa berhasil dibuat", created, location)
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("id harus berupa bilangan bulat positif")
	}

	currentStudent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return helper.Internal(err)
	}

	if !CanAccessStudent(current, currentStudent.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah data mahasiswa milik user lain")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	currentStudent.NIM = strings.TrimSpace(req.NIM)
	currentStudent.Name = strings.TrimSpace(req.Name)
	currentStudent.Grade = req.Grade
	currentStudent.IsActive = req.IsActive

	updated, err := s.repo.Update(ctx, currentStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM sudah terdaftar")
		}
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "data mahasiswa berhasil diperbarui", updated)
}

func (s *StudentService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("id harus berupa bilangan bulat positif")
	}

	currentStudent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return helper.Internal(err)
	}

	if !CanAccessStudent(current, currentStudent.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah data mahasiswa milik user lain")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if IsEmptyPatchStudent(req) {
		return helper.BadRequest("setidaknya satu field harus diisi pada PATCH")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	updatedStudent := ApplyPatchStudent(currentStudent, req)
	updated, err := s.repo.Update(ctx, updatedStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM sudah terdaftar")
		}
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "data mahasiswa berhasil diperbarui", updated)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("id harus berupa bilangan bulat positif")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return helper.Internal(err)
	}

	return helper.NoContent(c)
}
