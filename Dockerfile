FROM jdel/alpine:3.6

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}
ENV DOCKER_API_VERSION=
ARG DOCKER_VERSION=${DOCKER_VERSION:-latest}
ARG DOCKER_COMPOSE_VERSION=${DOCKER_COMPOSE_VERSION:-latest}
ARG ACDC_VERSION=${ACDC_VERSION:-master}
ARG ACDC_COMMIT=

LABEL maintainer=julien@del-piccolo.com

USER root

RUN apk add --update curl \
 # Install glibc on Alpine (required by docker-compose) from
 # https://github.com/sgerrand/alpine-pkg-glibc
 # See also https://github.com/gliderlabs/docker-alpine/issues/11 
 && GLIBC_VERSION='2.23-r3' \
 && curl -Lo /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
 && curl -Lo glibc.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$GLIBC_VERSION/glibc-$GLIBC_VERSION.apk \
 && curl -Lo glibc-bin.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$GLIBC_VERSION/glibc-bin-$GLIBC_VERSION.apk \
 && apk update \
 && apk add glibc.apk glibc-bin.apk \
 && rm glibc.apk glibc-bin.apk \
 && apk add --virtual build-dependencies go gcc build-base glide git openssh-client \
 && adduser acdc -D \
 && chown -R acdc:acdc /tmp /home/user \
 && curl -sL https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}.tgz -o docker.tgz \
 && tar xfvz docker.tgz --strip 1 -C /usr/local/bin/ docker/docker \
 && rm -f docker.tgz \
 && curl -sL https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64 -o /usr/local/bin/docker-compose \
 && curl -sL https://github.com/jdel/acdc/archive/${ACDC_VERSION}.zip -o acdc.zip \
 && mkdir -p ${GOPATH}/src/github.com/jdel/ \
 && unzip acdc.zip -d ${GOPATH}/src/github.com/jdel/ \
 && rm -f acdc.zip \
 && mv ${GOPATH}/src/github.com/jdel/acdc-* ${GOPATH}/src/github.com/jdel/acdc \
 && cd ${GOPATH}/src/github.com/jdel/acdc/ \
 && glide install -v \
 && go build -o /usr/local/bin/acdc -ldflags "-X github.com/jdel/acdc/cfg.Version=${ACDC_VERSION}-${ACDC_COMMIT}" \
 && chmod 755 /usr/local/bin/docker /usr/local/bin/docker-compose /usr/local/bin/acdc \
 && apk del build-dependencies \
 && rm -rf /var/cache/apk/* \
 && rm -rf /root/.glide/ \
 && rm -rf ${GOPATH}
 
USER user

WORKDIR /home/user/

EXPOSE 8080

VOLUME ["/tmp/", "/home/user/acdc"]
 
CMD ["/usr/local/bin/acdc"]
