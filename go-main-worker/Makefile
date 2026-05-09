migrate-up:
	migrate -path ./database/migrations -database "postgresql://dev_user:dev_password@localhost:5432/dev_database?sslmode=disable" up

make-migrate:
	migrate create -ext sql -dir database/migrations -seq $(word 2,$(MAKECMDGOALS))
%:
	@:

lint:
	golangci-lint run

format:
	gofmt -w ./..
