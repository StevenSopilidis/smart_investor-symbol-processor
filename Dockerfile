FROM golang:1.23-alpine AS Build

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /symbol-processor main.go

FROM alpine:latest

WORKDIR /
COPY --from=Build /symbol-processor /symbol-processor

ENTRYPOINT [ "/symbol-processor" ]