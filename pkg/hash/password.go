// Package hash 封装密码哈希与强度校验。
package hash

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// ErrWeakPassword 密码强度不足。
var ErrWeakPassword = errors.New("password too weak")

const (
	minPasswordLen = 8
	maxPasswordLen = 72 // bcrypt 仅使用前 72 字节，超长直接拒绝避免截断歧义
)

// Hasher 密码哈希接口，便于后续替换为 argon2id。
type Hasher interface {
	Hash(plain string) (string, error)
	Verify(hashed, plain string) bool
}

// BcryptHasher 基于 bcrypt 的实现。
type BcryptHasher struct {
	cost int
}

// NewBcrypt 创建 bcrypt 哈希器，cost 低于默认值时取默认值。
func NewBcrypt(cost int) *BcryptHasher {
	if cost < bcrypt.DefaultCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash 生成密码哈希。
func (h *BcryptHasher) Hash(plain string) (string, error) {
	if len(plain) > maxPasswordLen {
		return "", ErrWeakPassword
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Verify 校验明文密码与哈希是否匹配。
func (h *BcryptHasher) Verify(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

// CheckStrength 校验密码强度：长度 8-72，且至少包含字母与数字。
func CheckStrength(plain string) error {
	if len(plain) < minPasswordLen || len(plain) > maxPasswordLen {
		return ErrWeakPassword
	}
	var hasLetter, hasDigit bool
	for _, r := range plain {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return ErrWeakPassword
	}
	return nil
}
