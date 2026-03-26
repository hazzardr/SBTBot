-- name: ListGenres :many
select * from genres;

-- name: AddGenre :one
insert into genres(
    name
) values (
    ?
)
RETURNING *;

-- name: DeleteGenre :exec
delete from genres where name = ?;

-- name: ListThemes :many
select * from themes;

-- name: AddTheme :one
insert into themes(
    name, description
) VALUES (
    ?, ?
 )
returning *;

-- name: DeleteTheme :exec
delete from themes where name = ?;