FROM golang:1.20-alpine

# Required because go requires gcc to build
RUN apk add build-base git inotify-tools
RUN echo $GOPATH

RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . /ticket_api

WORKDIR /ticket_api

RUN go mod download

ENTRYPOINT CompileDaemon --build="go build main.go" --command="./main app:serve"
