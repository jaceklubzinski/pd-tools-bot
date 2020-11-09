FROM golang:1.13-alpine as builder

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    appuser

RUN apk update
RUN apk add \
  g++ \
  git \
  musl-dev \
  go \
  tesseract-ocr-dev

WORKDIR /app
COPY go.mod ./

RUN apk add build-base
COPY . /app

RUN go build -o pd-tools-bot .

# Build a small image
FROM alpine

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/pd-tools-bot /app/pd-tools-bot

USER appuser:appuser

ENTRYPOINT ["/app/pd-tools-bot"]
