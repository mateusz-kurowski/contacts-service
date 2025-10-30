FROM sqlc/sqlc:latest AS sqlc
WORKDIR /app
COPY sqlc.yml ./
COPY db/ ./db/
RUN sqlc generate

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=sqlc /app/db ./db

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o contactsService .

FROM scratch

COPY --from=builder /app/contactsService /contactsService

EXPOSE 8080

ENTRYPOINT ["/contactsService"]
