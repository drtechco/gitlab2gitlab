FROM alpine:latest

RUN apk update && apk add tzdata

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN echo "Asia/Shanghai"  > /etc/timezone

RUN apk update

RUN apk add pkgconfig go cmake make python3

COPY . /opt

WORKDIR /opt

RUN rm -rf /opt/vendor/github.com/libgit2/git2go/v34

#RUN /opt/.deploy/build2.sh
RUN  ls -la & go get github.com/libgit2/git2go/v34

RUN ln -s /opt/git2go /opt/vendor/github.com/libgit2/git2go/v34


WORKDIR /opt/git2go

RUN rm -rf CMakeLists.txt & rm -rf ./static-build

RUN make build-libgit2-static

RUN make install-static

WORKDIR /opt/

RUN go install -tags static

RUN go build  -tags static -a  ./main.go

RUN mv ./main /usr/bin/gl2gl

RUN mv ./.deploy/entrypoint.sh /entrypoint.sh

#COPY ./gl2gl /usr/bin
#COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT /entrypoint.sh

WORKDIR /

RUN rm -rf /opt

EXPOSE 8081