##!/bin/sh
##mac only
#
#CC=x86_64-linux-musl-gcc \
#CXX=x86_64-linux-musl-g++ \
#go install -tags static
#
#CC=x86_64-linux-musl-gcc \
#CXX=x86_64-linux-musl-g++ \
#GOARCH=amd64 \
#GOOS=linux \
#CGO_ENABLED=1 \
#go build  -ldflags "-linkmode external -extldflags -static" ./main.go
##go build -o main ./main.go
#VER=$(python3 ./.deploy/getVersion.py ./core/Versions.go)
#mv ./main ./.deploy/gl2gl
#echo $VER
#docker buildx build  --platform linux/amd64   --push -t harbor.wns8.io/public/gl2gl:$VER   -f ./.deploy/Dockerfile ./.deploy
#rm -rf ./.deploy/gl2gl
#

CC=x86_64-linux-musl-gcc \
CXX=x86_64-linux-musl-g++ \
GOARCH=amd64 \
GOOS=linux \
CGO_ENABLED=1 \
go build  ./main.go
VER=$(python3 ./.deploy/getVersion.py ./core/Versions.go)
docker buildx build  --platform linux/amd64  --push -t harbor.wns8.io/public/gl2gl:$VER  -f ./.deploy/Dockerfile ./
