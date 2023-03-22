FROM golang:alpine3.17 AS builder

RUN apk update && apk add --no-cache git

ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/app

COPY . .
RUN go env -w GO111MODULE=off
RUN go mod init dockergo
RUN go get -d -v ./...
RUN go get github.com/0x9ef/openai-go
RUN go get github.com/go-telegram-bot-api/telegram-bot-api

RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app main.go

FROM scratch
COPY --from=builder /go/bin/app/ /go/bin/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["go/bin/app"]