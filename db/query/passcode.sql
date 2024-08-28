-- name: CreateOtpsecret :execresult
INSERT INTO otpsecrets (
    email, secret
)  VALUES (
    ?, ?
);

-- name: GetOtpsecret :one
SELECT * FROM otpsecrets
WHERE email = ? LIMIT 1;

-- name: AddOtpsecretTryTime :exec
UPDATE otpsecrets
SET tried_times = tried_times + 1
WHERE email = ?;

-- name: DeleteOtpsecret :exec
DELETE FROM otpsecrets
WHERE email = ?;