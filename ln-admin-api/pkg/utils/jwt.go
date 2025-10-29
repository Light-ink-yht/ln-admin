package utils

import (
	"errors"
	"time"

	"github.com/Light-ink-yht/ln-admin/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenPair 双Token结构
type TokenPair struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int64  `json:"expiresIn"`        // AccessToken过期时间（秒）
	RefreshExpiresIn int64  `json:"refreshExpiresIn"` // RefreshToken过期时间（秒）
}

// GenerateToken 生成JWT AccessToken
func GenerateToken(userID string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(config.Cfg.JWT.Expire) * time.Second)
	expiresIn := int64(config.Cfg.JWT.Expire)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Cfg.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiresIn, nil
}

// GenerateRefreshToken 生成刷新Token
func GenerateRefreshToken(userID string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(config.Cfg.JWT.RefreshExpire) * time.Second)
	expiresIn := int64(config.Cfg.JWT.RefreshExpire)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Cfg.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiresIn, nil
}

// GenerateTokenPair 生成双Token
func GenerateTokenPair(userID string) (*TokenPair, error) {
	accessToken, expiresIn, err := GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresIn, err := GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(config.Cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
