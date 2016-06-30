#!/usr/bin/env sh
make linux
scp _build/haneda_linux_amd64 haneda-dev:/tmp/
scp server.dev.conf haneda-dev:/tmp/
