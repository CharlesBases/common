#!/bin/zsh

gopath=$GOPATH
gopath=${gopath%:*}

cd $gopath/src/github.com/

git clone https://github.com/wg/wrk.git

cd wrk

make

mv wrk $gopath/bin/

