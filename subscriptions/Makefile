DB_URL=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
# DB_URL=postgresql://adwera:mdcclxxvi@localhost:5432/test_migrations?sslmode=disable

migrate_create:
	migrate create -ext sql -dir db/migrations -seq ${MIGRATE_NAME}
migrate_up:
	migrate -path db/migrations -database "${DB_URL}" -verbose up
migrate_down:
	migrate -path db/migrations -database "${DB_URL}" -verbose down
sqlc:
	sqlc generate
tests:
	go test -v -cover ./...
mocks:
	mockgen -package mocksdb --destination db/mocks/store.go github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc TxSubscriptionsStore
start:
	go run main.go

.PHONY: create_db migrate_create migrate_up migrate_down sqlc tests mocks start