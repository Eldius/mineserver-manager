
TEST_SERVER ?= 192.168.0.42


test:
	go clean -cache
	go test ./... -cover -covermode=set -coverpkg=./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

install:
	$(eval TMP := $(shell mktemp -d))
	@echo "temp folder: $(TMP)"
	go run ./cmd/cli install --dest $(TMP) --headless --rcon-enabled --motd "My Awsome Server" --level-name "My Precious" --seed "$(TMP)"

versions:
	go run ./cmd/cli install --list

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run

validate: test lint vulncheck
	@echo "Validate finished."

snapshot-local:
	goreleaser release --snapshot --clean

release-local:
	goreleaser release --clean --skip=publish

put:
	echo 'rm /home/eldius/.bin/mineserver' | sftp $(USER)@$(TEST_SERVER)
	echo 'put ./dist/mineserver-cli_linux_arm64_v8.0/mineserver .bin/' | sftp $(USER)@$(TEST_SERVER)

install-remote:
	@echo "###################################"
	@echo "#  Backing up old instalation...  #"
	@echo "###################################"
	@echo
	ssh $(USER)@$(TEST_SERVER) '~/.bin/mineserver backup save --instance-folder /mineservers/vanila-1.21.3-aikar --backup-folder /mineservers/backup'
	@echo
	@echo
	@echo "##################################"
	@echo "#  Cleaning old installation...  #"
	@echo "##################################"
	@echo
	ssh $(USER)@$(TEST_SERVER) 'rm -rf /mineservers/vanila-1.21.3-aikar'
	@echo
	@echo
	@echo "###############################"
	@echo "#  Making a clean install...  #"
	@echo "###############################"
	@echo
	ssh $(USER)@$(TEST_SERVER) '~/.bin/mineserver install \
		--version 1.21.3 \
		--motd "A new test server" \
		--level-name "My test world" \
		--rcon-enabled \
		--rcon-passwd "MyStrongP@ss#123" \
		--headless \
		--dest /mineservers/vanila-1.21.3-aikar \
		--seed "5516949179205280665" \
		--memory-limit 2g \
		--whitelist-user Eldius'

backup-remote:
	@echo "###############################"
	@echo "#  Backing up instalation...  #"
	@echo "###############################"
	@echo
	ssh $(USER)@$(TEST_SERVER) '~/.bin/mineserver backup save \
		--instance-folder /mineservers/test-server-backup \
		--backup-folder /mineservers/backup \
		--max-backup-files 5'

restore-remote:
	@echo "###############################"
	@echo "#  Backing up instalation...  #"
	@echo "###############################"
	@echo
	ssh $(USER)@$(TEST_SERVER) '~/.bin/mineserver backup restore \
		--debug \
		--instance-folder /mineservers/test-server-backup-restore \
		--backup-file /mineservers/backup/server-mundo-da-duda_2024-12-06_00-00-01_backup.zip'

.tmp:
	-mkdir -p .tmp

get-remote-backups: .tmp
	sftp $(USER)@$(TEST_SERVER):/mineservers/backup/test-server-backup_*_backup.zip* .tmp/
