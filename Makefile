.PHONY: postgres createdb dropdb migrateup migratedown sqlc

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1079 -d postgres:15
redis:
	docker run -d --name redis -p 6379:6379 redis/redis-stack-server:latest
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres TestCase

dropdb:
	docker exec -it postgres dropdb TestCase

migrateup:
	migrate -path internal/db/migration -database "postgresql://postgres:1079@localhost:5432/TestCase?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://postgres:1079@localhost:5432/TestCase?sslmode=disable" -verbose down
