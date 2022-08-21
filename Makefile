postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

mysql:
	docker run --name mysql8 -p 3356:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8-debian
	
mysql_createdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "create database simple_bank"

mysql_dropdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "drop database simple_bank"

mysql_migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3356)/simple_bank" -verbose up

mysql_migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3356)/simple_bank" -verbose down

test:
	go test -v -cover ./...
sqlc:
	.\sqlc.exe generate
.PHONY: postgres createdb dropdb migrateup migratedown mysql_createdb mysql_dropdb mysql mysql_migrateup mysql_migratedown sqlc test