.PHONY: run-api
run-api:
	go run cmd/api/main.go

CONTAINER_NAME=mini-url-shortener-mysql-server
MYSQL_ROOT_PASSWORD=root
MYSQL_IMAGE=mysql:latest
MYSQL_PORT=3306
MYSQL_VOLUME=mini-url-shortener-mysql-data
MYSQL_DATABASE=mini-url-shortener

.PHONY: setup-db
setup-db:
	docker run -d \
		--name $(CONTAINER_NAME) \
		-e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) \
		-e MYSQL_DATABASE=$(MYSQL_DATABASE) \
		-p $(MYSQL_PORT):3306 \
		-v $(MYSQL_VOLUME):/var/lib/mysql \
		$(MYSQL_IMAGE)

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:$(MYSQL_PORT))/$(MYSQL_DATABASE)?parseTime=true" up

.PHONY: clean
clean:
	migrate -path ./migrations -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:$(MYSQL_PORT))/$(MYSQL_DATABASE)?parseTime=true" down

.PHONY: test
test:
	go test -v ./...