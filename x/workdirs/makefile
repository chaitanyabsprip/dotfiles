.DEFAULT_GOAL:=install
SOURCES := $(shell find . -type f -name '*.go')

clean:
	@rm -rd ./bin

install: ${SOURCES}
	@go install -v ./cmd/work
	@go install -v ./cmd/workdirs
	@go install -v ./cmd/worktrees

work: ${SOURCES}
	@go run ./cmd/work

uninstall: clean
	@rm "$(which workdirs)"
	@rm "$(which worktrees)"
