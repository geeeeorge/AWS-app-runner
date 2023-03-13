ENV_LOCAL=\
	APP_PORT=8080 \
	APP_HOST=127.0.0.1 \
	DB_HOST=localhost \
	DB_PORT=3306 \
	DB_NAME=book_review \
	DB_USER=mysql \
	DB_PASSWORD=passw0rd

.PHONY: app-start
app-start:
	$(ENV_LOCAL) \
	go run cmd/app/main.go
