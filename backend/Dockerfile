FROM golang:1.26-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o lgtm-backend ./cmd/server/main.go

# Install TinyGo
RUN apt-get update && apt-get install -y wget && \
    wget https://github.com/tinygo-org/tinygo/releases/download/v0.41.1/tinygo_0.41.1_amd64.deb && \
    dpkg -i tinygo_0.41.1_amd64.deb && \
    rm tinygo_0.41.1_amd64.deb

# Install Javy
RUN wget https://github.com/bytecodealliance/javy/releases/download/v9.0.0/javy-x86_64-linux-v9.0.0.gz && \
    gunzip javy-x86_64-linux-v9.0.0.gz && \
    chmod +x javy-x86_64-linux-v9.0.0 && \
    mv javy-x86_64-linux-v9.0.0 /usr/local/bin/javy

# Install Emcc
RUN apt-get update && apt-get install -y emscripten

ENV ENV=production

EXPOSE 4242

CMD ["./lgtm-backend"]