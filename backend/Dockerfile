FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o lgtm-backend ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/lgtm-backend .

EXPOSE 4242

CMD ["./lgtm-backend"]