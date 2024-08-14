package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker 是 Paseto token Maker
type PasetoMaker struct {	
	paseto *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker 用于生成 Paseto Maker
func NewPasetoMaker(symmertricKey string) (Maker, error) {
	if len(symmertricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size : must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmertricKey),
	}

	return maker, nil
}

// CreateToken 用于生成 token
func (maker *PasetoMaker) CreateToken(userAccount string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(userAccount, duration)

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken 用于验证 token
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}