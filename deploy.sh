#!/usr/bin/env sh

gox -osarch="linux/amd64"
scp haneda_linux_amd64 dev:/tmp/
