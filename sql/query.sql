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

