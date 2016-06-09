# haneda

Websocket server/lib in clojure.

## Running the server
```bash
lein run -m com.hello.haneda.server
```

This will start the server on port 11000.

## Using the client
Start a clojure REPL.
```bash
lein repl
```

```clojure
;; Load the client namespace
(require '[com.hello.haneda.client :as client])

;; Haneda server doesn't currently use key store so we can fake it
(def sense-id "my-sense")

;; Haneda server fakes id using the sense id
(def key (.getBytes sense-id))

(def url "ws://localhost:11000/dispatch")

;; Opens the websocket connection
(def connection (client/connect sense-id url))

;; Make a payload of bytes 0-99 (in reality will be serialized protobuf)
(def payload (byte-array (range 100)))

;; Wait for a message to be sent to server (leave off the `@` for fully async)
;; You should be able to see Haneda server print the preamble and payload
@(client/send-message connection key payload)

;; Get the next decoded message from the server
@(client/receive-message connection key)

;; Close the websocket connection
(.close connection)
```
