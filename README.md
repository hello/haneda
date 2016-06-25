# Haneda

Websocket Server for Sense 1.5


Message structure

```
+------------+-----------+----------+------------+-----------+
| header len | header pb | body len | body    pb | hmac sig  |
+------------+-----------+----------+------------+-----------+
```

# Pub/Sub

Redis is used a low latency pub/sub broker.
Out of band messages will be published through redis, and dispatched to the given sense


# How to run

```
go run $GOPATH/src/github.com/hello/haneda/server/server.go
```