FROM golang:1.12-alpine as builder

ENV CGO_ENABLED 0

EXPOSE 3499 40000

RUN apk update \
    && apk add --no-cache git make gcc musl-dev \
    && rm -rf /var/cache/apk/* \
    && adduser -D appuser \
    && mkdir -p /home/appuser/app \
    && chown -R appuser:appuser /home/appuser

RUN go get -u github.com/go-delve/delve/cmd/dlv \
	&& go get -u github.com/cespare/reflex

COPY --chown=appuser:appuser ./mockrift /home/appuser/app
WORKDIR /home/appuser/app

USER appuser

RUN make

CMD ["make", "run-reflex"]

####---- Production Image ----####

FROM alpine as production

EXPOSE 3499

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appuser /home/appuser/app /home/appuser/app
WORKDIR /home/appuser/app

USER appuser

ENTRYPOINT ["sh", "-c", "./bin/mockrift"]
