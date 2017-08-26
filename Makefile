VERSION=`git rev-parse --short HEAD`

.PHONY: default all x-acdc clean

default: acdc

all: clean acdc x-acdc

acdc:
	go build -ldflags "-X github.com/jdel/acdc/cfg.Version=${VERSION}"

x-acdc:
	gox -parallel=1 -osarch="linux/386 linux/amd64 linux/arm darwin/amd64 darwin/386 windows/amd64 windows/386" -output="out/{{.Dir}}-{{.OS}}-{{.Arch}}" -ldflags "-X github.com/jdel/acdc/cfg.Version=${VERSION}"

clean:
	@rm -rf acdc debug out
