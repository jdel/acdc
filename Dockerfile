FROM golang:1.13-alpine3.10 as builder

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}
ENV DOCKER_API_VERSION=
ARG DOCKER_VERSION=${DOCKER_VERSION:-19.03.5}
ARG DOCKER_COMPOSE_VERSION=${DOCKER_COMPOSE_VERSION:-1.25.0}
ARG ACDC_VERSION=${ACDC_VERSION:-master}
ARG ACDC_COMMIT=

COPY . /src

WORKDIR /src

RUN apk add --update curl gcc build-base \
 && go get -v ./... \
 && go test -v ./... \
 && go build -ldflags "-s -w -X github.com/jdel/acdc/cfg.Version=${ACDC_VERSION}" \
 && curl -sL https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz | tar xfvz - --strip 1 -C /usr/local/bin/ docker/docker \
 && curl -sL https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64 -o /usr/local/bin/docker-compose \
 && chmod +x /src/acdc /usr/local/bin/docker /usr/local/bin/docker-compose

FROM jdel/alpine-glibc:3.10
LABEL maintainer=julien@del-piccolo.com

COPY --from=builder /src/acdc /usr/local/bin/acdc
COPY --from=builder /usr/local/bin/docker /usr/local/bin/docker
COPY --from=builder /usr/local/bin/docker-compose /usr/local/bin/docker-compose

RUN mkdir /home/user/acdc \
 && chown user:user /home/user/acdc

EXPOSE 8080

VOLUME ["/tmp/", "/home/user/acdc"]
 
CMD ["/usr/local/bin/acdc"]
