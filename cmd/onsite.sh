#!/bin/bash

set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

echo "example: onsite.sh prod /path/to/ec2/key.pem"

LDFLAGS="-s -w"
TM=`date +%F@%T@%Z`
TM="${TM//:/-}" # change TM format to avoid extract error
OUT=gomail-$TM

GOARCH=amd64

CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
echo "${G}$OUT built${W}"

PKG_DIR=gomail
rm -rf $PKG_DIR
mkdir -p $PKG_DIR
mv $OUT $PKG_DIR
cp *.json $PKG_DIR

PKG_NAME=$OUT.tar.gz
tar -czvf ./$PKG_NAME $PKG_DIR

# send the package to AWS EC2, $2 is the absolute path of the key (pem)

if [ -f "$2" ]; then
    if [[ $1 == 'test' ]]; then IP='3.24.210.59'; fi
    if [[ $1 == 'prod' || $1 == 'product' ]]; then IP='54.66.224.148'; fi
        
    CD=`pwd`
    echo "scp -i $2 $CD/$PKG_NAME ubuntu@$IP:misc"
    scp -i $2 $CD/$PKG_NAME ubuntu@$IP:misc

    rm -rf $CD/$PKG_NAME

else
    echo "${Y}valid key file is not provided, cannot send package to EC2${W}"
fi