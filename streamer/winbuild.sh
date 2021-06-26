NC="\033[0;0m" # No Color
BOLD="\033[0;1m"
INVERTED="\033[0;7m"
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
LBLUE="\033[0;39m"
MAGENTA="\033[0;35m"
CYAN="\033[0;36m"
BGRED="\033[0;41m"
BGGREEN="\033[0;42m"
BGYELLOW="\033[0;42m"
BGLBLUE="\033[0;104m"


#Build WINDOWS
echo "${BGLBLUE}\t\t\tBuild WINDOWS\t\t\t${NC}"

#Build STATIC
source .env && GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
        CXX="x86_64-w64-mingw32-g++" CXX_FOR_TARGET="x86_64-w64-mingw32-g++" \
        CC="x86_64-w64-mingw32-gcc" CC_FOR_TARGET="x86_64-w64-mingw32-gcc" \
        go build -x -tags "osusergo netgo static_build static" \
        --ldflags '--extldflags "-static -static-libgcc -static-libstdc++"' \
        -o ./bin/streamer.static.exe

if [ $? -eq 0 ]; then echo "${BGGREEN}\t\t\tDone STATIC \t\t\t${NC}";
else echo "${BGRED}\t\t\tFail STATIC \t\t\t${NC}" >&2; fi

#exit

#Build STRIPPED
source .env && GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
        CXX="x86_64-w64-mingw32-g++" CXX_FOR_TARGET="x86_64-w64-mingw32-g++" \
        CC="x86_64-w64-mingw32-gcc" CC_FOR_TARGET="x86_64-w64-mingw32-gcc" \
        go build -x -tags "osusergo netgo static_build static" \
        --ldflags '-s -w --extldflags "-static -static-libgcc -static-libstdc++"' \
        -o ./bin/streamer.static.stripped.exe

if [ $? -eq 0 ]; then echo "${BGGREEN}\t\t\tDone STRIPPED\t\t\t${NC}";
else echo "${BGRED}\t\t\tFail STRIPPED\t\t\t${NC}" >&2; fi

exit

#Build MACOS
echo "${BGLBLUE}\t\t\tBuild MACOS  \t\t\t${NC}"

#Build STATIC
source .env && GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
        go build -x -v -tags "osusergo netgo static_build static" \
        --ldflags '-s -w --extldflags "-static -static-libgcc -static-libstdc++"' \
        -o ./bin/streamer.static.bin

if [ $? -eq 0 ]; then echo "${BGGREEN}\t\t\tDone STATIC \t\t\t${NC}";
else echo "${BGRED}\t\t\tFail STATIC \t\t\t${NC}" >&2; fi


#Build STRIPPED
echo "${BGLBLUE}\t\t\tBuild MACOS  \t\t\t${NC}"
source .env && GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
        go build -x -v -tags "osusergo netgo static_build static" \
        --ldflags '--extldflags "-static -static-libgcc -static-libstdc++"' \
        -o ./bin/streamer.static.stripped.bin
if [ $? -eq 0 ]; then echo "${BGGREEN}\t\t\tDone STATIC \t\t\t${NC}";
else echo "${BGRED}\t\t\tFail STATIC \t\t\t${NC}" >&2; fi
