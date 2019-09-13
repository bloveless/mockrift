FROM golang:1.13 as builder

ENV CGO_ENABLED 1
ENV GIN_MODE release

EXPOSE 3499 40000

RUN apt-get update \
    && apt-get install -y git make gcc \
    && rm -rf /var/cache/apk/* \
    && useradd appuser \
    && mkdir -p /home/appuser/app \
    && chown -R appuser:appuser /home/appuser \
    && rm -rf /var/lib/apt/lists/*

RUN go get -u github.com/go-delve/delve/cmd/dlv \
    && go get -u github.com/cortesi/modd/cmd/modd

COPY --chown=appuser:appuser ./mockrift /home/appuser/app
WORKDIR /home/appuser/app

USER appuser

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

ENTRYPOINT ["sh", "-c", "./bin/mockrift"]
