BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
BUILD_VERSION ?= $$(cat BUILD_VERSION)
STATIC_BINARY =-a -tags netgo -installsuffix netgo
LD_FLAGS_COMMON=-w -s \
	-X github.com/bcsimms/uipo/version.binaryBuildDate=$(BUILD_DATE)
LD_FLAGS =$(LD_FLAGS_COMMON) \
	-X github.com/bcsimms/uipo/version.binaryVersion=$(BUILD_VERSION)
LD_FLAGS_LINUX =$(LD_FLAGS)

GOSRC = $(shell find . -name "*.go" ! -name "*test.go" ! -name "*fake*" ! -path "./integration/*")


all: uipo win32.exe linux_i686 linux_x86-64

build: uipo

clean: ## Remove all files from the `out` directory
	rm -f $(wildcard out/uipo*)

# Native binary for the current platform
uipo: $(GOSRC)
	mkdir -p out
	go build -ldflags "$(LD_FLAGS)" -o out/uipo

win32.exe: $(GOSRC)
	GOARCH=386 GOOS=windows go build -tags="forceposix" -o out/uipo-win32.exe -ldflags "$(LD_FLAGS)" .

linux_i686: $(GOSRC)
	CGO_ENABLED=0 GOARCH=386 GOOS=linux go build \
							$(STATIC_BINARY) \
							-ldflags "$(LD_FLAGS_LINUX)" -o out/uipo-linux_i686 .

linux_x86-64: $(GOSRC)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
							$(STATIC_BINARY) \
							-ldflags "$(LD_FLAGS_LINUX)" -o out/uipo-linux_x86-64 .
