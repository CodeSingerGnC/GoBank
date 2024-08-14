package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt 生成在一个 min 和 max 范围内的随机值
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成一个长度为 n 的字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUser 随机生成长度为 6 的用户名称（用于测试）
func RandomUser() string {
	return RandomString(6)
}

// RandomMoney 随机生成 [1, 1000] 的银行余额（用于测试）
func RandomMoney() int64 {
	return RandomInt(1, 1000)
}

// RandomCurreny 随机选取一个货币简称
func RandomCurreny() string {
	curreny := []string{
		"CNY",
		"EUR",
		"USD",
	}
	n := len(curreny)
	return curreny[rand.Intn(n)]
}

// RandomEmail 随机生成一个邮箱地址
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}