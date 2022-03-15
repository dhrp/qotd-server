# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY wisdom.txt ./

RUN go build -o /qotd

EXPOSE 3333
EXPOSE 80

CMD [ "/qotd", "wisdom.txt" ]