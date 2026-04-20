FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/bot ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S bot && adduser -S -G bot bot

WORKDIR /app
COPY --from=builder /out/bot /app/bot

USER bot

ENTRYPOINT ["/app/bot"]
