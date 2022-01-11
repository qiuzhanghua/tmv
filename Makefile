all: pre build build-linux build-darwin build-windows post

pre:
	autotag write

build:
	go build -o tmv

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tmv_`autotag current`_linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o tmv_`autotag current`_linux_386

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o tmv_`autotag current`_darwin_amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o tmv_`autotag current`_darwin_arm64

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o tmv_`autotag current`_windows_amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o tmv_`autotag current`_windows_386.exe

post:
	git restore autotag.go

clean:
	rm tmv*
