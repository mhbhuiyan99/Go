CREATE TABLE IF NOT EXISTS users(
	id bigserial primary key,
	name text not null,
	email text not null,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone default current_timestamp
)