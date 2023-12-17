include .env

#migrations
gooseUp:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB} up
gooseReset:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB} reset
gooseVal:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB} validate

check:
	echo ${DB_HOST}
