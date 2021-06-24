@REM https://habr.com/ru/post/249449/

@SET GOOS=windows
@SET GOARCH=amd64
@SET CGO_ENABLED=1
@SET CXX=x86_64-w64-mingw32-g++.exe
@SET CXX_FOR_TARGET=x86_64-w64-mingw32-g++.exe
@SET CC=x86_64-w64-mingw32-gcc.exe
@SET CC_FOR_TARGET=x86_64-w64-mingw32-gcc.exe
@SET GCCGO=gccgo
@SET GOTOOLDIR=C:\Users\dbl0null\sdk\go1.16.5\pkg\tool\windows_amd64
@SET GOROOT=C:\Users\dbl0null\sdk\go1.16.5
@SET GOPATH=C:\Users\dbl0null\go
@SET GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -gno-record-gcc-switches

set GO111MODULE=on
set GOARCH=amd64
set GOBIN=
set GOCACHE=C:\Users\dbl0null\AppData\Local\go-build
set GOENV=C:\Users\dbl0null\AppData\Roaming\go\env
set GOEXE=.exe
set GOFLAGS=
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOINSECURE=
set GOMODCACHE=C:\Users\dbl0null\go\pkg\mod
set GONOPROXY=
set GONOSUMDB=
set GOOS=windows
set GOPATH=C:\Users\dbl0null\go
set GOPRIVATE=
set GOPROXY=https://proxy.golang.org,direct
set GOROOT=C:\Users\dbl0null\sdk\go1.16.5
set GOSUMDB=sum.golang.org
set GOTMPDIR=
set GOTOOLDIR=C:\Users\dbl0null\sdk\go1.16.5\pkg\tool\windows_amd64
set GOVCS=
set GOVERSION=go1.16.5
set GCCGO=gccgo
set AR=ar
set CC=x86_64-w64-mingw32-gcc.exe
set CXX=x86_64-w64-mingw32-g++.exe
set CGO_ENABLED=1
set GOMOD=C:\Users\dbl0null\GolandProjects\mediadevices\examples\streamer\go.mod
set CGO_CFLAGS=-g -O2
set CGO_CPPFLAGS=
set CGO_CXXFLAGS=-g -O2
set CGO_FFLAGS=-g -O2
set CGO_LDFLAGS=-g -O2
set PKG_CONFIG=pkg-config
set GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -gno-record-gcc-switches


go env

go build -v -x -o bin/streamer_amd64.exe
go build -v -x -ldflags "-s -w" -o bin/streamer_amd64_stripped.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-s -w -extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static_stripped.exe