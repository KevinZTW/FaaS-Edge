FROM golang:1.20-alpine3.17

WORKDIR     /app
COPY . ./
RUN go mod download
RUN go build -o main


RUN apk add --update redis && \
    rm -rf /var/cache/apk/* && \
    mkdir /data && \
    chown -R redis:redis /data && \
    sed -i 's#logfile /var/log/redis/redis.log#logfile ""#i' /etc/redis.conf && \
    sed -i 's#daemonize yes#daemonize no#i' /etc/redis.conf && \
    sed -i 's#dir /var/lib/redis/#dir /data#i' /etc/redis.conf && \
    echo -e "# placeholder for local options\n" > /etc/redis-local.conf && \
    echo -e "include /etc/redis-local.conf\n" >> /etc/redis.conf

EXPOSE 3038
EXPOSE 6379

CMD redis-server & ./main