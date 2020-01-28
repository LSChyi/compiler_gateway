FROM golang:1.13.6-alpine3.11 AS builder
RUN mkdir -p /go/compiler_gateway
WORKDIR /go/compiler_gateway
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build
#-------------------------------------
FROM alpine:3.11.3
COPY --from=builder /go/compiler_gateway/compiler_gateway .
ENTRYPOINT ["./compiler_gateway"]
