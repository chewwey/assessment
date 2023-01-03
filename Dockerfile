FROM golang:1.19-alpine as build-base

# Build base for Go application
WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o ./out/server .

# ====================

FROM alpine:3.16.2

# Copy built executable and run it
COPY --from=build-base /app/out/server /app/server

CMD ["/app/server"]