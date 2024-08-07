## createDB: create postgresDB
.PHONY: createDB
createDB:
	docker run --name songAPI -p 5430:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine
	docker exec -it songAPI createdb --username=root --owner=root musicDB
## run: run server
.PHONY: run
run:
	go run ./cmd/api/
