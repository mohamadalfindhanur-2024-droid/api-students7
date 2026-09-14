package service_test

import (
	"testing"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
)

func TestCanAccessStudent(t *testing.T) {
	permsMap := map[string][]string{
		"admin": {"student:list", "student:read:any", "student:create", "student:update:any", "student:delete"},
		"staff": {"student:list", "student:read:any", "student:create"},
		"user":  {},
	}
	perms := helper.NewPermissionSet(permsMap)

	userOwner := model.AuthUser{UserID: 10, Username: "pemilik", Role: "user"}
	userNonOwner := model.AuthUser{UserID: 20, Username: "orang_lain", Role: "user"}
	staffUser := model.AuthUser{UserID: 30, Username: "staf1", Role: "staff"}
	adminUser := model.AuthUser{UserID: 40, Username: "admin1", Role: "admin"}

	const studentOwnerID = 10

	tests := []struct {
		name          string
		current       model.AuthUser
		ownerID       int
		anyPermission string
		expected      bool
	}{
		{
			name:          "User biasa mengakses data miliknya sendiri (Ownership) -> Diizinkan",
			current:       userOwner,
			ownerID:       studentOwnerID,
			anyPermission: "student:read:any",
			expected:      true,
		},
		{
			name:          "User biasa mengakses data milik orang lain (Bukan pemilik & tidak punya :any) -> Ditolak",
			current:       userNonOwner,
			ownerID:       studentOwnerID,
			anyPermission: "student:read:any",
			expected:      false,
		},
		{
			name:          "Staff membaca data milik user lain (Punya permission student:read:any) -> Diizinkan",
			current:       staffUser,
			ownerID:       studentOwnerID,
			anyPermission: "student:read:any",
			expected:      true,
		},
		{
			name:          "Staff mengedit data milik user lain (Tidak punya student:update:any) -> Ditolak",
			current:       staffUser,
			ownerID:       studentOwnerID,
			anyPermission: "student:update:any",
			expected:      false,
		},
		{
			name:          "Admin mengedit data milik user lain (Punya permission student:update:any) -> Diizinkan",
			current:       adminUser,
			ownerID:       studentOwnerID,
			anyPermission: "student:update:any",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CanAccessStudent(tt.current, tt.ownerID, perms, tt.anyPermission)
			if result != tt.expected {
				t.Errorf("CanAccessStudent() = %v; expected %v", result, tt.expected)
			}
		})
	}
}
