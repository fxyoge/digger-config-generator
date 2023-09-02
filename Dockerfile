FROM docker.io/golang:1.21.0-alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/digger-config-generator

FROM scratch AS release
WORKDIR /app
COPY --from=builder /app/digger-config-generator /app/digger-config-generator

ENTRYPOINT ["/app/digger-config-generator"]
