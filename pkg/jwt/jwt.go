// Package jwt 封装 Access / Refresh 双令牌的签发与校验。
package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 令牌校验错误。
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrWrongType    = errors.New("wrong token type")
)

// TokenType 区分令牌用途，防止 refresh token 被当作 access token 使用。
type TokenType string

// 令牌类型枚举。
const (
	TypeAccess  TokenType = "access"
	TypeRefresh TokenType = "refresh"
)

// Claims 自定义载荷。
type Claims struct {
	UserID    int64     `json:"uid"`
	Username  string    `json:"uname,omitempty"`
	TokenType TokenType `json:"typ"`
	// FamilyID 标识一条刷新链，用于刷新令牌轮转与复用检测，仅 refresh 令牌携带。
	FamilyID string `json:"fid,omitempty"`
	jwtlib.RegisteredClaims
}

// Token 签发结果。
type Token struct {
	Value     string
	TokenID   string
	FamilyID  string
	ExpiresAt time.Time
}

// Options 管理器配置。
type Options struct {
	Issuer        string
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// Manager 令牌管理器。
type Manager struct {
	issuer        string
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewManager 创建令牌管理器。
func NewManager(opts Options) *Manager {
	return &Manager{
		issuer:        opts.Issuer,
		accessSecret:  []byte(opts.AccessSecret),
		refreshSecret: []byte(opts.RefreshSecret),
		accessTTL:     opts.AccessTTL,
		refreshTTL:    opts.RefreshTTL,
	}
}

// AccessTTL 返回访问令牌有效期。
func (m *Manager) AccessTTL() time.Duration { return m.accessTTL }

// RefreshTTL 返回刷新令牌有效期。
func (m *Manager) RefreshTTL() time.Duration { return m.refreshTTL }

// GenerateAccess 签发访问令牌。
func (m *Manager) GenerateAccess(userID int64, username string) (*Token, error) {
	return m.generate(userID, username, "", TypeAccess, m.accessSecret, m.accessTTL)
}

// GenerateRefresh 签发刷新令牌，familyID 为空时自动生成新的刷新链。
func (m *Manager) GenerateRefresh(userID int64, familyID string) (*Token, error) {
	if familyID == "" {
		familyID = uuid.NewString()
	}
	return m.generate(userID, "", familyID, TypeRefresh, m.refreshSecret, m.refreshTTL)
}

func (m *Manager) generate(
	userID int64, username, familyID string,
	typ TokenType, secret []byte, ttl time.Duration,
) (*Token, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	tokenID := uuid.NewString()

	claims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: typ,
		FamilyID:  familyID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ID:        tokenID,
			Issuer:    m.issuer,
			Subject:   fmt.Sprint(userID),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &Token{Value: signed, TokenID: tokenID, FamilyID: familyID, ExpiresAt: expiresAt}, nil
}

// ParseAccess 解析并校验访问令牌。
func (m *Manager) ParseAccess(token string) (*Claims, error) {
	return m.parse(token, m.accessSecret, TypeAccess)
}

// ParseRefresh 解析并校验刷新令牌。
func (m *Manager) ParseRefresh(token string) (*Claims, error) {
	return m.parse(token, m.refreshSecret, TypeRefresh)
}

func (m *Manager) parse(token string, secret []byte, want TokenType) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwtlib.ParseWithClaims(token, claims,
		func(t *jwtlib.Token) (any, error) {
			// 固定签名算法，防止 alg=none 等算法混淆攻击。
			if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return secret, nil
		},
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithIssuer(m.issuer),
	)
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if !parsed.Valid {
		return nil, ErrInvalidToken
	}
	if claims.TokenType != want {
		return nil, ErrWrongType
	}
	if claims.UserID <= 0 {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
