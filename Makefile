.PHONY: gen-src unit-test swagger proto docker-up docker-down test-docker-up test-docker-down docker-clear

gen-src:
	@go generate ./src/...

unit-test:
	@go test -v -short ./src/... -p 1 -count=1

swagger:
	@swagger init -d src

proto:
	@protoc --go_out=./src/ --go_opt=paths=source_relative \
	--go-grpc_out=./src/ --go-grpc_opt=paths=source_relative \
	--proto_path=../cocotola-translator-api \
	proto/translator_admin.proto
	@protoc --go_out=./src/ --go_opt=paths=source_relative \
	--go-grpc_out=./src/ --go-grpc_opt=paths=source_relative \
	--proto_path=../cocotola-translator-api \
	proto/translator_user.proto


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
