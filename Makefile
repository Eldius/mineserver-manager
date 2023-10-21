
test:
	go test ./... -cover

install:
	$(eval TMP := $(shell mktemp -d))
	@echo "temp folder: $(TMP)"
	go run ./cmd/cli install --dest $(TMP) --headless

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run

snapshot-local:
	goreleaser release --snapshot --clean

release-local:
	goreleaser release --clean --skip=publish
