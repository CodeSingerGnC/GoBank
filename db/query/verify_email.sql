-- name: CreateVerifyEmail :execresult
INSERT INTO verify_emails (
    user_account,
    email, 
    secret_code
) VALUES (?, ?, ?);

