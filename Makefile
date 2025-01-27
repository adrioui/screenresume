.PHONY: apply-schema

DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
DEV_URL=docker://postgres
SCHEMA_FILE=sqlc/schema.sql

apply-schema:
	atlas schema apply \
		--url "$(DB_URL)" \
		--dev-url "$(DEV_URL)" \
		--to "file://$(SCHEMA_FILE)"
