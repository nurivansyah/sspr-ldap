# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /app/sspr-ldap ./

# Production stage
FROM alpine:3.23
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/sspr-ldap /app/sspr-ldap
COPY --from=builder /src/templates /app/templates

ENV PORT=8080
EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app/sspr-ldap"]
