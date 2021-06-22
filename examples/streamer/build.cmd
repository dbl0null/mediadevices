@REM https://habr.com/ru/post/249449/

@SET GOOS=windows
@SET GOARCH=amd64
go build -v -x -ldflags "-s -w" -o bin/streamer_amd64.exe

@SET GOOS=darwin
@SET GOARCH=amd64
go build -v -x -ldflags "-s -w" -o bin/streamer_darwin

@SET GOOS=linux
@SET GOARCH=386
go build -v -x -ldflags "-s -w" -o bin/streamer_i386

@SET GOOS=linux
@SET GOARCH=amd64
go build -v -x -ldflags "-s -w" -o bin/streamer_amd64

@SET GOOS=linux
@SET GOARCH=arm
@SET GOARM=7
go build -v -x -ldflags "-s -w" -o bin/streamer_armv7

@SET GOOS=linux
@SET GOARCH=arm64
go build -v -x -ldflags "-s -w" -o bin/streamer_aarch64

@SET GOOS=darwin
@SET GOARCH=amd64
go build -v -x -ldflags "-s -w" -o bin/streamer_darwin