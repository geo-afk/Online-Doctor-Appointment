-- name: UserLogin :one
SELECT
    u.id,
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
        id,
        user_id, 
        user_role, 
        refresh_token, 
        is_revoked, 
        created_at, 
        expires_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7);



-- name: GetSession :one
SELECT * from "session"
where id=$1;



-- name: RevokeSession :exec
UPDATE "session" SET is_revoked=1
WHERE id=$1;


-- name: DeleteSession :exec
DELETE FROM "session"
WHERE id=$1;



-- name: IsUserLoggedIn :one
SELECT EXISTS (
    select 1
    from public.auth a
    join public.user u on a.user_id = u.id
    where u.user_role=$1 AND a.user_id=$2
) AS exists;    



-- name: BookAppointment :exec
INSERT INTO appointment ( 
    user_id,
    reason, 
    booked_at, 
    appointment_status, 
    created_at
    )
VALUES ($1,$2,$3,$4,$5);



-- name: UserByEmail :one
SELECT EXISTS (
    SELECT
        c.email
    FROM
        "user" AS u
        JOIN auth AS a ON a.user_id = u.id
        LEFT JOIN contact_details c ON u.contact_id=c.id
    WHERE c.email = $1
) AS exists;  




-- name: ForgetPassword :exec
UPDATE auth
SET password = $1
FROM "user" AS u
JOIN contact_details c ON u.contact_id = c.id
LEFT JOIN requests r ON c.email = r.user_email
WHERE auth.user_id = u.id
  AND r.token = $2;



-- name: ChangePassword :exec
UPDATE auth
SET password = $1
FROM "user" AS u
WHERE auth.user_id = u.id
    and u.id=$2;



-- name: CreateRequest :exec
INSERT INTO requests ( 
   request_type,
   token,
   expires_at,
   user_email
   )
VALUES ($1,$2,$3, $4);



-- name: GetRequest :one
SELECT 
    token, expires_at
FROM 
    requests
WHERE 
    token=$1;




-- name: RemoveRequest :exec
DELETE FROM 
    requests
WHERE 
    token=$1;



