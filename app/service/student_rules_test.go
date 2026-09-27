package service_test

import (
	"testing"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
)

func TestApplyPatchStudent(t *testing.T) {
	initial := model.Student{
		ID:       1,
		NIM:      "434241132",
		Name:     "Mohamad Alfin",
		Grade:    3.85,
		IsActive: true,
	}

	newName := "Mohamad Alfin D. R."
	patchReq := model.PatchStudentRequest{
		Name: &newName,
	}

	updated := service.ApplyPatchStudent(initial, patchReq)
	if updated.Name != "Mohamad Alfin D. R." {
		t.Errorf("ApplyPatchStudent() name = %q, expected %q", updated.Name, "Mohamad Alfin D. R.")
	}
	if updated.NIM != "434241132" {
		t.Errorf("ApplyPatchStudent() NIM modified unexpectedly: %q", updated.NIM)
	}
}

func TestIsEmptyPatchStudent(t *testing.T) {
	emptyReq := model.PatchStudentRequest{}
	if !service.IsEmptyPatchStudent(emptyReq) {
		t.Errorf("IsEmptyPatchStudent() for empty request expected true, got false")
	}

	val := "434241132"
	nonEmptyReq := model.PatchStudentRequest{NIM: &val}
	if service.IsEmptyPatchStudent(nonEmptyReq) {
		t.Errorf("IsEmptyPatchStudent() for non-empty request expected false, got true")
	}
}

func TestValidateCreateStudent(t *testing.T) {
	tests := []struct {
		name        string
		req         model.CreateStudentRequest
		expectedErr bool
		errorField  string
	}{
		{
			name: "Data mahasiswa valid",
			req: model.CreateStudentRequest{
				NIM:   "434241132",
				Name:  "Mohamad Alfin",
				Grade: 3.85,
			},
			expectedErr: false,
		},
		{
			name: "NIM format bukan angka / terlalu pendek",
			req: model.CreateStudentRequest{
				NIM:   "abc12",
				Name:  "Mohamad Alfin",
				Grade: 3.85,
			},
			expectedErr: true,
			errorField:  "nim",
		},
		{
			name: "Nama terlalu pendek",
			req: model.CreateStudentRequest{
				NIM:   "434241132",
				Name:  "Al",
				Grade: 3.85,
			},
			expectedErr: true,
			errorField:  "name",
		},
		{
			name: "Grade melebihi 4.00",
			req: model.CreateStudentRequest{
				NIM:   "434241132",
				Name:  "Mohamad Alfin",
				Grade: 5.00,
			},
			expectedErr: true,
			errorField:  "grade",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := helper.ValidateStruct(tt.req)
			if tt.expectedErr && len(errs) == 0 {
				t.Errorf("ValidateStruct() expected error on field %s, got none", tt.errorField)
			}
			if !tt.expectedErr && len(errs) > 0 {
				t.Errorf("ValidateStruct() expected no error, got: %v", errs)
			}
		})
	}
}

func TestValidatePatchStudentOmitnil(t *testing.T) {
	// Tugas D.2: {"name":""} DITOLAK (422), sedangkan body tanpa name diterima
	emptyName := ""
	invalidPatch := model.PatchStudentRequest{
		Name: &emptyName,
	}
	errs := helper.ValidateStruct(invalidPatch)
	if len(errs) == 0 || errs["name"] == "" {
		t.Errorf("ValidateStruct() on empty pointer string expected validation error, got: %v", errs)
	}

	validPatchWithoutName := model.PatchStudentRequest{
		Name: nil, // tidak dikirim -> omitnil lolos
	}
	errs2 := helper.ValidateStruct(validPatchWithoutName)
	if len(errs2) > 0 {
		t.Errorf("ValidateStruct() on nil pointer expected to pass, got: %v", errs2)
	}
}
