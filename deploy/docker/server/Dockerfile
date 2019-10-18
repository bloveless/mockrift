FROM golang:1.13-alpine as builder

ENV CGO_ENABLED 0

EXPOSE 3499

RUN apk update \
    && apk add --no-cache git \
    && rm -rf /var/cache/apk/* \
    && addgroup -g 1000 appuser \
    && adduser -S -u 1000 -G appuser appuser \
    && mkdir -p /home/appuser/app \
    && chown -R appuser:appuser /home/appuser

COPY --chown=appuser:appuser ./mockrift /home/appuser/app

USER appuser

WORKDIR /home/appuser/app

RUN go get -u github.com/go-delve/delve/cmd/dlv \
    && go get -u github.com/cortesi/modd/cmd/modd

RUN go mod download \
    && go build -o ./bin/mockrift cmd/web/*.go

CMD ["./bin/mockrift"]

####---- Production Image ----####

FROM alpine as production

EXPOSE 3499

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appuser /home/appuser/app /home/appuser/app
WORKDIR /home/appuser/app

USER appuser

CMD ["./bin/mockrift"]
