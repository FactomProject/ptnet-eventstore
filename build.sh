#!/usr/bin/env bash

PACKAGES=(ptnet contract identity finite x wsapi)

for PKG in ${PACKAGES[*]} ; do
  go build ./$PKG
  go install ./$PKG
done
