FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crawler .

FROM scratch
COPY --from=builder /app/crawler /crawler
ENTRYPOINT ["/crawler"]
