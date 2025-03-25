### Backend Dockerfile (Go Fiber)
FROM golang:1.21 AS builder
WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN go build -o shorten

### Final Image
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/shorten .
COPY backend/urls.db ./

EXPOSE 3000
CMD ["./shorten"]
