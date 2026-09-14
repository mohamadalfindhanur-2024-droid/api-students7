package helper

import "sort"

// PermissionSet menyimpan pemetaan role ke daftar permission miliknya di memori.
type PermissionSet struct {
	byRole map[string]map[string]struct{}
}

// NewPermissionSet mengubah hasil query database menjadi struktur hash set map di memori.
func NewPermissionSet(raw map[string][]string) *PermissionSet {
	byRole := make(map[string]map[string]struct{}, len(raw))
	for role, permissions := range raw {
		set := make(map[string]struct{}, len(permissions))
		for _, permission := range permissions {
			set[permission] = struct{}{}
		}
		byRole[role] = set
	}
	return &PermissionSet{byRole: byRole}
}

// Can menjawab apakah role ini memiliki permission tertentu.
// Menerapkan prinsip Fail Closed: SELALU mengembalikan false jika nil, role tidak dikenal, atau permission tidak ada.
func (p *PermissionSet) Can(role, permission string) bool {
	if p == nil {
		return false
	}
	permissions, ok := p.byRole[role]
	if !ok {
		return false
	}
	_, granted := permissions[permission]
	return granted
}

// PermissionsOf mengembalikan seluruh permission milik sebuah role, terurut alfabetis.
func (p *PermissionSet) PermissionsOf(role string) []string {
	result := []string{}
	if p == nil {
		return result
	}
	for permission := range p.byRole[role] {
		result = append(result, permission)
	}
	sort.Strings(result)
	return result
}

// KnownRoles mengembalikan daftar seluruh role yang diakui sistem, terurut.
func (p *PermissionSet) KnownRoles() []string {
	result := []string{}
	if p == nil {
		return result
	}
	for role := range p.byRole {
		result = append(result, role)
	}
	sort.Strings(result)
	return result
}

// IsKnownRole memeriksa apakah nama role terdaftar di sistem.
func (p *PermissionSet) IsKnownRole(role string) bool {
	if p == nil {
		return false
	}
	_, ok := p.byRole[role]
	return ok
}
