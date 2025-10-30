-- name: GetContacts :many
SELECT *
FROM contacts c
ORDER BY c.name ASC;
-- name: GetContactByID :one
SELECT *
FROM contacts
WHERE id = $1;
-- name: CreateContact :one
INSERT INTO contacts (name, phone)
VALUES ($1, $2)
RETURNING *;
-- name: UpdateContact :one
UPDATE contacts
SET name = $2,
    phone = $3
WHERE id = $1
RETURNING *;
-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1;