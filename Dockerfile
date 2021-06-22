FROM golang:1.15.4-alpine as build

WORKDIR /go/src/github.com/pianzm/arr/
ADD . /go/src/github.com/pianzm/arr
COPY go.mod go.sum .env /go/src/github.com/pianzm/arr/

ENV BUILD_PACKAGES="git curl build-base make openssh"
ENV GO111MODULE=on

RUN apk update && apk add --no-cache $BUILD_PACKAGES \
      && mkdir /root/.ssh/ && mv id_rsa /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa && touch /root/.ssh/known_hosts \
      && ssh-keyscan github.com >> /root/.ssh/known_hosts \
      && git config --global url."git@github.com:".insteadOf "https://github.com/" \
      && make app-linux \
      && apk del $BUILD_PACKAGES

FROM alpine:3.13.5
RUN apk update \
      && apk add rsyslog \
      && apk add supervisor \
      && apk add tzdata
ARG BINARY_PATH=/go/src/github.com/pianzm/arr
RUN mkdir -p /var/log/

ADD _build/supervisord.conf /etc/supervisord.conf

RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

VOLUME ["/var/log"]
COPY --from=build $BINARY_PATH/.env $BINARY_PATH/.env

EXPOSE 8082
EXPOSE 8081

COPY --from=build $BINARY_PATH/app-linux $BINARY_PATH/app-linux

ENTRYPOINT ["sh", "-c", "supervisord -nc /etc/supervisord.conf"]