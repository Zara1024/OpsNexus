package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims 定义业务JWT声明。
type CustomClaims struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	RoleCodes   []string `json:"role_codes"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"`
	jwtv5.RegisteredClaims
}

// TokenPair 定义登录返回的双Token。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Manager 管理JWT签发与解析。
type Manager struct {
	secret            []byte
	issuer            string
	accessExpire      time.Duration
	refreshExpire     time.Duration
	autoRefreshWindow time.Duration
}

// NewManager 创建JWT管理器。
func NewManager(secret, issuer string, accessExpire, refreshExpire time.Duration) *Manager {
	if accessExpire <= 0 {
		accessExpire = 30 * time.Minute
	}
	if refreshExpire <= 0 {
		refreshExpire = 7 * 24 * time.Hour
	}
	return &Manager{
		secret:            []byte(secret),
		issuer:            issuer,
		accessExpire:      accessExpire,
		refreshExpire:     refreshExpire,
		autoRefreshWindow: 10 * time.Minute,
	}
}

// GenerateTokenPair 生成Access与Refresh Token。
func (m *Manager) GenerateTokenPair(userID int64, username string, roleCodes []string, permissions []string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := CustomClaims{
		UserID:      userID,
		Username:    username,
		RoleCodes:   roleCodes,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   username,
			ID:        uuid.NewString(),
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(m.accessExpire)),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}

	refreshClaims := CustomClaims{
		UserID:      userID,
		Username:    username,
		RoleCodes:   roleCodes,
		Permissions: permissions,
		TokenType:   "refresh",
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   username,
			ID:        uuid.NewString(),
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(m.refreshExpire)),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}

	accessToken, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpire.Seconds()),
	}, nil
}

// ParseAccessToken 解析Access Token。
func (m *Manager) ParseAccessToken(token string) (*CustomClaims, error) {
	claims, err := m.parse(token)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, errors.New("token类型错误")
	}
	return claims, nil
}

// ParseRefreshToken 解析Refresh Token。
func (m *Manager) ParseRefreshToken(token string) (*CustomClaims, error) {
	claims, err := m.parse(token)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("token类型错误")
	}
	return claims, nil
}

// NeedAutoRefresh 判断Access Token是否需要自动续期。
func (m *Manager) NeedAutoRefresh(claims *CustomClaims) bool {
	if claims == nil || claims.RegisteredClaims.ExpiresAt == nil {
		return false
	}
	return time.Until(claims.RegisteredClaims.ExpiresAt.Time) <= m.autoRefreshWindow
}

// AccessTTL 返回Access Token剩余时长。
func (m *Manager) AccessTTL(claims *CustomClaims) time.Duration {
	if claims == nil || claims.RegisteredClaims.ExpiresAt == nil {
		return 0
	}
	ttl := time.Until(claims.RegisteredClaims.ExpiresAt.Time)
	if ttl < 0 {
		return 0
	}
	return ttl
}

func (m *Manager) parse(token string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	parsed, err := jwtv5.ParseWithClaims(token, claims, func(token *jwtv5.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("token无效")
	}
	return claims, nil
}
