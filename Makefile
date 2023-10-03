
test:
	go test ./... -cover

install:
	$(eval TMP := $(shell mktemp -d))
	@echo "temp folder: $(TMP)"
	go run ./cmd/cli install --dest $(TMP) --headless
