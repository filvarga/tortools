#
FROM alpine:latest AS app
LABEL maintainer="filipvarga89@gmail.com"

RUN apk add --no-cache rtorrent

ENV USER_HOME=/app
ENV USER_NAME=app

ARG IDU=1000
ARG IDG=1000

RUN addgroup -g $IDG $USER_NAME && \
adduser -h $USER_HOME -G $USER_NAME -u $IDU -D $USER_NAME

COPY ./files/rtorrent.rc $USER_HOME/.rtorrent.rc

ADD ./files/bin /usr/local/bin

WORKDIR $USER_HOME

USER $IDU:$IDG

RUN mkdir -p session downloads

CMD ["/usr/local/bin/entrypoint.sh"]

#
FROM nginx:latest AS web
LABEL maintainer="filipvarga89@gmail.com"

ADD ./files/etc/nginx/conf.d /etc/nginx/conf.d

RUN mkdir -p /app/session
