FROM golang:alpine AS build_container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o app .
FROM alpine:latest
WORKDIR /app
COPY --from=build_container /app/app .
EXPOSE 8000
CMD ["./app"]