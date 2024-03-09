# Use multi-stages to avoid large image.
# It's required to have only docker image with binary, other code files and dependencies are actually not required.

# build stage
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]