dev.up:
	docker compose up --build

dev.down:
	docker compose down

dev.generate:
	docker compose exec server go generate ./...