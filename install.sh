#!/bin/bash
#export GOPATH=/opt/go
#curpath=`pwd`
echo "building the mock service"
#cd $curpath
echo "Listing $GOPATH/bin"
mkdir -p "$GOPATH/bin"
ls -l "$GOPATH/bin"
go install -v
ret=$?
if [ $ret != 0 ]
then
    exit $ret
else
    echo "After Installing Mocker"
    ls -l "$GOPATH/bin"
fi


