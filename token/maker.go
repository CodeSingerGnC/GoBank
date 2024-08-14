package token

import "time"

type Maker interface {
	// CreateToken 是生成 token 的函数接口
	CreateToken(username string, duration time.Duration) (string, *Payload, error) 
	// VerifyToken 是验证 token 的函数接口
	VerifyToken(token string) (*Payload, error)
}