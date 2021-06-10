package config

const (
	RoleAdmin = "ADMIN"
	RoleSEC   = "SECURITY"
	RoleHSSE  = "HSSE"
)

func GetRolesAvailable() []string {
	return []string{RoleAdmin, RoleSEC, RoleHSSE}
}
