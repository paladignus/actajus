FROM golang:1.25.5-alpine AS builder

# Instalar dependências necessárias para compilar
RUN apk add --no-cache git ca-certificates tzdata bash

WORKDIR /app

# Copiar arquivos de dependência primeiro para aproveitar o cache da camada
COPY go.mod go.sum ./

# Baixar dependências com tentativas extras em caso de falha
RUN --mount=type=cache,target=/go/pkg/mod \
  for i in 1 2 3; do go mod download && break || sleep 5; done && \
  go mod verify

# Copiar código fonte
COPY . .

# Compilar o binário com otimizações
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -buildid=" -o main cmd/main.go

# Imagem final
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copiar o binário compilado
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
