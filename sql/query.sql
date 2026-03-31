-- name: ListGenres :many
select * from genres;

-- name: GetGenreByName :one
select * from genres
where name == ?;

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
-- name: AddBookIdea :one
insert into book_ideas (title, author, submitter, theme_id, genre_id)
values (
   sqlc.arg(title),
   sqlc.arg(author),
   sqlc.arg(submitter),
   (select t.id from themes t where t.name = sqlc.arg(theme_name)),
   (select g.id from genres g where g.name = sqlc.arg(genre_name))
) RETURNING *;