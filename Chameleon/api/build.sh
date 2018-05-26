#! /bin/sh

linuxRelease=chameleon.linux

GOOS=linux GOARCH=amd64 go build -o $linuxRelease
