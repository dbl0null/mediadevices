ifneq (,$(wildcard ./.env))
    include .env
    export
endif

RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: mac-static
all: clean mac win clean
mac: clean mac-basic mac-strip stat
win: clean win-static win-static-strip stat

cmd-exists-%:
	@hash $(*) > /dev/null 2>&1 || (echo "ERROR: '$(*)' must be installed and available on your PATH."; exit 1)

clean:
	@echo "${RED}CLEANUP: ${NC}"
	@rm -rf ./bin
	@echo "${RED}CURRENT ENV:${NC}"
	@export -p

stat: 
	ls -hl ./bin

base:
	@echo "${RED}BUILDING 'base'${NC}"
	go build -x -v .

mac-basic:
	@echo "${RED}BUILDING 'mac-basic'${NC}"
	GOOS=darwin GOARCH=amd64 go build -x -v -tags ${SS_TAGS} -o ./bin/streamer.bin

mac-strip:
	@echo "${RED}BUILDING 'mac-strip'${NC}"
	GOOS=darwin GOARCH=amd64 go build ${SS_TOOL_ARGS} -tags ${SS_TAGS} \
	--ldflags '-s -w' -o ./bin/streamer.strip.bin

mac-static:
	@echo "${RED}BUILDING 'mac-static'${NC}"
	GOOS=darwin GOARCH=amd64 go build ${SS_TOOL_ARGS} -tags ${SS_TAGS} --ldflags ${SS_LDFLAGS_STATIC}  -o ./bin/streamer.static.bin

mac-static-strip:
	@echo "${RED}BUILDING 'mac-static-strip'${NC}"
	GOOS=darwin GOARCH=amd64 \
	go build ${SS_TOOL_ARGS} -tags ${SS_TAGS} --ldflags ${SS_LDFLAGS_STRIP} -o ./bin/streamer.static.strip.bin

win-static:
	@echo "${RED}BUILDING 'win-static'${NC}"
	GOOS=windows GOARCH=amd64 go build -o ./bin/streamer.static.exe

win-static-strip:
	@echo "${RED}BUILDING 'win-static-strip'${NC}"
	GOOS=windows GOARCH=amd64 go build --ldflags="-s -w" -o ./bin/streamer.static.strip.exe
