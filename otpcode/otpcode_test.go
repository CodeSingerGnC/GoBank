package otpcode

import (
	"fmt"
	"testing"
	"time"

	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
)

const (
	Issuer = "https://github.com/CodeSingerGnC"
)

func generateKey(accountName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      Issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func TestPassCode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	secret, err := generateKey(util.RandomEmail())
	require.NoError(t, err)
	passcode, err := GeneratePassCode(secret.Secret())
	require.NoError(t, err)
	require.Equal(t, int(Digits), len(passcode))
	wrongPasscode := util.RandomString(6)
	err = VerifyPassCode(wrongPasscode, secret.Secret())
	require.ErrorIs(t, err, ErrPassCodeMismatch)
	ssecret := secret.Secret()
	err = VerifyPassCode(passcode, ssecret)
	require.NoError(t, err)
	fmt.Println("passcode 延期测试，需要等待 30 s。")
	err = VerifyPassCode(passcode, secret.Secret())
	require.NoError(t, err)
	fmt.Println("passcode 过期测试，需要等待 60 s。")
	time.Sleep(60 * time.Second)
	err = VerifyPassCode(passcode, secret.Secret())
	require.ErrorIs(t, err, ErrPassCodeMismatch)
}
