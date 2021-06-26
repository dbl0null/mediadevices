export GOPATH="/Users/dbl0null/go"
export GOROOT="/Users/dbl0null/.go"

GOOS=darwin GOARCH=amd64 \
	go build -x -v \
	--ldflags '--extldflags "-static -static-libgcc -static-libstdc++" -s -w' \
	-o ./bin/streamer.static.stripped.bin

GOOS=darwin GOARCH=amd64 \
	go build -x -v \
	--ldflags '--extldflags "-static -static-libgcc -static-libstdc++"' \
	-o ./bin/streamer.static.bin

GOOS=windows GOARCH=amd64 \
	CGO_ENABLED=1 \
	CXX="x86_64-w64-mingw32-g++" \
	CXX_FOR_TARGET="x86_64-w64-mingw32-g++" \
	CC="x86_64-w64-mingw32-gcc" \
	CC_FOR_TARGET="x86_64-w64-mingw32-gcc" \
	go build -x -v \
	--ldflags '--extldflags "-static -static-libgcc -static-libstdc++" -s -w' \
	-o ./bin/streamer.static.sw.exe

GOOS=windows GOARCH=amd64 \
	CGO_ENABLED=1 \
	CXX="x86_64-w64-mingw32-g++" \
	CXX_FOR_TARGET="x86_64-w64-mingw32-g++" \
	CC="x86_64-w64-mingw32-gcc" \
	CC_FOR_TARGET="x86_64-w64-mingw32-gcc" \
	go build -x -v \
	--ldflags '--extldflags "-static -static-libgcc -static-libstdc++"' \
	-o ./bin/streamer.static.exe