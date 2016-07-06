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
go get github.com/golang/protobuf/proto

cd $GOPATH/src/github.com/hello/
git clone git@github.com:hello/haneda.git
cd haneda
go run main.go -c config-file-here.conf
```

