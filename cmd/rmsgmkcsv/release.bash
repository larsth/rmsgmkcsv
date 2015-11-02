#!/bin/bash

mkdir -p ./release
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./release/rmsgmkcsv-linux-i386 rmsgmkcsv.go
env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o ./release/rmsgmkcsv-linux-ARM6 rmsgmkcsv.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/rmsgmkcsv-linux-amd64 rmsgmkcsv.go
env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./release/rmsgmkcsv-windows-i386.exe rmsgmkcsv.go
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/rmsgmkcsv-windows-amd64.exe rmsgmkcsv.go
cd ./release
sha256sum -b * > SHA256SUMS
sha512sum -b * > SHA512SUMS
cd ..