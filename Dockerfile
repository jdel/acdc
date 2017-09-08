FROM jdel/alpine-glibc:3.6

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}
ENV DOCKER_API_VERSION=
ARG DOCKER_VERSION=${DOCKER_VERSION:-17.06.2-ce}
ARG DOCKER_COMPOSE_VERSION=${DOCKER_COMPOSE_VERSION:-1.16.1}
ARG ACDC_VERSION=${ACDC_VERSION:-master}
ARG ACDC_COMMIT=

LABEL maintainer=julien@del-piccolo.com

USER root

RUN apk add --update curl \
 && apk add --virtual build-dependencies go gcc build-base glide git openssh-client \
 && curl -sL https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz -o docker.tgz \
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
 && mkdir /home/user/acdc \
 && chown user:user /home/user/acdc \
 && rm -rf /var/cache/apk/* \
 && rm -rf /root/.glide/ \
 && rm -rf ${GOPATH}
 
USER user

WORKDIR /home/user/

EXPOSE 8080

VOLUME ["/tmp/", "/home/user/acdc"]
 
CMD ["/usr/local/bin/acdc"]
