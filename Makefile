GOBIN := $(shell go env GOPATH)/bin/
MODULE := $(shell head -n 1 go.mod | cut -d' ' - -f2)
PACKRGO := "pkg/assets/assets-packr.go"

build:
	$(MAKE) packr-up
	CGO_ENABLED=0 go build -v -trimpath -o bin/ghstsbot cmd/bot.go
	$(MAKE) packr-down
test:
	go test -v ./...
format:
	go fmt ./...
	golines -w .
check-format:
	./ensure_formatted.sh
clean:
	rm -r bin/ || true
	$(MAKE) packr-down || true

setup-dep:
	go mod download -x
setup-devdep:
	go get -u github.com/segmentio/golines
	go get -u github.com/gobuffalo/packr/v2/packr2
	go mod tidy
packr-up:
	$(GOBIN)packr2 --ignore-imports
	sed -i 's!^import _.*!import _ "$(MODULE)/packrd"!' $(PACKRGO)
packr-down:
	$(GOBIN)packr2 clean
	rm $(PACKRGO)

.PHONY: build
.PHONY: test
.PHONY: format
.PHONY: check-format
.PHONY: clean
# Adapted from https://gitlab.com/NickCao/RAIT/-/blob/master/Makefile
# Ref: https://superuser.com/questions/429693/git-list-all-files-currently-under-source-control
# Ref: https://stackoverflow.com/questions/39792766/checking-to-find-out-if-go-code-has-been-formatted/39796269
# Ref: https://stackoverflow.com/questions/2145590/what-is-the-purpose-of-phony-in-a-makefile
# Ref: https://github.com/gobuffalo/packr/issues/175
# Ref: https://github.com/gobuffalo/packr/issues/169#issuecomment-474205649
# Ref: https://github.com/golang/go/issues/31273#issuecomment-480196029

# WTF/TODO: How to go get without touching go.{mod,sum?}
#	ref: https://github.com/golang/go/issues/31273
# 	ref: https://github.com/golang/go/issues/30515