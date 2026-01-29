FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk add --no-cache curl unzip

WORKDIR /app

COPY --from=builder /app/main .

ARG VOICE_ZIP_URL
RUN test -n "$VOICE_ZIP_URL" || { echo "VOICE_ZIP_URL build arg is required"; exit 1; } \
    && curl -fSL -o /tmp/voice.zip "$VOICE_ZIP_URL" \
    && mkdir -p internal/quote/data \
    && unzip -qo /tmp/voice.zip -d /tmp/voice \
    && mv /tmp/voice/voice internal/quote/data/audio \
    && rm -rf /tmp/voice.zip /tmp/voice \
    && apk del curl unzip

EXPOSE 3000

CMD ["./main"]
