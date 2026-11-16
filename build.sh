#!/bin/bash -ex

export GO111MODULE=on

EXENAME=$1
if [ -z "$EXENAME" ]; then
  EXENAME=unoclient
fi

VERSION=1.0.7
if [ -z "$REVISION" ]; then
  REVISION=`git rev-parse --short HEAD`
fi

BUILDTIME=`date -u "+%Y-%m-%dT%H:%M:%SZ"`

PKG=github.com/askasoft/gounoconv
LDF="-X ${PKG}.Version=${VERSION} -X ${PKG}.Revision=${REVISION} -X ${PKG}.Buildtime=${BUILDTIME}"

go build -ldflags "${LDF}" -o ${EXENAME}

