FROM golang:alpine as builder
LABEL maintainer="amoCRM"
RUN apk update && apk add --no-cache git
WORKDIR ./app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8080

