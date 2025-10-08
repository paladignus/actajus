FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/http/server.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /
EXPOSE 8080 9090
CMD ["/main"]



# FROM golang:alpine AS builder
#
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/http/server.go
#
# FROM scratch 
# COPY --from=builder /app/main /
# EXPOSE 8080
# CMD [ "/main" ]
