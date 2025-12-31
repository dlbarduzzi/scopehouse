.PHONY: run
run:
	@go run ./cmd/scopehouse

.PHONY: local/run
local/run:
	@SH_CONFIG_PATH='.local/etc/scopehouse/configs' SH_SECRET_PATH='.local/etc/scopehouse/secrets' go run ./cmd/scopehouse

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: lint
lint:
	@golangci-lint run -c ./.golangci.yml ./...

.PHONY: test
test:
	@go test -count=1 ./... --cover --coverprofile=coverage.out

.PHONY: test/verbose
test/verbose:
	@go test -count=1 ./... -v --cover --coverprofile=coverage.out

.PHONY: test/coverage
test/coverage: test
	@go tool cover -html=coverage.out
