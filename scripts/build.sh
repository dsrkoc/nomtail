#!/usr/bin/env bash

#
# Building and packaging binaries for supported architectures
#

version=${1:-v1.0.0}

rootDir=$(dirname "${0}")/..
buildDir=${rootDir}/build

rm -rf ${buildDir}

os=darwin
arch=amd64
GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

# arch=arm64
# GOOS=${os} GOARCH=${arch} go build -o ${buildDir}/nomtail ${rootDir}/*.go
# tar czf ${buildDir}/nomtail_${version}_${os}_${arch}.tar.gz -C ${buildDir} nomtail

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
