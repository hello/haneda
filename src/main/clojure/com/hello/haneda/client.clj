(ns com.hello.haneda.client
  "Haneda websocket client."
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
