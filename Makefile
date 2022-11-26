
protoc:
	protoc --go_out=. --go-grpc_out=. ./sdk/*.proto

build-dev-docker:
	docker build -t babel:dev .
