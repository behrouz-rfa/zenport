install-tools:
	@echo installing tools
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@go install github.com/cucumber/godog/cmd/godog@latest
	@echo done

generate:
	@echo running code generation
	@go generate ./...
	@echo done


run:
	echo "Starting docker environment"
	docker compose  up --build

microservice:
	echo "Starting local environment"
	 docker compose --profile microservices up --build

build: build-monolith build-services

rebuild: clean-monolith clean-services build

clean-monolith:
	docker image rm zenports-monolith

clean-services:
	docker image rm zenports-gates zenports-ntps zenports-notifications

build-monolith:
	docker build -t zenports-monolith --file docker/Dockerfile .

build-services:
	docker build -t zenports-gates --file docker/Dockerfile.microservices --build-arg=service=gates .
	docker build -t zenports-ntps --file docker/Dockerfile.microservices --build-arg=service=ntps .
	docker build -t zenports-notifications --file docker/Dockerfile.microservices --build-arg=service=notifications .


