#!/usr/bin/env bash

#
# Building and packaging binaries for supported architectures
#

if [[ $# -eq 0 ]]; then
    echo 'Build version needed. Usage: build.sh <version>'
    exit 1
fi

version="${1}"

rootDir=$(dirname "${0}")/..
buildDir=${rootDir}/build

rm -rf ${buildDir}

os=darwin
arch=amd64
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

arch=arm64
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

os=linux
arch=386
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

arch=amd64
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

os=windows
arch=386
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail.exe ${rootDir}/*.go
zip -j ${buildDir}/nomtail_${version}_${os}_${arch}.zip ${buildDir}/nomtail.exe

arch=amd64
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail.exe ${rootDir}/*.go
zip -j ${buildDir}/nomtail_${version}_${os}_${arch}.zip ${buildDir}/nomtail.exe
