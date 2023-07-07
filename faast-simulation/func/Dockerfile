FROM golang:1.20-alpine3.17

WORKDIR     /app
COPY . ./
RUN go mod download
RUN go build -o main

EXPOSE 3038

CMD ./main