package config

type (
	SuperAdminConfig struct {
		Username string `mapstructure:"username" description:"用户名"`
		Email    string `mapstructure:"email" description:"邮箱"`
		Password string `mapstructure:"password" description:"密码"`
		RoleID   string `mapstructure:"role_id" description:"角色ID"`
	}
)
