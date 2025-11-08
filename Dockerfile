FROM sqlc/sqlc:1.30.0 AS sqlc_binary

FROM golang:1.25-alpine AS builder

COPY --from=sqlc_binary /workspace/sqlc /usr/local/bin/sqlc

WORKDIR /app

COPY sqlc.yml .
COPY sql/ ./sql/
RUN sqlc generate

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o contactsService .

FROM scratch

COPY --from=builder /app/contactsService /contactsService

EXPOSE 33500

ENTRYPOINT ["/contactsService"]
