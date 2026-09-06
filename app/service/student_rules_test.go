package service

import (
	"testing"

	"api-students/app/model"
)

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
		{5, 0, 0},
	}
	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{
		ID:       1,
		NIM:      "5025211001",
		Name:     "Alfin",
		Grade:    3.75,
		IsActive: true,
	}

	inactive := false
	newName := "Alfin Updated"
	newGrade := 3.90

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{
		Name:     &newName,
		Grade:    &newGrade,
		IsActive: &inactive,
	})

	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error: %v", errs)
	}
	if result.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}
	if result.Name != "Alfin Updated" {
		t.Errorf("name diharapkan 'Alfin Updated', didapat '%s'", result.Name)
	}
	if result.Grade != 3.90 {
		t.Errorf("grade diharapkan 3.90, didapat %f", result.Grade)
	}
	if result.NIM != "5025211001" {
		t.Error("field nim yang tidak dikirim seharusnya tidak berubah")
	}

	// Test validation failure on invalid patch
	emptyNIM := "   "
	invalidGrade := 4.50
	_, invalidErrs := ApplyPatch(initial, model.PatchStudentRequest{
		NIM:   &emptyNIM,
		Grade: &invalidGrade,
	})
	if len(invalidErrs) != 2 {
		t.Errorf("diharapkan 2 error validasi, didapat %d", len(invalidErrs))
	}
}

func TestValidateCreate(t *testing.T) {
	// Kasus valid
	validReq := model.CreateStudentRequest{
		NIM:   "5025211002",
		Name:  "Budi Santoso",
		Grade: 3.50,
	}
	if errs := ValidateCreate(validReq); len(errs) != 0 {
		t.Errorf("data valid tidak seharusnya menghasilkan error: %v", errs)
	}

	// Kasus invalid
	invalidReq := model.CreateStudentRequest{
		NIM:   "",
		Name:  "   ",
		Grade: 4.50,
	}
	errs := ValidateCreate(invalidReq)
	if len(errs) != 3 {
		t.Errorf("diharapkan 3 error validasi, didapat %d", len(errs))
	}
	if errs["nim"] != "wajib diisi" {
		t.Errorf("pesan error nim tidak sesuai: %s", errs["nim"])
	}
	if errs["name"] != "wajib diisi" {
		t.Errorf("pesan error name tidak sesuai: %s", errs["name"])
	}
	if errs["grade"] != "nilai harus berada di antara 0.00 dan 4.00" {
		t.Errorf("pesan error grade tidak sesuai: %s", errs["grade"])
	}
}

func TestValidateReplace(t *testing.T) {
	// Kasus valid
	validReq := model.ReplaceStudentRequest{
		NIM:      "5025211003",
		Name:     "Citra",
		Grade:    3.80,
		IsActive: true,
	}
	if errs := ValidateReplace(validReq); len(errs) != 0 {
		t.Errorf("data valid tidak seharusnya menghasilkan error: %v", errs)
	}

	// Kasus invalid
	invalidReq := model.ReplaceStudentRequest{
		NIM:      "",
		Name:     "",
		Grade:    -1.0,
		IsActive: true,
	}
	errs := ValidateReplace(invalidReq)
	if len(errs) != 3 {
		t.Errorf("diharapkan 3 error validasi, didapat %d", len(errs))
	}
}

func TestIsEmptyPatch(t *testing.T) {
	if !IsEmptyPatch(model.PatchStudentRequest{}) {
		t.Error("struct kosong seharusnya mengembalikan true")
	}
	name := "Alfin"
	if IsEmptyPatch(model.PatchStudentRequest{Name: &name}) {
		t.Error("struct dengan data seharusnya mengembalikan false")
	}
}
