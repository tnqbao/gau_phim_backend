FROM golang:1.23-alpine AS builder
WORKDIR /gau_phim
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -tags 'prod' -o main .

FROM alpine:latest
WORKDIR /gau_phim
COPY --from=builder /gau_phim/main .
EXPOSE 8083
CMD ["./main"]