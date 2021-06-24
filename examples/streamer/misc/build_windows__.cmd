@REM https://habr.com/ru/post/249449/

@GOOS=windows
@GOARCH=amd64
@CGO_ENABLED=1
@CXX=x86_64-w64-mingw32-g++.exe
@CXX_FOR_TARGET=x86_64-w64-mingw32-g++.exe
@CC=x86_64-w64-mingw32-gcc.exe
@CC_FOR_TARGET=x86_64-w64-mingw32-gcc.exe
@GCCGO=gccgo
@GOTOOLDIR=C:\Users\dbl0null\sdk\go1.16.5\pkg\tool\windows_amd64
@GOROOT=C:\Users\dbl0null\sdk\go1.16.5
@GOPATH=C:\Users\dbl0null\go
@GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -gno-record-gcc-switches

GO111MODULE=on
GOARCH=amd64
GOBIN=
GOCACHE=C:\Users\dbl0null\AppData\Local\go-build
GOENV=C:\Users\dbl0null\AppData\Roaming\go\env
GOEXE=.exe
GOFLAGS=
GOHOSTARCH=amd64
GOHOSTOS=windows
GOINSECURE=
GOMODCACHE=C:\Users\dbl0null\go\pkg\mod
GONOPROXY=
GONOSUMDB=
GOOS=windows
GOPATH=C:\Users\dbl0null\go
GOPRIVATE=
GOPROXY=https://proxy.golang.org,direct
GOROOT=C:\Users\dbl0null\sdk\go1.16.5
GOSUMDB=sum.golang.org
GOTMPDIR=
GOTOOLDIR=C:\Users\dbl0null\sdk\go1.16.5\pkg\tool\windows_amd64
GOVCS=
GOVERSION=go1.16.5
GCCGO=gccgo
AR=ar
CC=x86_64-w64-mingw32-gcc.exe
CXX=x86_64-w64-mingw32-g++.exe
CGO_ENABLED=1
GOMOD=C:\Users\dbl0null\GolandProjects\mediadevices\examples\streamer\go.mod
CGO_CFLAGS=-g -O2
CGO_CPPFLAGS=
CGO_CXXFLAGS=-g -O2
CGO_FFLAGS=-g -O2
CGO_LDFLAGS=-g -O2
PKG_CONFIG=pkg-config
GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -gno-record-gcc-switches


go env

go build -v -x -o bin/streamer_amd64.exe
go build -v -x -ldflags "-s -w" -o bin/streamer_amd64_stripped.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-s -w -extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static_stripped.exe