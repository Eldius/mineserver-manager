
test:
	go clean -cache
	go test ./... -cover

install:
	$(eval TMP := $(shell mktemp -d))
	@echo "temp folder: $(TMP)"
	go run ./cmd/cli install --dest $(TMP) --headless --enable-rcon --motd "My Awsome Server" --level-name "My Precious" --seed "$(TMP)"

versions:
	go run ./cmd/cli install --list

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run

snapshot-local:
	goreleaser release --snapshot --clean

release-local:
	goreleaser release --clean --skip=publish

put:
	echo 'rm ~/.bin/mineserver' | sftp $(USER)@192.168.100.183
	echo 'put ./dist/mineserver-cli_linux_arm64/mineserver .bin/' | sftp $(USER)@192.168.100.183
