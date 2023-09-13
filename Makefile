install-tools:
	@echo installing tools
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@go install github.com/nats-rpc/nrpc/protoc-gen-nrpc@latest
	@echo done

generate:
	@echo running code generation
	@go generate ./...
	@sed -i 's/_ \"labnrpc\/provisioning\/provisioningpb\"//' provisioning/provisioningpb/api.pb.go
	@rpl -q --encoding 'UTF-8' "Server" "NRPCServer" provisioning/provisioningpb/api.nrpc.go
	@rpl -q --encoding 'UTF-8' "Client" "NRPCClient" provisioning/provisioningpb/api.nrpc.go
	@echo done
