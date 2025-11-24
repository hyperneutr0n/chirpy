.PHONY: migration query migrate

migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name is required. Usage: make migration name=your_migration_name"; \
		exit 1; \
	fi
	cd sql/schema && touch $(name)

query:
	@if [ -z "$(name)" ]; then \
		echo "Error: query name is required. Usage: make query name=your_query_name"; \
		exit 1; \
	fi
	cd sql/queries && touch $(name)

migrate:
	@if [ -z "$(direction)" ]; then \
		echo "Error: need migrate direction when migrating"; \
		exit 1; \
	fi
	cd sql/schema && goose postgres "postgres://postgres:changeme@localhost:5432/chirpy" $(direction)

generate:
	sqlc generate
