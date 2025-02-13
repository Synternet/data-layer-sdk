.PHONY: proto
proto:
	cd proto && buf generate 

.PHONY: mocks
mocks:
	mockery

.PHONY: vulncheck
vulncheck: $(BUILD_PATH)/
	GOBIN=$(BUILD_PATH) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILD_PATH)/govulncheck ./...

.PHONY: cover
cover:
	go test -coverprofile=.coverage.tmp ./...
	cat .coverage.tmp | grep -Ev '/mock_|/.*options.go' > .coverage
	go tool cover -func=.coverage

.PHONY: install-tools
install-tools:
	go install github.com/vektra/mockery/v2@v2.42.3
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	echo Please install buf
