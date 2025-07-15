Create migration
- migrate create -ext sql -dir migrations -seq create_users_table

Apply Migrations
- migrate -database 'postgres://myuser:mypassword@localhost:5432/etsydb?sslmode=disable' -path migrations up

Revert Migrations
- migrate -database 'postgres://myuser:mypassword@localhost:5432/etsydb?sslmode=disable' -path migrations down