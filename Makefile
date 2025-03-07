BIN=go
OUTPATH=./bin

.PHONY: build test watch-test bench watch-bench coverage tools lint lint-fix audit outdated weight latest proto-all
build: create_build_folder copy_env
	${BIN} build -v -o ${OUTPATH} ./...

create_build_folder:
	mkdir -p bin

copy_env:
	@envfile=$$(find . -name ".env" -print -quit); \
	if [ -f "$${envfile}" ] && [ -f "${OUTPATH}/.env" ]; then \
		rm -f "${OUTPATH}/.env" || exit 1; \
	fi; \
	if [ -f "$${envfile}" ]; then \
		cp -f "$${envfile}" "${OUTPATH}/" || exit 1; \
	fi


test:
	go test -race -v ./...
watch-test:
	reflex -t 50ms -s -- sh -c 'go test -race -v ./...'

bench:
	go test -benchmem -count 3 -bench ./...
watch-bench:
	reflex -t 50ms -s -- sh -c 'go test -benchmem -count 3 -bench ./...'

coverage:
	${BIN} test -v ./... -coverprofile=cover.out -covermode=atomic
	${BIN} tool cover -html=cover.out -o cover.html

tools:
	${BIN} install github.com/cespare/reflex@latest
	${BIN} install github.com/rakyll/gotest@latest
	${BIN} install github.com/psampaz/go-mod-outdated@latest
	${BIN} install github.com/jondot/goweight@latest
	${BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	${BIN} get -t -u golang.org/x/tools/cmd/cover
	${BIN} install github.com/sonatype-nexus-community/nancy@latest
	go mod tidy

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...
lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...

audit:
	${BIN} list -json -m all | nancy sleuth

outdated:
	${BIN} list -u -m -json all | go-mod-outdated -update -direct

weight:
	goweight

latest:
	${BIN} get -t -u ./...

proto-all:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/*.proto
