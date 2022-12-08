 
run: 
	go run main.go
migrateup:
	migrate -path ./storage/migrations -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable' up
migratedown:
	migrate -path ./storage/migrations -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable' down
pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule foreach --recursive git pull
	# git submodule update --remote --merge
pull-proto:
	git submodule update --init --recursive