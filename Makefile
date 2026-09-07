.PHONY: docker-build docker-up seed migrate-up migrate-down

docker-build:
	docker compose up --build

docker-up:
	docker compose up -d --force-recreate

seed:
	docker exec -it url-shortener-api ./seed

migrate-up:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@postgres:5432/url_shortener?sslmode=disable" \
		up

migrate-down:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@postgres:5432/url_shortener?sslmode=disable" \
		down --all