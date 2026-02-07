FROM golang:1.25-alpine AS builder

WORKDIR /chirpy

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 1001 \
    chirpyuser

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/chirpy .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /bin/chirpy /chirpy

USER chirpyuser:chirpyuser

ENTRYPOINT [ "/chirpy" ]
