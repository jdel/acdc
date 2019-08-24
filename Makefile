VERSION=`git rev-parse --short HEAD`

.PHONY: default all x-acdc docker clean

default: acdc

all: clean acdc x-acdc docker

acdc:
	go build -ldflags "-X github.com/jdel/acdc/cfg.Version=${VERSION}"

x-acdc:
	gox -parallel=1 -osarch="linux/386 linux/amd64 linux/arm linux/arm64 darwin/amd64 darwin/386 windows/amd64 windows/386" -output="out/{{.Dir}}-{{.OS}}-{{.Arch}}" -ldflags "-X github.com/jdel/acdc/cfg.Version=${VERSION}"

docker:
	docker build --no-cache -t jdel/acdc:local --build-arg ACDC_VERSION=${VERSION} .

clean:
	@rm -rf acdc debug out
