-- name: UserLogin :one
SELECT
    u.user_role,
    a.password
FROM
    "user" AS u
    JOIN auth AS a ON a.user_id = u.id
WHERE
    a.user_name = $1;

-- Create 
--
--
-- name: CreateContact :one
INSERT INTO
    contact_details (primary_number, secondary_number, email)
VALUES
    ($1, $2, $3) RETURNING id;

--
--
-- name: RegisterUser :one
INSERT INTO
    "user" (
        first_name,
        last_name,
        location,
        user_role,
        contact_id,
        emergency_contact,
        created_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

--
--
-- name: CreateUserAuth :exec
INSERT into
    auth(user_id, user_name, password)
VALUES
    ($1, $2, $3);

--
--
-- name: CreateAddress :exec
INSERT INTO
    address (street, city, state, country, zip)
VALUES
    ($1, $2, $3, $4, $5);


--
--
-- name: InsertSession :exec
INSERT INTO 
    "session" (
        user_id, 
        user_role, 
        refresh_token, 
        is_revoked, 
        created_at, 
        expired_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6);






-- name: IsUserLoggedIn :one
SELECT EXISTS (
    select 1
    from public.auth a
    join public.user u on a.user_id = u.id
    where u.user_role=$1 AND a.user_name=$2
) AS exists;    