FROM golang:1.13-alpine as builder

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

COPY --from=builder /app/pd-tools-bot /app/pd-tools-bot

ENTRYPOINT ["/app/pd-tools-bot"]
