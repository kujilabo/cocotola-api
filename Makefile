.PHONY: unit-test docker-up docker-down test-docker-up test-docker-down docker-clear

unit-test: dependency
	@go test -v -short ./...

docker-up:
	@docker-compose -f docker/development/docker-compose.yml up -d
	sleep 10
	# @chmod -R 777 docker/development

docker-down:
	@chmod -R 777 docker/development
	@docker-compose -f docker/development/docker-compose.yml down
	@chmod -R 777 docker/development

test-docker-up:
	@docker-compose -f docker/test/docker-compose.yml up -d
	sleep 10
	@chmod -R 777 docker/test

test-docker-down:
	@docker-compose -f docker/test/docker-compose.yml down

docker-clear:
	@rm -rf docker/development/mysql-data
	@rm -rf docker/test/mysql-data
