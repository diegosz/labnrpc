.PHONY: install-tools generate example

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


example:
	tmux new-session -d -s ntest
	tmux split-window -t "ntest:0"   -v
	tmux split-window -t "ntest:0.0" -h -p 66
	tmux split-window -t "ntest:0.1" -h -p 50
	tmux select-pane -t "ntest:0.3"
	tmux send-keys -t "ntest:0.0" "sleep 2 && go run cmd/client/main.go" Enter
	tmux send-keys -t "ntest:0.1" "sleep 2 && go run cmd/client/main.go" Enter
	tmux send-keys -t "ntest:0.2" "sleep 3 && go run cmd/client/main.go" Enter
	tmux send-keys -t "ntest:0.3" "go run cmd/server/main.go" Enter
	tmux attach -t ntest
	tmux kill-session -t ntest

example2:
	tmux new-session -d -s ntest
	tmux split-window -t "ntest:0"   -v
	tmux split-window -t "ntest:0.1" -h -p 66
	tmux split-window -t "ntest:0.2" -h -p 50
	tmux select-pane -t "ntest:0.0"
	tmux send-keys -t "ntest:0.0" "go run cmd/server/main.go" Enter
	tmux send-keys -t "ntest:0.1" "sleep 2 && go run cmd/client/main.go" Enter
	tmux send-keys -t "ntest:0.2" "sleep 2 && go run cmd/client/main.go" Enter
	tmux send-keys -t "ntest:0.3" "sleep 3 && go run cmd/client/main.go" Enter
	tmux attach -t ntest
	tmux kill-session -t ntest
