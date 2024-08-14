-- name: CreateUser :execresult
INSERT INTO users (
    user_account, 
    hash_password, 
    username, 
    email, 
    password_chaged_at, 
    created_at
) VALUES (
    ?, ?, ?, ?, DEFAULT, DEFAULT
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_account = ? LIMIT 1;