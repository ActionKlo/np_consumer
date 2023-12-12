include .env

#migrations
gooseUp:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${PORT_DB}/${POSTGRES_DB} up
gooseReset:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${PORT_DB}/${POSTGRES_DB} reset
gooseVal:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${PORT_DB}/${POSTGRES_DB} validate

check:
	echo ${HOST}