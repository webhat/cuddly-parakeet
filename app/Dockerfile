FROM golang

ENV GO111MODULE=on

# DEBUG
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y netcat

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080
ENTRYPOINT ["/app/fantastic-lamp"]
