# gitlab2gitlab
```shell
git submodule update --init 
cd ./git2go
git submodule update --init 
make install-static
cd ..
go mod vendor
rm -rf ./vendor/github.com/libgit2/git2go/v34
ln -s $(pwd)/git2go $(pwd)/vendor/github.com/libgit2/git2go/v34
go install -tags static
```