package service

import (
	red "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
	jwtx "github.com/Zara1024/OpsNexus/cloudops-server/pkg/jwt"
)

// Services 聚合auth模块业务服务。
type Services struct {
	Auth       *AuthService
	User       *UserService
	Role       *RoleService
	Menu       *MenuService
	Department *DepartmentService
}

// New 创建auth模块服务集合。
func New(db *gorm.DB, redisClient *red.Client, jwtManager *jwtx.Manager) *Services {
	repos := repository.New(db)
	menuService := NewMenuService(repos.Menu)
	return &Services{
		Auth:       NewAuthService(repos, redisClient, jwtManager),
		User:       NewUserService(repos.User),
		Role:       NewRoleService(repos.Role),
		Menu:       menuService,
		Department: NewDepartmentService(repos.Department),
	}
}
