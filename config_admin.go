package config

type SuperAdminConfig struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
	RoleID   string `mapstructure:"role_id"`
}
