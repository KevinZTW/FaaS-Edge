
# temporarily change to use ubuntu image to install couchdb
FROM ubuntu:bionic

# ENV GO_VERSION=1.20

# RUN apt-get update
# RUN apt-get install -y wget git gcc

# #  install golang
# RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

# RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
# RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

# ENV GOPATH /go
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"


# install couchdb =====================

RUN apt update && apt install -y curl apt-transport-https gnupg

WORKDIR     /

COPY keys.asc /tmp/keys.asc
RUN cat /tmp/keys.asc

RUN  gpg --dearmor < /tmp/keys.asc | tee /usr/share/keyrings/couchdb-archive-keyring.gpg >/dev/null 2>&1

# RUN source /etc/os-release
# Since can't use source in dockerfile, hardcode version as temporary workaround
ENV VERSION_CODENAME=bionic

RUN echo "deb [signed-by=/usr/share/keyrings/couchdb-archive-keyring.gpg] https://apache.jfrog.io/artifactory/couchdb-deb/ ${VERSION_CODENAME} main" | tee /etc/apt/sources.list.d/couchdb.list >/dev/null

# Couchdb install requires user input, so set noninteractive as workaround
ENV DEBIAN_FRONTEND noninteractive

RUN apt update -y

RUN apt install -y apt-utils
RUN apt install -f -y apt-utils couchdb

# =====================================

# build the app
WORKDIR     /app
COPY . ./
RUN go mod download
RUN go build -o main



EXPOSE 3037
EXPOSE 6379

CMD ./main


# FROM golang:1.20-alpine3.17

# WORKDIR     /app
# COPY . ./
# RUN go mod download
# RUN go build -o main

# RUN apk add --update redis && \
#     rm -rf /var/cache/apk/* && \
#     mkdir /data && \
#     chown -R redis:redis /data && \
#     sed -i 's#logfile /var/log/redis/redis.log#logfile ""#i' /etc/redis.conf && \
#     sed -i 's#daemonize yes#daemonize no#i' /etc/redis.conf && \
#     sed -i 's#dir /var/lib/redis/#dir /data#i' /etc/redis.conf && \
#     echo -e "# placeholder for local options\n" > /etc/redis-local.conf && \
#     echo -e "include /etc/redis-local.conf\n" >> /etc/redis.conf

# EXPOSE 3037
# EXPOSE 6379

# CMD redis-server & ./main