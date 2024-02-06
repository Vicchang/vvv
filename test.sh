#!/bin/sh

set -e
echo mode: set > coverage.out

for dir in ./...;
do 
    go test -coverprofile=coverage.tmp -p 1 $dir -args $@

if [ -f coverage.tmp ]
then
    cat coverage.tmp | tail -n +2 >> coverage.out
    rm coverage.tmp
fi
done

go tool cover -func coverage.out
rm coverage.out
