(ns com.hello.haneda.client
  "Haneda websocket client.

  Example usage:

  (require '[com.hello.haneda.client :as client])
  (def sense-id \"my-sense\")
  (def key (.getBytes sense-id))  ; Haneda server fakes id using the sense id
  (def url \"ws://localhost:11000/dispatch\")
  (def connection (client/connect sense-id url))

  ;; Make a payload of bytes 0-99 (in reality will be serialized protobuf)
  (def payload (byte-array (range 100)))

  ;; Wait for a message to be sent to server (leave off the `@` for fully async)
  @(client/send-message connection key payload)

  ;; Get the next decoded message from the server
  @(client/receive-message connection key)

  (.close connection)"
  (:require
    [aleph.http :as http]
    [com.hello.haneda
      [message-format :as mf]
      [protobuf :as pb]]
    [manifold [stream :as s] [deferred :as d]]))

;; Allow .close to be called on connection to shutdown websocket connection
(defrecord Connection [conn last-id]
  java.io.Closeable
  (close [_]
    (.close conn)))

(defn connect
  "Returns a Connection record that can be passed to receive-message and send-message."
  ([sense-id url]
    (map->Connection
      {:conn @(http/websocket-client url {:headers {"X-Hello-Sense-Id" sense-id}})
       :last-id (atom 0)}))
  ([sense-id]
    (connect sense-id "ws://localhost:11000/dispatch")))

(defn receive-message
  "Given a Connection record and a key, return a Deferred of the next decoded message from the server."
  [{:keys [conn]} key]
  (d/chain
    (s/take! conn)
    (partial mf/decode-message key)))

(defn send-message
  "Given a Connection record, key, and payload, send an encoded message to the server.

  Will increment to the next message-id using the state of the connection."
  [{:keys [conn last-id]} ^bytes key ^bytes request-payload]
  (let [message (mf/encode-message key (:messeji pb/pb-type) request-payload (swap! last-id inc))]
    (s/put! conn message)))
