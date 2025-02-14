FROM golang:alpine AS build_container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o app .
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
FROM alpine:latest

WORKDIR /app
COPY --from=build_container /app/app .
COPY --from=build_container /go/bin/dlv .
EXPOSE 8000 2345
ENTRYPOINT ["./dlv", "--listen=:2345", "--continue", "--accept-multiclient", "--headless=true", "--api-version=2", "exec", "./app"]