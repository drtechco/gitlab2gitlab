FROM alpine:latest

RUN apk update && apk add tzdata

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN echo "Asia/Shanghai"  > /etc/timezone

RUN apk update

RUN apk add go git

COPY ./main /usr/bin/gl2gl

COPY ./.deploy/entrypoint.sh /entrypoint.sh

#COPY ./gl2gl /usr/bin
#COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT /entrypoint.sh

WORKDIR /

EXPOSE 8081