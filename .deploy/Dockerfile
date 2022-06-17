FROM alpine:latest

RUN apk update && apk add tzdata

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN echo "Asia/Shanghai"  > /etc/timezone

COPY ./gl2gl /usr/bin
COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT /entrypoint.sh

EXPOSE 8081