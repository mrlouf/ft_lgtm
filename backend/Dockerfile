FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o lgtm-backend ./cmd/server/main.go

RUN wget https://github.com/bytecodealliance/javy/releases/download/v9.0.0/javy-x86_64-linux-v9.0.0.gz && \
    gunzip javy-x86_64-linux-v9.0.0.gz && \
    chmod +x javy-x86_64-linux-v9.0.0 && \
    mv javy-x86_64-linux-v9.0.0 /usr/local/bin/javy

# TinyGo for now since we only need to compile Go
# but this may change in the future if we want to support other languages like Python, JavaScript, etc.
FROM tinygo/tinygo:0.41.1

WORKDIR /root/

ENV ENV=production

COPY --from=builder /app/lgtm-backend .
COPY --from=builder /usr/local/bin/javy /usr/local/bin/javy

EXPOSE 4242

CMD ["./lgtm-backend"]