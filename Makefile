fmt:
	go fmt ./...
test: fmt
	# go test -mod=vendor -race -cover -coverprofile=cover.out `go list ./... | grep -v ./mocks`
	export NODE_ENV=test && (go test -cover ./v0/tests || (echo "test failing" && exit 1))
clean:
	sudo rm -rf volume_docker
	docker compose rm -f
proto:
	protoc  --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  internal/offers_platform/core/salesforce-grpc-client/proto/salesforceadapter.proto