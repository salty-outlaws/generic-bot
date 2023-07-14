FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git make

WORKDIR /app
COPY . .

RUN make vendor
RUN make build

FROM alpine:3
RUN apk update \
    && apk add --no-cache curl wget \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder /app/bin/generic-bot /generic-bot

CMD ["/generic-bot"]

ENV PORT=8080
EXPOSE $PORT