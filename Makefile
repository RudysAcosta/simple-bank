postgres:
	docker run --name postgres_db -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin123 -d postgres

createdb:
	docker exec -it postgres_db createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres_db dropdb --username=admin simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://admin:admin123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://admin:admin123@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown