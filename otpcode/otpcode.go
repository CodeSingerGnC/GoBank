package otpcode

import (
	"encoding/base32"
	"errors"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	// PeriodSecond 定义 passcode 过期时间
	PeriodSecond = 30
	// Skew 允许过期时间便宜来适应时钟偏移或者网络延迟
	Skew = 1
	// Digits 定义 passcode 的数字长度
	Digits = otp.DigitsSix
	// Algorithm 定义 TOTP 算法
	Algorithm = otp.AlgorithmSHA256
)

var (
	// ErrPassCodeMismatch 定义 passcode 不匹配错误
	ErrPassCodeMismatch = errors.New("passcode mismatch")
)

// GeneratePassCode 用于生成 totp passcode，其中密钥格式应该是 utf8string。
func GeneratePassCode(secret string) (string, error) {
	secret = base32.StdEncoding.EncodeToString([]byte(secret))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    PeriodSecond,
		Skew:      Skew,
		Digits:    Digits,
		Algorithm: Algorithm,
	})
	if err != nil {
		return "", err
	}
	return passcode, nil
}

// VerifyPassCode 用于验证 GeneratePassCode 生成的 totp passcode。
func VerifyPassCode(passcode, secret string) error {
	secret = base32.StdEncoding.EncodeToString([]byte(secret))
	ok, err := totp.ValidateCustom(passcode, secret, time.Now(), totp.ValidateOpts{
		Period:    PeriodSecond,
		Skew:      Skew,
		Digits:    Digits,
		Algorithm: Algorithm,
	})
	if err != nil {
		return err
	}
	if !ok {
		return ErrPassCodeMismatch
	}
	return nil
}
