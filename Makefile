GOPATH := $(HOME)/go/bin/
MODULE := $(shell head -n 1 go.mod | cut -d' ' - -f2)
PACKRGO := "pkg/assets/assets-packr.go"

build:
	$(MAKE) packr-up
	CGO_ENABLED=0 go build -trimpath -o bin/ghstsbot cmd/bot.go
	$(MAKE) packr-down
test:
	go test ./...
format:
	go fmt ./...
check-format:
	./ensure_formatted.sh
clean:
	rm -r bin/ || true
	$(MAKE) packr-down || true

packr-up:
	$(GOPATH)packr2 --ignore-imports
	sed -i 's!^import _.*!import _ "$(MODULE)/packrd"!' $(PACKRGO)

packr-down:
	$(GOPATH)packr2 clean
	rm $(PACKRGO)

.PHONY: build
.PHONY: test
.PHONY: format
.PHONY: check-format
.PHONY: clean
# Adapted from https://gitlab.com/NickCao/RAIT/-/blob/master/Makefile
# Ref: https://superuser.com/questions/429693/git-list-all-files-currently-under-source-control
# Ref: https://stackoverflow.com/questions/39792766/checking-to-find-out-if-go-code-has-been-formatted/39796269
# https://stackoverflow.com/questions/2145590/what-is-the-purpose-of-phony-in-a-makefile