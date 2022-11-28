# Build
FROM golang:1.19-alpine AS builder
ARG CGO_ENABLED=0
ARG GO111MODULE=on
ARG GOARCH=amd64
ARG GOOS=linux

RUN ln -fs /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

WORKDIR /source

COPY . .
RUN go mod download

RUN go build -o /source/app cmd/app/main.go

# Deploy
FROM alpine:3.14 as deployer

RUN ln -fs /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

WORKDIR /source

COPY --from=builder /source/app /source/app

COPY scripts /source/scripts
COPY configs /source/configs
COPY migrations /source/migrations

RUN chmod +x /source/scripts/run.sh

CMD ["/source/scripts/run.sh", "app"]
