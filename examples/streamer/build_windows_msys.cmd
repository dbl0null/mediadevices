@REM https://habr.com/ru/post/249449/

@SET GOOS=windows
@SET GOARCH=amd64
@SET CGO_ENABLED=1
@SET CXX=x86_64-w64-mingw32-g++.exe
@SET CC=x86_64-w64-mingw32-gcc.exe
@SET CXX_FOR_TARGET=x86_64-w64-mingw32-g++.exe
@SET CC_FOR_TARGET=x86_64-w64-mingw32-gcc.exe
@SET GCCGO=gccgo

@SET GOTOOLDIR=C:\msys64\mingw64\lib\go\pkg\tool\windows_amd64
@SET GOROOT=C:\msys64\mingw64\lib\go

@SET GOPATH=C:\Users\dbl0null\go

@REM @SET GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -gno-record-gcc-switches

go env

go build -v -x -o bin/streamer_amd64.exe
go build -v -x -ldflags "-s -w" -o bin/streamer_amd64_stripped.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static.exe
go build -v -x -tags "osusergo netgo static_build" -ldflags "-s -w -extldflags \"-fno-PIC -static -static-libgcc -static-libstdc++\"" -o bin/streamer_amd64_static_stripped.exe