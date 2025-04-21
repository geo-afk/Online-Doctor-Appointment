-- name: UserLogin :one
SELECT 
    u.user_role,
    u.id
FROM 
    "user" AS u
JOIN auth AS a ON 
    a.user_id = u.id
WHERE a.user_name = $1 AND a.password = $2;
