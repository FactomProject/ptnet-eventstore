#!/usr/bin/env bash

PACKAGES=(ptnet contract wsapi gen blockchain identity sim finite x )

for PKG in ${PACKAGES[*]} ; do
  go build ./$PKG
  go install ./$PKG
done
