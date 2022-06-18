#!/bin/sh

cd /opt

rm -rf ./vendor/github.com/libgit2/git2go/v33
ln -s $(pwd)/git2go $(pwd)/vendor/github.com/libgit2/git2go/v33

cd /opt/git2go

make install-static

cd ../

go install -tags static

go build  -tags static -a  ./main.go

mv ./main /usr/bin/gl2gl

mv ./.deploy/entrypoint.sh /entrypoint.sh