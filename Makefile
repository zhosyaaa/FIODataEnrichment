.PHONY: postgres createdb dropdb migrateup migratedown sqlc

postgres:
	docker run --name TestCase -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1079 -d postgres:15

createdb:
	docker exec -it TestCase createdb --username=root --owner=root TestCase

dropdb:
	docker exec -it TestCase dropdb TestCase

migrateup:
	migrate -path internal/db/migration -database "postgresql://root:1079@localhost:5432/TestCase?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://root:1079@localhost:5432/TestCase?sslmode=disable" -verbose down

setupkafka:
	kafka-topics.sh --bootstrap-server=localhost:9092 --create --topic coordinates --replication-factor 1 --partitions 1