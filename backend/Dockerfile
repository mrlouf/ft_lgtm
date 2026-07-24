FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o lgtm-backend ./cmd/server/main.go

# TinyGo for now since we only need to compile Go
# but this may change in the future if we want to support other languages like Python, JavaScript, etc.
FROM tinygo/tinygo:0.41.1

WORKDIR /root/

ENV ENV=production

COPY --from=builder /app/lgtm-backend .

EXPOSE 4242

CMD ["./lgtm-backend"]