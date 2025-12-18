ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.24-alpine AS builder

ARG APP_RELATIVE_PATH

COPY .. /data/app
WORKDIR /data/app

RUN rm -rf /data/app/bin/
RUN go mod tidy && go build -ldflags="-s -w" -o ./bin/server ${APP_RELATIVE_PATH}
RUN mv config /data/app/bin/


FROM ${REGISTRY}/alpine:3.16

ARG APP_ENV
ENV APP_ENV=${APP_ENV}

WORKDIR /data/app
COPY --from=builder /data/app/bin /data/app

EXPOSE 8291
ENTRYPOINT ./server -conf ./config/${APP_ENV}.yml
# ENTRYPOINT [ "./server", "-conf", "./config/${APP_ENV}.yml" ]
