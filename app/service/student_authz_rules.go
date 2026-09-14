package service

import (
	"api-students/app/model"
	"api-students/helper"
)

// CanAccessStudent memutuskan apakah user berhak mengakses / memodifikasi data mahasiswa tertentu.
// Dua jalur akses yang diizinkan:
// 1. Kepemilikan (Ownership): User adalah pemilik data tersebut (current.UserID == ownerID).
// 2. Permission RBAC: Role user memiliki permission cakupan ':any' yang sesuai.
func CanAccessStudent(
	current model.AuthUser,
	ownerID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == ownerID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}
