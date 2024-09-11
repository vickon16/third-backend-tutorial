-- name: CreateProject :exec
INSERT INTO projects (id, userId, name, description, repoURL, siteURL, status, dependencies, createdAt, updatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, DEFAULT, DEFAULT);

-- name: GetProjectById :one
SELECT * FROM projects  WHERE id = ?;

-- name: GetProjectsByUserId :many
SELECT * FROM projects  WHERE userId = ?;

-- name: UpdateProject :exec
UPDATE projects SET name = ?, description = ?, repoURL = ?, siteURL = ?, status = ?, dependencies = ? WHERE id = ? AND userId = ?;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ? AND userId = ?;

