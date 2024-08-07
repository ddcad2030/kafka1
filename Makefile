docker:
	@docker compose down && docker compose up --build
run-producer:
	@cd producer && go build -o bin/app && ./bin/app
run-consumer:
	@cd consumer && go build -o bin/app && ./bin/app