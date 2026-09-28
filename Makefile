VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS := -X github.com/Chaitanyabsprip/dotfiles/internal/core/version.Version=$(VERSION)

.PHONY: dot x both

dot:
	go install -v -buildvcs=false -ldflags "$(LDFLAGS)" ./cmd/dot

x:
	go install -v -buildvcs=false -ldflags "$(LDFLAGS)" ./cmd/x

both: dot x
