#!/usr/bin/env sh

curl -i -N  \
-X "GET" \
-H "X-Hello-Sense-Id: XXXXXXXXXXXXXXXX" \
-H "Connection: upgrade" \
-H "Upgrade: websocket" \
-H "Host: localhost:8082" \
-H "Sec-Websocket-Version: 13" \
-H "Sec-Websocket-Key: foo" \
http://localhost:8082/bench