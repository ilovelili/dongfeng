#!/bin/bash

set -e

read -r -p "update server? [y/n] " response
case "$response" in
    [yY][eE][sS]|[yY])
        echo "building server"
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dongfeng .

        echo "scping server"
        scp dongfeng dongfeng:/root/dongfeng/bin/
        scp dongfeng dongfeng-2:/root/dongfeng/bin/

        echo "cleaning server"
        rm dongfeng 
        ;;
    *)
esac

echo "done"