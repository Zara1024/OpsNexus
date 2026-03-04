package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	red "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	jwtx "github.com/Zara1024/OpsNexus/cloudops-server/pkg/jwt"
)

const (
	loginFailMax         = 5
	loginFailWindow      = 15 * time.Minute
	captchaTTL           = 3 * time.Minute
	captchaKeyPrefix     = "auth:captcha:"
	loginFailKeyPrefix   = "auth:login_fail:"
	tokenBlacklistPrefix = "auth:blacklist:"
)

// AuthService 提供认证相关能力。
type AuthService struct {
	repos  *repository.Repositories
	redis  *red.Client
	jwtMgr *jwtx.Manager
}

// LoginRequest 定义登录请求。
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

// LoginResponse 定义登录响应。
type LoginResponse struct {
	TokenPair *jwtx.TokenPair `json:"token"`
}

// UserInfo 定义当前用户信息。
type UserInfo struct {
	UserID      int64                `json:"user_id"`
	Username    string               `json:"username"`
	Nickname    string               `json:"nickname"`
	RoleCodes   []string             `json:"role_codes"`
	Permissions []string             `json:"permissions"`
	Menus       []model.MenuTreeNode `json:"menus"`
}

// CaptchaResponse 定义验证码响应。
type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

// NewAuthService 创建认证服务。
func NewAuthService(repos *repository.Repositories, redisClient *red.Client, jwtMgr *jwtx.Manager) *AuthService {
	return &AuthService{repos: repos, redis: redisClient, jwtMgr: jwtMgr}
}

// GetCaptcha 生成数字验证码。
func (s *AuthService) GetCaptcha(ctx context.Context) (*CaptchaResponse, error) {
	id := randomString(16)
	code := randomDigits(4)
	key := captchaKeyPrefix + id
	if err := s.redis.Set(ctx, key, code, captchaTTL).Err(); err != nil {
		return nil, err
	}
	return &CaptchaResponse{CaptchaID: id, Captcha: code}, nil
}

// Login 执行登录流程。
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return nil, apperrors.ErrValidation
	}

	if locked, err := s.isLocked(ctx, username); err != nil {
		return nil, err
	} else if locked {
		return nil, &apperrors.AppError{Code: 10007, HTTPCode: 400, Message: "登录失败次数过多，请15分钟后再试"}
	}

	if ok, err := s.verifyCaptcha(ctx, req.CaptchaID, req.Captcha); err != nil {
		return nil, err
	} else if !ok {
		return nil, &apperrors.AppError{Code: 10008, HTTPCode: 400, Message: "验证码错误或已过期"}
	}

	user, err := s.repos.User.GetByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			_ = s.increaseLoginFail(ctx, username)
			return nil, &apperrors.AppError{Code: 10002, HTTPCode: 400, Message: "用户名或密码错误"}
		}
		return nil, err
	}

	if user.Status != 1 {
		return nil, &apperrors.AppError{Code: 10009, HTTPCode: 403, Message: "用户已被禁用"}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		_ = s.increaseLoginFail(ctx, username)
		return nil, &apperrors.AppError{Code: 10002, HTTPCode: 400, Message: "用户名或密码错误"}
	}
	_ = s.resetLoginFail(ctx, username)

	roleIDs, _ := s.repos.User.GetRoleIDsByUserID(user.ID)
	roleCodes, _ := s.repos.User.GetRoleCodesByUserID(user.ID)
	permissions, _ := s.repos.Menu.ListPermissionsByRoleIDs(roleIDs)

	tokenPair, err := s.jwtMgr.GenerateTokenPair(user.ID, user.Username, roleCodes, permissions)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{TokenPair: tokenPair}, nil
}

// Logout 将Access Token加入黑名单。
func (s *AuthService) Logout(ctx context.Context, accessToken string, claims *jwtx.CustomClaims) error {
	if accessToken == "" {
		return nil
	}
	ttl := s.jwtMgr.AccessTTL(claims)
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return s.redis.Set(ctx, tokenBlacklistPrefix+accessToken, 1, ttl).Err()
}

// IsBlacklisted 判断Token是否在黑名单。
func (s *AuthService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	n, err := s.redis.Exists(ctx, tokenBlacklistPrefix+token).Result()
	return n > 0, err
}

// Refresh 刷新Token。
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*jwtx.TokenPair, error) {
	claims, err := s.jwtMgr.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	user, err := s.repos.User.GetByID(claims.UserID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if user.Status != 1 {
		return nil, apperrors.ErrForbidden
	}

	roleIDs, _ := s.repos.User.GetRoleIDsByUserID(user.ID)
	roleCodes, _ := s.repos.User.GetRoleCodesByUserID(user.ID)
	permissions, _ := s.repos.Menu.ListPermissionsByRoleIDs(roleIDs)

	return s.jwtMgr.GenerateTokenPair(user.ID, user.Username, roleCodes, permissions)
}

// GetUserInfo 查询当前用户信息与权限。
func (s *AuthService) GetUserInfo(userID int64) (*UserInfo, error) {
	user, err := s.repos.User.GetByID(userID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	roleIDs, _ := s.repos.User.GetRoleIDsByUserID(user.ID)
	roleCodes, _ := s.repos.User.GetRoleCodesByUserID(user.ID)
	permissions, _ := s.repos.Menu.ListPermissionsByRoleIDs(roleIDs)
	menus, _ := s.repos.Menu.ListByRoleIDs(roleIDs)

	return &UserInfo{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		RoleCodes:   roleCodes,
		Permissions: permissions,
		Menus:       buildMenuTree(menus),
	}, nil
}

func (s *AuthService) verifyCaptcha(ctx context.Context, captchaID, captcha string) (bool, error) {
	if captchaID == "" || captcha == "" {
		return false, nil
	}
	key := captchaKeyPrefix + captchaID
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == red.Nil {
			return false, nil
		}
		return false, err
	}
	_ = s.redis.Del(ctx, key).Err()
	return strings.EqualFold(strings.TrimSpace(stored), strings.TrimSpace(captcha)), nil
}

func (s *AuthService) isLocked(ctx context.Context, username string) (bool, error) {
	count, err := s.redis.Get(ctx, loginFailKeyPrefix+username).Int()
	if err != nil && err != red.Nil {
		return false, err
	}
	return count >= loginFailMax, nil
}

func (s *AuthService) increaseLoginFail(ctx context.Context, username string) error {
	key := loginFailKeyPrefix + username
	pipe := s.redis.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, loginFailWindow)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *AuthService) resetLoginFail(ctx context.Context, username string) error {
	return s.redis.Del(ctx, loginFailKeyPrefix+username).Err()
}

func randomDigits(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	for i := range buf {
		buf[i] = '0' + (buf[i] % 10)
	}
	return string(buf)
}

func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)[:n]
}

func buildMenuTree(menus []model.Menu) []model.MenuTreeNode {
	nodeMap := make(map[int64]*model.MenuTreeNode, len(menus))
	roots := make([]model.MenuTreeNode, 0)

	for _, m := range menus {
		copyMenu := m
		nodeMap[m.ID] = &model.MenuTreeNode{
			ID:         copyMenu.ID,
			ParentID:   copyMenu.ParentID,
			MenuName:   copyMenu.MenuName,
			RoutePath:  copyMenu.RoutePath,
			Component:  copyMenu.Component,
			Icon:       copyMenu.Icon,
			SortOrder:  copyMenu.SortOrder,
			MenuType:   copyMenu.MenuType,
			Permission: copyMenu.Permission,
			Children:   make([]model.MenuTreeNode, 0),
		}
	}

	for _, node := range nodeMap {
		if node.ParentID == 0 {
			roots = append(roots, *node)
			continue
		}
		if parent, ok := nodeMap[node.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}
	return roots
}
