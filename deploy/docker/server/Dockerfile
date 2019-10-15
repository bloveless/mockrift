FROM golang:1.13-alpine as builder

ENV CGO_ENABLED 0

EXPOSE 3499 40000

RUN apk update \
    && apk add --no-cache git \
    && rm -rf /var/cache/apk/* \
    && addgroup -g 1000 appuser \
    && adduser -S -u 1000 -G appuser appuser \
    && mkdir -p /home/appuser/app \
    && chown -R appuser:appuser /home/appuser

COPY --chown=appuser:appuser ./mockrift /home/appuser/app

USER appuser

RUN mkdir -p /tmp/dependencies

WORKDIR /tmp/dependencies

# go mod init dependencies is an unfortunate hack to get go to respect the dependencies within modd
# seems like go get will only respect the dependencies of a package if it is in mod mode which means it need to be in
# a go mod initialized directory
RUN go mod init dependencies \
    && go get -u github.com/go-delve/delve/cmd/dlv \
    && go get -u github.com/cortesi/modd/cmd/modd

WORKDIR /home/appuser/app

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
