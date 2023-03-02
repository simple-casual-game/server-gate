.PHONY: proto 

proto:
	@for file in `ls ./protobuf/`; do \
		protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false  --proto_path=./protobuf ./protobuf/$${file}/*.proto; \
	done ;
	cp -r github.com/simple-casual-game/server-gate/protobuf/* ./protobuf/
	rm -r github.com


deprecated:
	protoc -I=. --go_out=. protobuf/common.proto
	protoc -I=. --go_out=. protobuf/common_game.proto
	protoc -I=. --go_out=. protobuf/diamonds/common.proto
	protoc -I=. --go_out=. protobuf/diamonds/game.proto


	protoc -I=. --go_out=. protobuf/gate.proto
	protoc -I=. --go_out=. protobuf/gate_game.proto
	protoc -I=. --go_out=plugins=grpc:. common.proto

	cp -r github.com/simple-casual-game/server-gate/protobuf/* ./protobuf/
	rm -r github.com

	@for file in `ls ./protobuf/`; do \
		protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false  --proto_path=./protobuf ./protobuf/$${file}/*.proto; \
	done ;
	cp -r github.com/paper-trade-chatbot/be-proto/* ./
	rm -r github.com

	for file in `ls ./protobuf/`; do   protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false