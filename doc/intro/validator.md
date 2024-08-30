# validator

目前项目中用于验证请求参数正确性的函数一共 6 个，其函数签名如下：

```go
ValidateString(value string, minLength int, maxLength int) error
ValidateUserAccount(value string) error
ValidateUsername(value string) error
ValidatePassword(value string) error
ValidateEmail(value string) error
ValidatePasscode(value string) error
```

ValidateString 是一个较为通用的 API，用于验证字符串的长度是否在 [minLength, maxLength] 范围内。

ValidateUserAccount 用于验证创建账户请求中的账户参数的正确性，该字符串的长度在 [3, 100] 之间，并且只包含小写字母 a-z，阿拉伯数字 0-9，以及下划线 _ 。

ValidateUsername 用于验证创建账户请求中的用户名的正确性，该字符串的长度在 [3, 100] 之间，并且只包含小写字母 a-z，大写字母 A-Z，以及下划线 _ 。

ValidatePassword 用于验证密码正确性，该字符串长度在 [6, 18] 之间，并且只包含小写字母 a-z，阿拉伯数字 0-9，大写字母 A-Z，以及特殊字符 &$。

ValidateEmail 用于验证邮箱地址的合法性，该字符串长度在 [3, 100] 之间，并且使用 mail.ParseAddress() 解析邮箱地址格式。

ValidatePasscode 用于验证 totp passcode 的格式是否正确，totp passcode 长度为 otpcode.Digits 并且只能是数字。

其中 ValidateUserAccount, ValidateUsername, ValidatePassword, ValidatePasscode 在验证内部包含什么字符时，用到正则表达式，类似：`isValidUserAccount = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString`。

函数具体实现见：[/val/validator.go](../val/validator.go)
